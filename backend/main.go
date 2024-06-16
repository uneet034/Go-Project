package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/patrickmn/expirable-cache"
)

//  represents an entry in the LRU cache
type CacheEntry struct {
	Key         string    `json:"key"`
	Value       string    `json:"value"`
	Expiration  time.Time `json:"expiration"`
}

// LRUCache represents the LRU cache structure
type LRUCache struct {
	mutex       sync.Mutex
	cache       map[string]*expirable.CacheItem
	lruList     []string // Maintains the order of keys for LRU eviction
	capacity    int
}

// NewLRUCache creates a new instance of LRUCache
func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		cache:    make(map[string]*expirable.CacheItem),
		lruList:  make([]string, 0),
		capacity: capacity,
	}
}

// Get retrieves a value from the cache
func (c *LRUCache) Get(key string) (*CacheEntry, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	item, found := c.cache[key]
	if !found || item.Expired() {
		return nil, fmt.Errorf("key '%s' not found in cache", key)
	}

	return &CacheEntry{
		Key:        key,
		Value:      item.Value().(string),
		Expiration: item.Expiration(),
	}, nil
}

// Set add or update a value in the cache with expiration time
func (c *LRUCache) Set(key, value string, expiration time.Duration) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Check if key already exists, update if it does
	if item, found := c.cache[key]; found {
		item.Update(value, expiration)
		return nil
	}

	// Add a new item
	c.cache[key] = expirable.New(value, expiration)
	c.lruList = append(c.lruList, key)

	// Maintain LRU order
	if len(c.lruList) > c.capacity {
		oldKey := c.lruList[0]
		c.lruList = c.lruList[1:]
		delete(c.cache, oldKey)
	}

	return nil
}

// Delete removes a key from the cache
func (c *LRUCache) Delete(key string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	_, found := c.cache[key]
	if !found {
		return fmt.Errorf("key '%s' not found in cache", key)
	}

	delete(c.cache, key)
	return nil
}

// API Handler

func getHandler(cache *LRUCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := mux.Vars(r)["key"]

		entry, err := cache.Get(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(entry)
	}
}

func setHandler(cache *LRUCache) http.HandlerFunc {
	type requestPayload struct {
		Key         string `json:"key"`
		Value       string `json:"value"`
		Expiration  int    `json:"expiration"` // in seconds
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var payload requestPayload
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		expiration := time.Duration(payload.Expiration) * time.Second
		err = cache.Set(payload.Key, payload.Value, expiration)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func deleteHandler(cache *LRUCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := mux.Vars(r)["key"]

		err := cache.Delete(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func main() {
	r := mux.NewRouter()

	// Create a new LRUCache with capacity 100
	cache := NewLRUCache(100)

	// Define API endpoints
	r.HandleFunc("/cache/{key}", getHandler(cache)).Methods("GET")
	r.HandleFunc("/cache", setHandler(cache)).Methods("POST")
	r.HandleFunc("/cache/{key}", deleteHandler(cache)).Methods("DELETE")

	fmt.Println("Starting server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}

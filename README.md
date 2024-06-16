# Go-Project

Project Name
LRU Cache with Golang Backend and React Frontend

Overview
This project implements a Least Recently Used (LRU) Cache system using Golang for the backend and React for the frontend. The backend provides RESTful APIs for GET, SET, and DELETE operations on the cache, while the frontend consumes these APIs to interact with the cache.

Features
Backend (Golang)
Implements an LRU Cache using concurrency-safe data structures.
Exposes REST APIs (GET, POST, DELETE) for cache operations.
Handles expiration of cache entries automatically.
Frontend (React)
Provides a user-friendly interface to interact with the cache.
Consumes backend APIs to perform GET, SET, and DELETE operations.
Displays real-time updates using WebSocket for current key-value pairs and their expiration times (Good to have).
Setup Instructions
Prerequisites
Golang: Ensure Golang is installed. You can download it from golang.org.
Node.js: Ensure Node.js and npm are installed. You can download them from nodejs.org.
Backend Setup
Clone this repository:

bash -

git clone <repository_url>
cd <repository_name>/backend
Install dependencies:

bash -

go mod tidy
Run the backend server:

bash -

go run main.go
The server will start at http://localhost:8080.

# Frontend Setup
Navigate to the frontend directory:

bash -

cd <repository_name>/frontend
Install dependencies:

bash -

npm install
Start the React development server:

bash - 

npm start
The React application will be accessible at http://localhost:3000.


Backend APIs

GET /cache/{key}

Retrieves the value associated with key from the cache.
Example: curl -X GET http://localhost:8080/cache/mykey
POST /cache

Adds or updates a key-value pair in the cache with a customizable expiration time.
Request body:
json

{
  "key": "mykey",
  "value": "myvalue",
  "expiration": 10 // expiration time in seconds
}

DELETE /cache/{key}

Deletes the key-value pair from the cache.

Frontend Interface
Open the React application in your browser (http://localhost:3000 by default).
Enter a key, value, and expiration time (in seconds) in the input fields.
Click Set to add the key-value pair to the cache.
Click Get to retrieve the value associated with a key from the cache.
Click Delete to remove a key from the cache.

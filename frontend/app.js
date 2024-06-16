import React, { useState } from 'react';
import './App.css';

function App() {
  const [key, setKey] = useState('');
  const [value, setValue] = useState('');
  const [expiration, setExpiration] = useState(5); // Default expiration in seconds
  const [response, setResponse] = useState('');

  const handleGet = async () => {
    try {
      const res = await fetch(`http://localhost:8080/cache/${key}`);
      const data = await res.json();
      setResponse(JSON.stringify(data));
    } catch (error) {
      setResponse(`Error: ${error.message}`);
    }
  };

  const handleSet = async () => {
    try {
      const res = await fetch(`http://localhost:8080/cache`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ key, value, expiration }),
      });
      if (res.ok) {
        setResponse('Key-Value pair added successfully.');
      } else {
        const data = await res.json();
        setResponse(`Error: ${JSON.stringify(data)}`);
      }
    } catch (error) {
      setResponse(`Error: ${error.message}`);
    }
  };

  const handleDelete = async () => {
    try {
      const res = await fetch(`http://localhost:8080/cache/${key}`, {
        method: 'DELETE',
      });
      if (res.ok) {
        setResponse('Key deleted successfully.');
      } else {
        const data = await res.json();
        setResponse(`Error: ${JSON.stringify(data)}`);
      }
    } catch (error) {
      setResponse(`Error: ${error.message}`);
    }
  };

  return (
    <div className="App">
      <h1>LRU Cache App</h1>
      <div>
        <label>
          Key:
          <input type="text" value={key} onChange={(e) => setKey(e.target.value)} />
        </label>
        <br />
        <label>
          Value:
          <input type="text" value={value} onChange={(e) => setValue(e.target.value)} />
        </label>
        <br />
        <label>
          Expiration (seconds):
          <input type="number" value={expiration} onChange={(e) => setExpiration(parseInt(e.target.value))} />
        </label>
        <br />
        <button onClick={handleSet}>Set</button>
        <button onClick={handleGet}>Get</button>
        <button onClick={handleDelete}>Delete</button>
      </div>
      <div className="response">
        <h2>Response:</h2>
        <p>{response}</p>
      </div>
    </div>
  );
}

export default App;

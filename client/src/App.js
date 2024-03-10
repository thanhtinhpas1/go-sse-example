import logo from './logo.svg';
import './App.css';
import { useEffect, useState } from 'react';
import { EventSource } from 'event-source-polyfill';

function App() {
  const [statuses, setStatuses] = useState([]);

  useEffect(() => {
    const es = new EventSource('http://127.0.0.1:3000/api/sse')
    es.onmessage = (msg) => {
      console.info("received message: ", msg)
      const data =JSON.parse(msg.data);
      setStatuses(statuses => [...statuses, data]);
    }

    es.onerror = (err ) => {
      // console.error("received error: ", err)
    }
  }, []);

  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <ul>
          {statuses.map(status => <li>App ID: {status.AppId}, AppTransId: {status.AppTransId}</li>)}
        </ul>
        <a
          className="App-link"
          href="https://reactjs.org"
          target="_blank"
          rel="noopener noreferrer"
        >
          Learn React
        </a>
      </header>
    </div>
  );
}

export default App;

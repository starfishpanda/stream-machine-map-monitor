import React, { useEffect, useState } from 'react'
import { GoogleMap, LoadScript, Marker, InfoWindow } from '@react-google-maps/api'
import './App.css'

interface GPS {
  lat: number;
  lon: number;
  alt: number;
}

interface Machine {
  id: number;
  location: GPS;
  fuelLevel: number;
  isPaused: boolean;
}

const mapContainerStyle = {
  width: '100%',
  height: '100vh'
};

const center = {
  lat: 47.695185, // Sammamish Valley
  lon: -122.145161
};

function App() {
  // Initialize state for machines and socket
  const [machines, setMachines ] = useState<Map<number,Machine>>(new Map())
  const [socket, setSocket] = useState<WebSocket | null>(null);
  
  // Connecto WebSocket for machine updates
  useEffect(() => {
    const ws = new WebSocket('ws://localhost:3001/machine');

    ws.onopen = () => {
      console.log('WebSocket connection established');
    }

    ws.onmessage = (event) => {
      const machine = JSON.parse(event.data) as Machine;
      console.log('Machine message received:', machine)
      setMachines((prev) => {
        const updated = new Map(prev)
        updated.set(machine.id,machine)
        return updated
      });
    };

    ws.onerror = (error) => {
      console.error('WebSocket error:', error);
    };

    ws.onclose = () => {
      console.log('WebSocket connection closed')
    };

    setSocket(ws);

    return () => {
      ws.close();
    }
  }, []);
  return (
    <>
      <h1>Vite + React</h1>
      <div className="card">
        <p>
          Edit <code>src/App.jsx</code> and save to test HMR
        </p>
      </div>
      <p className="read-the-docs">
        Click on the Vite and React logos to learn more
      </p>
    </>
  )
}

export default App

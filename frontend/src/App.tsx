import React, { useEffect, useState } from 'react'
import { GoogleMap, LoadScript, Marker, InfoWindow } from '@react-google-maps/api'
const mapsApiKey = import.meta.env.VITE_GOOGLE_MAPS_API_KEY;
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
  lng: -122.145161 // Maps API types uses LatLng
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
    <div className="App">
      <LoadScript googleMapsApiKey={mapsApiKey}>
        <GoogleMap
          mapContainerStyle={mapContainerStyle}
          center={center}
          zoom={13}
          >
            {/* Render machine markers on map */}
            {Array.from(machines.values()).map((machine)=>(
              <Marker
                key={machine.id}
                position={{
                  lat: machine.location.lat,
                  lng: machine.location.lon
                }}
                icon={{
                  path: window.google.maps.SymbolPath.CIRCLE,
                  scale: 8,
                  fillColor: machine.isPaused ? 'red' : 'green',
                  fillOpacity: 0.8,
                  strokeWeight: 2,
                  strokeColor: 'white'
                }}
                />
            ))}

        </GoogleMap>
      </LoadScript>
    </div>
  )
}

export default App

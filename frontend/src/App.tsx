// App.tsx
import React, { useEffect, useState } from 'react';
import { GoogleMap, LoadScript, Marker, InfoWindow } from '@react-google-maps/api';
import './App.css';

const mapsApiKey = import.meta.env.VITE_GOOGLE_MAPS_API_KEY;

interface GPS {
  lat: number;
  lon: number;
  alt: number;
}

interface Machine {
  id: number;
  location: GPS;
  fuel_level: number;
  is_paused: boolean;
}

const mapContainerStyle = {
  width: '100%',
  height: '50vh',
};

const center = {
  lat: 47.695185, // Sammamish Valley
  lng: -122.145161 // Maps API types uses LatLng
};

function App() {
  const [machines, setMachines] = useState<Map<number, Machine>>(new Map()); // machine.id: data
  const [selectedMachine, setSelectedMachine] = useState<Machine | null>(null);
  const [socketConnections, setSocketConnections] = useState<Map<number,WebSocket>>(new Map()); // machine.id: socketConnection

  // // Initialize machine on page load
  // useEffect(() => {
  //   const ws = new WebSocket('ws://localhost:3001/machine');
    
  //   ws.onopen = () => {
  //     console.log('WebSocket connection established');
  //   };
    
  //   ws.onmessage = (event) => {
  //     const machine = JSON.parse(event.data) as Machine;
  //     console.log('Machine message received:', machine);
      
  //     setMachines((prev) => {
  //       const updated = new Map(prev);
  //       updated.set(machine.id, machine);
  //       return updated;
  //     });
  //   };
    
  //   ws.onerror = (error) => {
  //     console.error('WebSocket error:', error);
  //   };
    
  //   ws.onclose = () => {
  //     console.log('WebSocket connection closed');
  //   };
    
  //   setSocket(ws);
    
  //   // Clean up WebSocket on unmount
  //   return () => {
  //     ws.close();
  //   };
  // }, []);

  // Function to pause/unpause a machine
  const togglePause = (machine: Machine) => {
    console.log(`Socket Connections in togglePause for ${machine.id}`, socketConnections)
    const socket = socketConnections.get(machine.id)

    if (!socket || socket.readyState !== WebSocket.OPEN) {
      console.error(`WebSocket not connected for machine ${machine.id}`);
      return;
    }
    
    try {
      const command = {
        type: machine.is_paused ? 'unpause' : 'pause',
        id: machine.id
      };
      
      console.log(`Sending ${command.type} command for machine ${command.id}`);
      socket.send(JSON.stringify(command));
      
      // Note: The actual state update will come through the WebSocket stream
    } catch (error) {
      console.error('Error sending command via WebSocket:', error);
    }
  };

  // Function to add a new machine by creating a new websocket connection
  const addMachine = () => {
    const newSocket = new WebSocket('ws://localhost:3001/machine');

    newSocket.onopen = () => {
      console.log(`New machine connection established.`);
    };
    
    newSocket.onmessage = (event) => {
      const newMachine = JSON.parse(event.data) as Machine;
      console.log('New machine data received from gRPC Server to WebSocket Proxy:', newMachine);

      // Update machines state
      setMachines((prev) => {
        const updatedMachines = new Map(prev); // machine.id: machine
        updatedMachines.set(newMachine.id,newMachine)
        return updatedMachines
      });

      // Update socketConnections state
      setSocketConnections((prev) => {
        const updatedSocketConnections = new Map(prev)

        if (!updatedSocketConnections.has(newMachine.id)){
          updatedSocketConnections.set(newMachine.id,newSocket)

        }
        return updatedSocketConnections
      });
    };
    
    newSocket.onerror = (error) => {
      console.error('Error in new machine connection:', error);
    };
    
    newSocket.onclose = () => {
      console.log('Machine connection closed');
    };
  };

  const removeMachine = (machineId: number) => {
    console.log(`Cannot remove machine ${machineId} - not supported by current gRPC service`);
    alert("Machine removal isn't supported by the current service implementation");

  };

  return (
    <div className="App">
      {/* Google Maps Section */}
      <div className="map-container">
        <LoadScript googleMapsApiKey={mapsApiKey}>
          <GoogleMap
            mapContainerStyle={mapContainerStyle}
            center={center}
            zoom={17}
          >
            {/* Render markers for each machine */}
            {Array.from(machines.values()).map((machine) => (
              <Marker
                key={machine.id}
                position={{
                  lat: machine.location.lat,
                  lng: machine.location.lon,
                }}
                icon={{
                  path: window.google.maps.SymbolPath.CIRCLE,
                  scale: 8,
                  fillColor: machine.is_paused ? 'red' : 'green',
                  fillOpacity: 0.8,
                  strokeWeight: 2,
                  strokeColor: 'white',
                }}
                onClick={() => setSelectedMachine(machine)}
              />
            ))}
            
            {/* Info window for selected machine */}
            {selectedMachine && (
              <InfoWindow
                position={{
                  lat: selectedMachine.location.lat,
                  lng: selectedMachine.location.lon,
                }}
                onCloseClick={() => setSelectedMachine(null)}
              >
                <div className="info-window">
                  <h3>Machine #{selectedMachine.id}</h3>
                  <p>Fuel Level: {selectedMachine.fuel_level.toFixed(2)}%</p>
                  <p>Status: {selectedMachine.is_paused ? 'Paused' : 'Active'}</p>
                  <p>
                    Location: {selectedMachine.location.lat.toFixed(6)}, {selectedMachine.location.lon.toFixed(6)}
                  </p>
                  <button
                    onClick={() => togglePause(selectedMachine)}
                    className={selectedMachine.is_paused ? 'play-btn' : 'pause-btn'}
                  >
                    {selectedMachine.is_paused ? 'Resume' : 'Pause'}
                  </button>
                </div>
              </InfoWindow>
            )}
          </GoogleMap>
        </LoadScript>
      </div>
      
      {/* Dashboard Section */}
      <div className="dashboard">
        <div className="dashboard-header">
          <h2>Dashboard</h2>
          <button className="add-machine-btn" onClick={addMachine}>
            Add Machine
          </button>
        </div>
        <div className="machine-list">
          {machines.size == 0 ? (
            <div className="no-machines-message">
              No machines connected. Click "Add Machine" to create one.
            </div>
          ) : (
            Array.from(machines.values()).map((machine) => (
              <div key={machine.id} className="machine-item-row">
                <div className="machine-info">
                  <div className="machine-id">
                    Machine # {machine.id}
                  </div>
                  <div className="machine-location">
                    <span>Lat: {machine.location.lat.toFixed(6)}</span>
                    <span>Lon: {machine.location.lon.toFixed(6)}</span>
                    <span>Alt: {machine.location.alt.toFixed(1)}m</span>
                  </div>
                  <div className="fuel-level">
                    Fuel: {machine.fuel_level.toFixed(1)}%
                  </div>
                  <div className="button-container">
                    <button
                      className={machine.is_paused ? "unpause-btn" : "pause-btn"}
                      onClick={() => togglePause(machine)}
                    >
                      {machine.is_paused ? 'Resume' : 'Pause'}
                    </button>
                    <button
                      className="remove-btn"
                      onClick={() => removeMachine(machine.id)}
                    >
                      ✕
                    </button>
                  </div>
                </div>
              </div>
            ))
          )}
        </div>
      </div>
    </div>
  );
}

export default App;
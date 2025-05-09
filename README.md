# Machine Map Monitor

A real-time monitoring application for tracking machines with gRPC streams, WebSockets, and Google Maps integration.

## Architecture

The application consists of three main components:

1. **gRPC Server**: Written in Go, manages machines and their state using Brownian motion simulation
2. **WebSocket Proxy**: Go server that bridges between gRPC streams and WebSockets for browser communication
3. **Frontend**: React application with Google Maps integration for visualizing and controlling machines

## Quick Start with Docker

The easiest way to run this application is using Docker Compose:

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)
- Google Maps API Key

### Running the Application

1. Clone the repository:
   ```
   git clone https://github.com/starfishpanda/stream-machine-map-monitor.git
   cd machine-map-monitor
   ```

2. Create a `.env` file in the root directory with your Google Maps API key:
   ```
   GOOGLE_MAPS_API_KEY=your_google_maps_api_key_here
   ```

3. Build and start all services:
   ```
   docker-compose up -d
   ```

4. Open your browser and navigate to:
   ```
   http://localhost
   ```

5. To stop the application:
   ```
   docker-compose down
   ```

## Manual Setup

If you prefer to run the components individually:

### gRPC Server

```bash
cd server
go run main.go
```

### WebSocket Proxy

```bash
cd server/ws-proxy
go run main.go
```

### Frontend

```bash
cd frontend
npm install
npm run dev
```

## Features

- Real-time tracking of machine positions on Google Maps
- Control panel for managing multiple machines
- Pause/resume machine movement functionality
- Fuel level monitoring
- Brownian motion simulation for realistic movement patterns

## Usage

1. Click "Add Machine" to create a new machine on the map
2. Use the pause/resume buttons to control machine movement
3. Click on map markers to view detailed machine information
4. Monitor fuel levels as machines move around

## Troubleshooting

If you encounter any issues:

- Ensure all ports (50051, 3001, and 80) are available
- Check Docker logs with `docker-compose logs`
- Verify your Google Maps API key has the necessary permissions

## Project Structure

```
├── frontend/              # React frontend application
├── server/                # Go gRPC server
│   ├── main.go            # Server implementation
│   └── ws-proxy/          # WebSocket proxy service
│       └── main.go        # Proxy implementation
├── proto/                 # Protocol Buffer definitions
├── docker-compose.yml     # Docker Compose configuration
└── README.md              # This file
```
# Machine Map Monitor

A real-time monitoring application for tracking machines with gRPC streams, WebSockets, and Google Maps integration.

## Architecture

The application consists of three main components:

1. **gRPC Server**: Written in Go, manages machines and their state using Brownian motion simulation
2. **WebSocket Proxy**: Go server that bridges between gRPC streams and WebSockets for browser communication
3. **Frontend**: React application with Google Maps integration for visualizing and controlling machines

![Architecture Diagram](./assets/stream-machine-diagram-1.png)

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)
- Google Maps API Key

## Features

- Real-time tracking of machine positions on Google Maps
- Control panel for managing multiple machines
- Pause/resume machine movement functionality
- Fuel level monitoring
- Brownian motion simulation for realistic movement patterns

![Features](./assets/stream-machine-mock-1.png)

## Usage

1. Click "Add Machine" to create a new machine on the map
2. Use the pause/resume buttons to control machine movement
3. Click on map markers to view detailed machine information
4. Monitor fuel levels as machines move around

![Demo](./assets/demo.gif)

## Quick Start with Docker

The easiest way to run this application is using Docker Compose:

## Setup and Running Instructions

1. **Download and Extract the Repository**
   - Download the ZIP file from the GitHub repository
   - Extract it to a directory on your local machine

2. **Configure Environment Variables**
   - Copy the example environment file to create your own:
     ```bash
     cp .env.example .env
     ```
   - Edit the `.env` file and add your Google Maps API key:
     ```
     GOOGLE_MAPS_API_KEY=your_google_maps_api_key_here
     ```
     (Note: Do not use quotes around the API key)

3. **Build and Start the Application**
   - Make sure you're in the project root directory
   - Build the Docker containers:
     ```bash
     make build
     ```
   - Start the application:
     ```bash
     make start
     ```

4. **Access the Application**
   - Open your browser and navigate to:
     ```
     http://localhost
     ```
   - You should see the map interface
   - Click "Add Machine" to create and monitor machines

5. **Usage**
   - Each machine will appear on the map with a colored marker, in the paused state
   - Green markers indicate active machines, red markers indicate paused machines
   - Click on a machine marker to view details and control options
   - Use the dashboard to manage multiple machines

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

## Troubleshooting

If you encounter any issues:

- Ensure all ports (50051, 3001, and 80) are available
- Check Docker logs with `docker-compose logs`
- Verify your Google Maps API key has the necessary permissions

## Project Structure

```
├── frontend/ # React frontend application
│   ├── dist/
│   ├── src/ # Source code
│   │   ├── App.css
│   │   ├── App.tsx
│   │   └── main.tsx
│   ├── .env.development
│   ├── .gitignore
│   ├── config.ts
│   ├── Dockerfile
│   ├── eslint.config.js
│   ├── index.html
│   ├── nginx.conf
│   ├── package-lock.json
│   ├── package.json
│   ├── tsconfig.json
│   ├── vite-env.d.ts
│   └── vite.config.js
├── server/ # Backend
│   ├── proto/ # Protocol Buffer definitions
│   │   ├── machine_stream_grpc.pb.go
│   │   ├── machine_stream.pb.go
│   │   └── machine_stream.proto
│   ├── ws-proxy/ # WebSocket proxy service
│   │   ├── Dockerfile
│   │   └── main.go # Proxy implementation
│   ├── go.mod # Go gRPC ierver implementation
│   ├── go.sum
│   ├── main.go
│   └── .env.example
├── .gitignore
├── docker-compose.yml # Docker Compose configuration
├── Makefile
└── README.md # This file
```

## Possible Extensions
1. **Remove Machines**
   - The ability to remove machines would require selectively closing the specific machine's WebSocket connection
2. **Pause All Machines**
   - Useful in emergency situations
3. **Routine to Refuel Machines**
   - Ability to route machines to specific refueling location, or home base
4. **Comprehensive Alerts and Monitoring**
   - Geo-fenced and fuel-related notifications, job completion status and estimated completion time
5. **Route Planning and Visualization**
   - Sketch out a route on the map, which the vehicle will follow
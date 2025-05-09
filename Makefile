.PHONY: build start stop restart logs clean help

# Variables
CONTAINER_PREFIX = stream-machine-map

# Default target
help:
	@echo "Stream Machine Map Monitor - Docker Management"
	@echo ""
	@echo "Usage:"
	@echo "  make build    - Build all Docker images"
	@echo "  make start    - Start all containers"
	@echo "  make stop     - Stop all containers"
	@echo "  make restart  - Restart all containers"
	@echo "  make logs     - View logs from all containers"
	@echo "  make clean    - Stop and remove containers, networks, and images"
	@echo ""

# Build Docker images
build:
	@echo "Building Docker images..."
	docker-compose build

# Start containers
start:
	@echo "Starting containers..."
	docker-compose up -d
	@echo "Application is now running at http://localhost"

# Stop containers
stop:
	@echo "Stopping containers..."
	docker-compose down

# Restart containers
restart: stop start

# View logs
logs:
	@echo "Showing logs..."
	docker-compose logs -f

# Clean up
clean:
	@echo "Cleaning up Docker resources..."
	docker-compose down --rmi all --volumes --remove-orphans
	@echo "Cleanup complete."
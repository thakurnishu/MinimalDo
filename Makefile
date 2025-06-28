.PHONY: help build up down logs clean dev-backend dev-frontend

# Default target
help:
	@echo "Available commands:"
	@echo "  build        - Build all Docker images"
	@echo "  up           - Start all services"
	@echo "  down         - Stop all services"
	@echo "  logs         - Show logs from all services"
	@echo "  clean        - Clean up containers and volumes"
	@echo "  dev-backend  - Run backend in development mode"
	@echo "  dev-frontend - Serve frontend for development"

# Build all Docker images
build:
	docker compose build

# Start all services
up:
	docker compose up -d
	@echo "Services started!"
	@echo "Frontend: http://localhost:3000"
	@echo "Backend API: http://localhost:8080"
	@echo "PostgreSQL: localhost:5432"

# Stop all services
down:
	docker compose down

# Show logs
logs:
	docker compose logs -f

# Clean up everything
clean:
	docker compose down -v --remove-orphans
	docker rmi minimaldo-frontend minimaldo-backend

# Development mode for backend
dev-backend:
	@echo "Starting backend in development mode..."
	@echo "Make sure PostgreSQL is running on localhost:5432"
	cd backend && go mod tidy && go run main.go

# Development mode for frontend (simple HTTP server)
dev-frontend:
	@echo "Starting frontend development server..."
	@echo "Frontend will be available at http://localhost:8000"
	cd frontend && python3 -m http.server 8000

# Run tests (placeholder for future tests)
test:
	@echo "Running tests..."
	cd backend && go test ./...

# Install dependencies
deps:
	@echo "Installing backend dependencies..."
	cd backend && go mod tidy
	@echo "Dependencies installed!"

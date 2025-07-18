.PHONY: help build up down logs clean dev-backend dev-frontend

# Default target
help:
	@echo "Available commands:"
	@echo "  build        - Build all Docker images"
	@echo "  up           - Start all services"
	@echo "  down         - Stop all services"
	@echo "  logs         - Show logs from all services"
	@echo "  clean        - Clean up containers and volumes"

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

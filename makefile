.PHONY: up down restart logs ps help cleanup server

# Start services in the background
up:
	@echo "Starting services..."
	@docker-compose --env-file ./cred.env up -d

# Stop services
down:
	@echo "Stopping services..."
	@docker-compose --env-file ./cred.env down

# Restart services
restart: down up
	@echo "Services restarted."

# Show logs of services
logs:
	@docker-compose --env-file ./cred.env logs -f

# Show status of services
ps:
	@docker-compose --env-file ./cred.env ps

cleanup:
	@echo "Cleaning up stopped containers, networks, volumes, and dangling images..."
	@docker system prune -f

server:
	@echo "Starting Go server..."
	@go run app-backend/main.go

# Show help
help:
	@echo "-------------------------------------------"
	@echo "Makefile commands:"
	@echo "-------------------------------------------"
	@echo "make up       - Start services"
	@echo "make down     - Stop services"
	@echo "make restart  - Restart services"
	@echo "make logs     - Show logs of services"
	@echo "make ps       - Show status of services"
	@echo "make cleanup  - Clean up Docker (remove stopped containers, networks, volumes, and dangling images)"
	@echo "make server   - Start Go server located under app/main.go"
	@echo "make help     - Show help message"
	@echo "-------------------------------------------"

.DEFAULT_GOAL := help

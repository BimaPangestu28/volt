# Volt API Client - Development Makefile

.PHONY: help install dev build test clean docker-up docker-down docker-build docker-logs

# Default target
help: ## Show this help message
	@echo "⚡ Volt API Client - Development Commands"
	@echo ""
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Installation
install: ## Install all dependencies
	@echo "📦 Installing dependencies..."
	@echo "Installing backend dependencies..."
	cd backend && go mod tidy
	@echo "Installing frontend dependencies..."
	cd frontend && pnpm install --config.engine-strict=false
	@echo "✅ All dependencies installed!"

# Development
dev: ## Start development servers (backend + frontend)
	@echo "🚀 Starting development servers..."
	@make -j2 dev-backend dev-frontend

dev-backend: ## Start backend development server with hot reload
	@echo "🔧 Starting backend server with hot reload..."
	cd backend && cp -n .env.example .env 2>/dev/null || true && air

dev-backend-simple: ## Start backend development server (simple)
	@echo "🔧 Starting backend server..."
	cd backend && cp -n .env.example .env 2>/dev/null || true && go run cmd/server/main.go

dev-frontend: ## Start frontend development server
	@echo "🎨 Starting frontend server..."
	cd frontend && pnpm run dev -- --host 0.0.0.0 --port 5173

dev-backend-only: ## Start only backend server
	@make dev-backend

dev-frontend-only: ## Start only frontend server
	@make dev-frontend

# Building
build: ## Build both frontend and backend
	@echo "🏗️ Building applications..."
	@make build-backend
	@make build-frontend

build-backend: ## Build backend binary
	@echo "Building backend..."
	cd backend && go build -o bin/volt-api cmd/server/main.go
	@echo "✅ Backend built to backend/bin/volt-api"

build-frontend: ## Build frontend for production
	@echo "Building frontend..."
	cd frontend && pnpm run build
	@echo "✅ Frontend built to frontend/build/"

# Testing
test: ## Run all tests
	@echo "🧪 Running tests..."
	@make test-backend
	@make test-frontend

test-backend: ## Run backend tests
	@echo "Testing backend..."
	cd backend && go test ./...
	cd backend && go test -race ./...

test-frontend: ## Run frontend tests
	@echo "Testing frontend..."
	cd frontend && pnpm run check
	cd frontend && pnpm run type-check

# Code Quality
lint: ## Run linters
	@echo "🔍 Running linters..."
	cd frontend && pnpm run lint
	@echo "✅ Lint complete!"

format: ## Format code
	@echo "💅 Formatting code..."
	cd frontend && pnpm run format
	cd backend && go fmt ./...
	@echo "✅ Code formatted!"

# Docker Development
docker-up: ## Start all services with Docker Compose
	@echo "🐳 Starting Docker development environment..."
	docker-compose up -d
	@echo "✅ Services started! Frontend: http://localhost:5173, Backend: http://localhost:8080"

docker-down: ## Stop all Docker services
	@echo "🛑 Stopping Docker services..."
	docker-compose down
	@echo "✅ Services stopped!"

docker-build: ## Build Docker images
	@echo "🔨 Building Docker images..."
	docker-compose build
	@echo "✅ Images built!"

docker-logs: ## View Docker logs
	@echo "📋 Viewing Docker logs..."
	docker-compose logs -f

docker-logs-backend: ## View backend logs
	docker-compose logs -f backend

docker-logs-frontend: ## View frontend logs
	docker-compose logs -f frontend

# Database
db-setup: ## Setup MongoDB with Docker
	@echo "🗄️ Setting up MongoDB..."
	docker-compose up -d mongodb
	@echo "✅ MongoDB started on localhost:27017"

db-shell: ## Open MongoDB shell
	@echo "🐚 Opening MongoDB shell..."
	docker-compose exec mongodb mongosh volt

# Environment
env-setup: ## Setup environment files
	@echo "⚙️ Setting up environment files..."
	cd backend && cp -n .env.example .env 2>/dev/null || true
	cd frontend && cp -n .env.example .env 2>/dev/null || true
	@echo "✅ Environment files created! Please edit them with your configuration."

# Cleanup
clean: ## Clean build artifacts and dependencies
	@echo "🧹 Cleaning up..."
	rm -rf backend/bin/
	rm -rf frontend/build/
	rm -rf frontend/node_modules/
	rm -rf frontend/.svelte-kit/
	cd backend && go clean -cache
	@echo "✅ Cleanup complete!"

clean-docker: ## Clean Docker containers and images
	@echo "🧹 Cleaning Docker..."
	docker-compose down -v
	docker system prune -f
	@echo "✅ Docker cleanup complete!"

# Health checks
health: ## Check if services are running
	@echo "🩺 Checking service health..."
	@curl -s http://localhost:8080/health && echo "✅ Backend is healthy" || echo "❌ Backend is down"
	@curl -s http://localhost:5173 > /dev/null && echo "✅ Frontend is healthy" || echo "❌ Frontend is down"

# Quick commands
quick-start: install docker-up ## Quick start: install dependencies and start with Docker
	@echo "🚀 Volt is ready! Visit http://localhost:5173"

quick-dev: env-setup install dev ## Quick dev: setup environment, install deps, and start dev servers
	@echo "🚀 Development servers started!"

# Production
prod-build: ## Build for production
	@echo "🏭 Building for production..."
	@make build
	@echo "✅ Production build complete!"

# Information
info: ## Show project information
	@echo "⚡ Volt API Client"
	@echo "Frontend: SvelteKit + TypeScript"
	@echo "Backend: Go + Gin Framework"
	@echo "Database: MongoDB"
	@echo ""
	@echo "Development URLs:"
	@echo "  Frontend: http://localhost:5173"
	@echo "  Backend:  http://localhost:8080"
	@echo "  MongoDB:  mongodb://localhost:27017"
	@echo ""
	@echo "Run 'make help' for available commands"
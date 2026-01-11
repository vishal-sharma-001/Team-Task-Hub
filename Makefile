.PHONY: help build up down logs clean ps

help:
	@echo "🚀 Team Task Hub - Makefile Commands"
	@echo "===================================="
	@echo ""
	@echo "make build          - Build Docker images"
	@echo "make up             - Start all services (build + run)"
	@echo "make down           - Stop all services"
	@echo "make logs           - View logs from all services"
	@echo "make logs-backend   - View backend logs only"
	@echo "make logs-frontend  - View frontend logs only"
	@echo "make logs-db        - View database logs only"
	@echo "make ps             - Show running containers"
	@echo "make clean          - Remove containers, images, volumes"
	@echo ""
	@echo "Quick Start:"
	@echo "  make up"
	@echo ""
	@echo "Access:"
	@echo "  Frontend:  http://localhost:3000"
	@echo "  Backend:   http://localhost:8080/api"
	@echo "  Database:  localhost:5432"
	@echo ""

build:
	@echo "🔨 Building Docker images..."
	docker-compose build

up:
	@echo "🚀 Starting all services..."
	docker-compose up --build

down:
	@echo "🛑 Stopping all services..."
	docker-compose down

logs:
	@echo "📋 Showing all logs..."
	docker-compose logs -f

logs-backend:
	@echo "📋 Backend logs..."
	docker-compose logs -f backend

logs-frontend:
	@echo "📋 Frontend logs..."
	docker-compose logs -f frontend

logs-db:
	@echo "📋 Database logs..."
	docker-compose logs -f postgres

ps:
	@echo "📦 Running containers:"
	docker-compose ps

clean:
	@echo "🧹 Cleaning up Docker resources..."
	docker-compose down -v
	@echo "✅ Cleanup complete"

restart:
	@echo "🔄 Restarting services..."
	docker-compose restart

restart-backend:
	@echo "🔄 Restarting backend..."
	docker-compose restart backend

restart-frontend:
	@echo "🔄 Restarting frontend..."
	docker-compose restart frontend

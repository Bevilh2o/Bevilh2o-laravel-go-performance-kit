.PHONY: help up down restart logs ps migrate seed test-laravel bench-sync bench-async bench-go bench-all tidy-go

# Default target: display help documentation
help:
	@echo "Laravel + Go Performance Kit — Available Commands:"
	@echo ""
	@echo "  Infrastructure & Lifecycle:"
	@echo "    make up             Start all containers in detached mode"
	@echo "    make down           Stop and remove all project containers"
	@echo "    make restart        Rebuild and restart all containers"
	@echo "    make logs           Follow aggregated container logs"
	@echo "    make ps             Display status of running containers"
	@echo ""
	@echo "  Database & Migrations:"
	@echo "    make migrate        Run database migrations for Laravel"
	@echo "    make migrate-fresh  Drop all tables and re-run migrations"
	@echo ""
	@echo "  Go Service Maintenance:"
	@echo "    make tidy-go        Tidy and synchronize Go dependencies"
	@echo ""
	@echo "  Benchmarking Suite (k6):"
	@echo "    make bench-sync     Run k6 baseline benchmark against Laravel Sync (POST /api/events)"
	@echo "    make bench-async    Run k6 benchmark against Laravel Async Queue (POST /api/events/async)"
	@echo "    make bench-go       Run k6 benchmark against Go Ingestion Service (POST /events)"
	@echo "    make bench-all      Execute all benchmark scenarios sequentially"

up:
	docker compose up -d --build

down:
	docker compose down

restart:
	docker compose down && docker compose up -d --build

logs:
	docker compose logs -f

ps:
	docker compose ps

migrate:
	docker compose exec app php artisan migrate

migrate-fresh:
	docker compose exec app php artisan migrate:fresh

tidy-go:
	MSYS_NO_PATHCONV=1 docker run --rm -v "$$(pwd)/go:/app" -w /app golang:alpine go mod tidy

bench-sync:
	MSYS_NO_PATHCONV=1 docker compose run --rm k6 run /benchmarks/k6_sync.js

bench-async:
	MSYS_NO_PATHCONV=1 docker compose run --rm k6 run /benchmarks/k6_async.js

bench-go:
	MSYS_NO_PATHCONV=1 docker compose run --rm k6 run /benchmarks/k6_go.js

bench-all: bench-sync bench-async bench-go
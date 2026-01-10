# Makefile for Exam App Queue Project

.PHONY: help setup start-infra migrate-up migrate-down run-direct run-queue run-worker test-direct test-queue clean

help:
	@echo "Exam App Queue - Available Commands"
	@echo ""
	@echo "Setup:"
	@echo "  make setup          - Start infrastructure and run migrations"
	@echo "  make start-infra    - Start PostgreSQL and NATS containers"
	@echo "  make migrate-up     - Run database migrations"
	@echo "  make migrate-down   - Rollback database migrations"
	@echo ""
	@echo "Run Applications:"
	@echo "  make run-direct     - Start direct insert API (port 8080)"
	@echo "  make run-queue      - Start queue-based API (port 8081)"
	@echo "  make run-worker     - Start background workers"
	@echo ""
	@echo "Load Testing:"
	@echo "  make test-direct    - Run k6 load test on direct approach"
	@echo "  make test-queue     - Run k6 load test on queue approach"
	@echo ""
	@echo "Utilities:"
	@echo "  make logs           - View all container logs"
	@echo "  make db-shell       - Connect to PostgreSQL shell"
	@echo "  make check-db       - Check submission counts"
	@echo "  make clean          - Stop containers and remove volumes"
	@echo "  make clean-db       - Delete all submissions from database"

setup: start-infra migrate-up
	@echo "✅ Setup complete!"

start-infra:
	@echo "Starting PostgreSQL and NATS..."
	docker-compose up -d
	@echo "Waiting for services to be ready..."
	@sleep 8
	@echo "✅ Infrastructure ready"

migrate-up:
	@echo "Running database migrations..."
	cd migrate && go run migrate.go up

migrate-down:
	@echo "Rolling back migrations..."
	cd migrate && go run migrate.go down

run-direct:
	@echo "Starting Direct Insert API on :8080..."
	cd api/direct && go run main.go

run-queue:
	@echo "Starting Queue-based API on :8081..."
	cd api/queue && go run main.go

run-worker:
	@echo "Starting 10 background workers..."
	cd worker && go run main.go

test-direct:
	@echo "Running load test on Direct approach..."
	@echo "Make sure api/direct is running!"
	cd load-test && k6 run test_direct.js

test-queue:
	@echo "Running load test on Queue approach..."
	@echo "Make sure api/queue AND workers are running!"
	cd load-test && k6 run test_queue.js

logs:
	docker-compose logs -f

db-shell:
	docker exec -it exam-db psql -U postgres -d exam_app

check-db:
	docker exec -it exam-db psql -U postgres -d exam_app -c "SELECT status, COUNT(*) FROM exam_submissions GROUP BY status;"

clean-db:
	docker exec -it exam-db psql -U postgres -d exam_app -c "TRUNCATE TABLE exam_submissions;"
	@echo "✅ All submissions deleted"

clean:
	@echo "Stopping containers..."
	docker-compose down -v
	@echo "✅ Cleanup complete"

install-deps:
	go mod download
	@echo "✅ Dependencies installed"

check-health:
	@echo "Checking services..."
	@echo "PostgreSQL:"
	@docker-compose ps postgres
	@echo ""
	@echo "NATS:"
	@docker-compose ps nats
	@echo ""
	@echo "NATS Monitoring: http://localhost:8222"

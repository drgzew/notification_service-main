.PHONY: up down logs build run run-win run-docker create-topics psql test

COMPOSE_FILE=docker-compose.dev.yml
CONFIG=example/config.dev.yaml
IMAGE_NAME=notification_service:dev

up:
	docker-compose -f $(COMPOSE_FILE) up -d --build

down:
	docker-compose -f $(COMPOSE_FILE) down

logs:
	docker-compose -f $(COMPOSE_FILE) logs -f

build:
	docker build -t $(IMAGE_NAME) .

run:
	@echo "Run locally (Linux/macOS) with CONFIG=$(CONFIG)"
	@CONFIG=$(CONFIG) env configPath=$(CONFIG) go run ./cmd/app

run-win:
	@echo "Run locally (Windows PowerShell). Set configPath env before running:"
	@echo "  $env:configPath = '$(CURDIR)/$(CONFIG)'"

run-docker:
	docker-compose -f $(COMPOSE_FILE) up --build app

create-topics:
	# Create Kafka topics for local dev
	docker-compose -f $(COMPOSE_FILE) exec kafka kafka-topics.sh --create --topic event_notifications --bootstrap-server kafka:9092 --replication-factor 1 --partitions 1 || true
	docker-compose -f $(COMPOSE_FILE) exec kafka kafka-topics.sh --create --topic notification_status --bootstrap-server kafka:9092 --replication-factor 1 --partitions 1 || true

psql:
	docker-compose -f $(COMPOSE_FILE) exec postgres psql -U postgres -d notifications

test:
	go test ./...

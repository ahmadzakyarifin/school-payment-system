# Load .env if exists
ifneq (,$(wildcard .env))
	include .env
	export
endif

# Database URL
DB_URL := mysql://$(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)
# Migration tool
MIGRATE := $(HOME)/go/bin/migrate -path migrations -database "$(DB_URL)"

.DEFAULT_GOAL := help

.PHONY: help
help:
	@echo ""
	@echo "Migration Commands:"
	@echo "  make up              - Run all pending migrations"
	@echo "  make down            - Rollback last migration"
	@echo "  make reset           - Drop & re-run all migrations"
	@echo "  make create          - Create new migration"
	@echo "  make force           - Force database to specific version"
	@echo "  make version         - Show current database migration version"
	@echo ""
	@echo "App Commands:"
	@echo "  make dev             - Run with hot reload (air)"
	@echo "  make seed            - Create default admin account"
	@echo ""


.PHONY: up
up:
	@echo "Running migrations..."
	$(MIGRATE) up

.PHONY: down
down:
	@echo "Rolling back last migration..."
	$(MIGRATE) down 1

.PHONY: reset
reset:
	@echo "Resetting database..."
	$(MIGRATE) drop -f
	$(MIGRATE) up

.PHONY: create
create:
	@read -p "Migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

.PHONY: force
force:
	@if [ -z "$(version)" ]; then \
		read -p "Force database to version: " ver; \
	else \
		ver=$(version); \
	fi; \
	echo "Forcing database to version $$ver ..."; \
	$(MIGRATE) force $$ver

.PHONY: version
version:
	@echo "Checking current migration version..."
	$(MIGRATE) version

# ============================================================
# APP
# ============================================================

.PHONY: dev
dev:
	air

.PHONY: seed
seed:
	go run cmd/seed/main.go

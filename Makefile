# Makefile

.PHONY: help build run clean db-start db-stop db-up db-down db-reset

BIN_DIR := .bin
BINARY := flashcards
BUILD_OUT := $(BIN_DIR)/$(BINARY)

help:
	@echo "Available commands:"
	@echo "  build     - Build the application into $(BIN_DIR)"
	@echo "  run       - Run the application"
	@echo "  clean     - Clean build artifacts"
	@echo "  db-start  - Start Supabase local development"
	@echo "  db-stop   - Stop Supabase local development"
	@echo "  db-up     - Run database migrations"
	@echo "  db-down   - Rollback database migrations"
	@echo "  db-reset  - Reset database (stop, start, migrate)"

build:
	-@mkdir -p $(BIN_DIR) 2>/dev/null || mkdir $(BIN_DIR) 2>nul || true
	go build -o $(BUILD_OUT) cmd/main.go
	@echo "Built $(BUILD_OUT)"

run:
	go run cmd/main.go

clean:
	-@rm -rf $(BIN_DIR) 2>/dev/null || rmdir /s /q $(BIN_DIR) 2>nul || true
	@echo "Removed $(BIN_DIR)"

db-start:
	@echo "Starting Supabase local development..."
	supabase start

db-stop:
	@echo "Stopping Supabase local development..."
	supabase stop

db-up:
	@echo "Running database migrations..."
	supabase migration up

db-down:
	@echo "Rolling back database migrations..."
	supabase migration down

db-reset: db-stop db-start db-up
	@echo "Database reset complete"

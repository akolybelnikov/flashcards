# Makefile

BIN_DIR := .bin

.PHONY: help build run clean db-start db-stop db-up db-down db-reset db-link db-status db-push db-pull db-remote-reset db-remote-reset-confirm db-new db-list test clean-test generate clean-mocks

help:

BIN_DIR := .bin
BINARY := flashcards
BUILD_OUT := $(BIN_DIR)/$(BINARY)

help:
	@echo "Available commands:"
	@echo ""
	@echo "Build & Run:"
	@echo "  build       - Build the application into $(BIN_DIR)"
	@echo "  run         - Run the application"
	@echo "  test        - Run all tests"
	@echo "  clean       - Clean build artifacts"
	@echo "  clean-test  - Clear Go test cache"
	@echo ""
	@echo "Code Generation:"
	@echo "  generate    - Generate mocks using mockgen"
	@echo "  clean-mocks - Delete all generated mocks"
	@echo ""
	@echo "Local Database (Docker):"
	@echo "  db-start    - Start Supabase local development"
	@echo "  db-stop     - Stop Supabase local development"
	@echo "  db-up       - Run database migrations (local)"
	@echo "  db-down     - Rollback database migrations (local)"
	@echo "  db-reset    - Reset local database (stop, start, migrate)"
	@echo ""
	@echo "Remote Database (Production):"
	@echo "  db-link          - Link to remote Supabase project"
	@echo "  db-status        - Check migration status (local vs remote)"
	@echo "  db-push          - Push migrations to remote Supabase"
	@echo "  db-pull          - Pull remote schema to local migrations"
	@echo "  db-remote-reset  - Reset remote database and re-apply all migrations (DESTRUCTIVE!)"
	@echo ""
	@echo "Migration Management:"
	@echo "  db-new      - Create a new migration file"
	@echo "  db-list     - List all migrations"

build:
	-@mkdir -p $(BIN_DIR) 2>/dev/null || mkdir $(BIN_DIR) 2>nul || true
	go build -o $(BUILD_OUT) cmd/main.go
	@echo "Built $(BUILD_OUT)"

run:
	go run cmd/main.go

test:
	go test ./... -v

clean:
	-@rm -rf $(BIN_DIR) 2>/dev/null || rmdir /s /q $(BIN_DIR) 2>nul || true
	@echo "Removed $(BIN_DIR)"

clean-test:
	go clean -testcache
	@echo "Cleared Go test cache"

# Local database commands
db-start:
	@echo "Starting Supabase local development..."
	supabase start

db-stop:
	@echo "Stopping Supabase local development..."
	supabase stop

db-up:
	@echo "Running database migrations (local)..."
	supabase migration up

db-down:
	@echo "Rolling back database migrations (local)..."
	supabase migration down

db-reset: db-stop db-start db-up
	@echo "Database reset complete"

# Remote database commands
db-link:
	@echo "Linking to remote Supabase project..."
	@echo "You'll need your project ref (found in Supabase dashboard URL)"
	supabase link

db-status:
	@echo "Checking migration status (local vs remote)..."
	supabase migration list

db-push:
	@echo "Pushing migrations to remote Supabase..."
	@echo "WARNING: This will apply migrations to production!"
	supabase db push

db-pull:
	@echo "Pulling remote schema to local migrations..."
	supabase db pull

db-remote-reset-confirm:
	@echo "========================================"
	@echo "REMOTE DATABASE RESET"
	@echo "========================================"
	@echo "This will:"
	@echo "  1. Drop all tables in remote database"
	@echo "  2. Re-apply all migrations from scratch"
	@echo "  3. All data will be LOST!"
	@echo ""
	@echo "WARNING: This affects PRODUCTION database!"
	@echo "========================================"
	@echo ""
	@echo "Press Ctrl+C now to cancel, or"
	@pause
	supabase db reset --linked

db-remote-reset:
	@echo "Resetting remote database without confirmation..."
	supabase db reset --linked

# Migration management
db-new:
	@echo "Creating new migration file..."
	@echo "Usage: supabase migration new <migration_name>"
	@echo "Example: supabase migration new add_user_table"

db-list:
	@echo "Listing all migrations..."
	supabase migration list

# Code generation
generate:
	@echo "Generating mocks..."
	go generate ./...
	@echo "✓ Mocks generated successfully"

clean-mocks:
	@echo "Cleaning generated mocks..."	@rm -rf services/mocks db/mocks 2>/dev/null || true
	@echo "✓ Mocks cleaned"


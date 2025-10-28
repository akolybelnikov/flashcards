#!/usr/bin/env bash
# exit on error
set -o errexit

# Install dependencies
go mod download
go mod verify

# Run tests (optional - comment out if you want faster builds)
# go test ./... -v

# Build the application from cmd/main.go
go build -tags netgo -ldflags '-s -w' -o app ./cmd

# Make the binary executable
chmod +x app

echo "Build completed successfully!"

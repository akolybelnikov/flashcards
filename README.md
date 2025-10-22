# Go Project Template

A comprehensive Go project template with REST API functionality, PostgreSQL database integration, and Supabase local development setup. This template is designed to be used throughout the course as a starting point for building Go applications.

## Prerequisites

- [Go 1.24.1+](https://golang.org/dl/)
- [Docker](https://www.docker.com/get-started)
- [Supabase CLI](https://supabase.com/docs/guides/cli)

## Quick Start

### Environment Configuration

The project requires a `.env` file for local development. This file is already provided when you clone the repository and contains the necessary environment variables:

- `DB_URL` - PostgreSQL database connection string
- `PORT` - Application port (defaults to 8080)

### Database Setup

Make sure Docker is running, then use the provided Makefile commands:

```bash
# Start Supabase local development environment
make db-start

# Run database migrations
make db-up
```

### Run the Application

```bash
# Run directly
make run

# Or build and run
make build
./todo-api
```

The application will start on `http://localhost:8080` (or the port specified in your `.env` file).

## Available Commands

### Application Commands
- `make build` - Build the application binary
- `make run` - Run the application directly
- `make clean` - Clean build artifacts

### Database Commands
- `make db-start` - Start Supabase local development
- `make db-stop` - Stop Supabase local development  
- `make db-up` - Run database migrations

## API Endpoints

The template includes a complete REST API with the following endpoints:

### Health Check
- `GET /health` - Application health status

### Exported calls for REST client
You can find an exported HAR archive which you can import into a REST client for easily interacting with the API in `./artifacts`

## Configuration

The application uses environment-based configuration managed through the `config` package. Key configuration options:

- **DB_URL**: PostgreSQL database connection string (required)
- **PORT**: Application port (optional, defaults to 8080)

## Database

The project uses PostgreSQL with Supabase for local development:

- **Local Database**: Accessible at `localhost:54322`
- **Supabase Studio**: Available at `http://localhost:54323`
- **API**: Available at `http://localhost:54321`

### Migration Management

Database schema is managed through SQL migrations located in: `supabase/migrations/`.

## Creating Your Own Project

When you're ready to build your own application using this template, you can delete the existing todo API implementation and replace it with your own business logic. The template provides the foundation with database connectivity, configuration management, and API structure.

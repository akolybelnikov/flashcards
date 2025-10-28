# Go Project (AI app demo)

A comprehensive Go project with 
REST API functionality, 
PostgresQL database integration, 
and Supabase local development setup. 

## Prerequisites

- [Go 1.25.0+](https://golang.org/dl/)
- [Docker](https://www.docker.com/get-started)
- [Supabase CLI](https://supabase.com/docs/guides/cli)

## Quick Start

### Environment Configuration

The project requires a `.env` file for local development. This file is already provided when you clone the repository and contains the necessary environment variables:

- `DB_URL` - PostgresQL database connection string
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

## Configuration

The application uses environment-based configuration managed through the `config` package. Key configuration options:

- **DB_URL**: PostgresQL database connection string (required)
- **PORT**: Application port (optional, defaults to 8080)

## Database

The project uses PostgresQL with Supabase for local development:

- **Local Database**: Accessible at `localhost:54322`
- **Supabase Studio**: Available at `http://localhost:54323`
- **API**: Available at `http://localhost:54321`

### Migration Management

Database schema is managed through SQL migrations located in: `supabase/migrations/`.

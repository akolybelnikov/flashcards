# Flashcards API

AI-engineered flashcard application for language learning with intelligent translation features.

A comprehensive Go REST API with PostgresQL database integration, OpenAI-powered translations, and Supabase local development setup. 

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
./.bin/flashcards
```

The application will start on `http://localhost:10000` (or the port specified in your `.env` file).

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

Full API documentation is available in the [OpenAPI specification](./openapi.yaml).

The API includes the following endpoints:

### Flashcards
- `POST /flashcards` - Create a new flashcard (with optional AI translation)
- `GET /flashcards` - Get all flashcards
- `GET /flashcards/{id}` - Get a specific flashcard by ID
- `PUT /flashcards/{id}` - Update a flashcard
- `DELETE /flashcards/{id}` - Delete a flashcard
- `GET /flashcards/random` - Get a random flashcard for study (with an optional AI hint)

### Health Check
- `GET /health` - Application health status

### API Features

- **AI Translation**: Automatically translate flashcards between English and Greek
- **Validation**: Smart validation ensures language parameters are provided when needed
- **Study Mode**: Random flashcard endpoint for practicing
- **Full CRUD**: Complete create, read, update, delete operations


## Configuration

The application uses environment-based configuration managed through the `config` package. Key configuration options:

- **DB_URL**: PostgresQL database connection string (required)
- **PORT**: Application port (optional, defaults to 8080)
- **OPENAI_API_KEY**: OpenAI API key for AI translation features (optional)

## Database

The project uses PostgresQL with Supabase for local development:

- **Local Database**: Accessible at `localhost:54322`
- **Supabase Studio**: Available at `http://localhost:54323`
- **API**: Available at `http://localhost:54321`

### Migration Management

Database schema is managed through SQL migrations located in: `supabase/migrations/`.

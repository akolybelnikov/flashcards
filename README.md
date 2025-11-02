# Flashcards API

AI-engineered flashcard application for language learning with intelligent translation features.

A comprehensive Go REST API with PostgreSQL database integration, OpenAI-powered translations with caching, and Supabase local development setup.

## Features

- 🤖 **AI Translation Generation** - Generate translations on-demand with 1-hour caching
- 📝 **Manual Card Creation** - Full control over flashcard content
- 🎯 **Translation Tracking** - Know which fields were AI-generated
- 🌍 **Multi-language Support** - Any ISO 639-1 language pair
- 📚 **Study Mode** - Random flashcard selection with optional AI hints
- 🔄 **Full CRUD** - Complete REST API for flashcard management 

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

Full API documentation is available in the [OpenAPI specification](./openapi.yaml) (v2.0.0).

### Translation (New in v2.0)
- `POST /flashcards/translate` - Generate AI translation with caching

### Flashcards
- `POST /flashcards` - Create a new flashcard (both question and answer required)
- `GET /flashcards` - Get all flashcards
- `GET /flashcards/{id}` - Get a specific flashcard by ID
- `PUT /flashcards/{id}` - Update a flashcard
- `DELETE /flashcards/{id}` - Delete a flashcard
- `GET /flashcards/random` - Get a random flashcard with optional AI hint

### Health Check
- `GET /health` - Application health status

### API Workflow (v2.0)

**Creating a Flashcard with AI Translation:**

1. **Generate Translation:**
   ```bash
   POST /flashcards/translate
   {
     "content": "hello",
     "from_lang": "en",
     "to_lang": "el"
   }
   # Returns: {"translation": "γεια σας", "cached": false}
   ```

2. **Create Flashcard:**
   ```bash
   POST /flashcards
   {
     "question": "hello",
     "answer": "γεια σας",
     "question_lang": "en",
     "answer_lang": "el",
     "ai_translated_answer": true
   }
   ```

See [Testing Guide](.ai/TESTING_GUIDE.md) for more examples.

### API Features

- **Translation Caching**: Generated translations are cached for 1 hour
- **AI Tracking**: Track which fields were AI-generated vs manually entered
- **Multi-language**: Support for any ISO 639-1 language pair (en, el, fr, de, es, etc.)
- **Validation**: Smart input validation with helpful error messages
- **Study Mode**: Random flashcard selection for practice sessions


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

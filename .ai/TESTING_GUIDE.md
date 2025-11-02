# Testing the Stateful Translation Feature

## Overview
This document provides curl commands to test the new stateful translation workflow.

## Prerequisites
- Server running on http://localhost:8080
- Database migrated with new columns
- OpenAI API key configured (for actual translation)

## Test Workflow

### 1. Generate Translation (First Call - Cache Miss)

Generate a translation from English to Greek:

```bash
curl -X POST http://localhost:8080/flashcards/translate \
  -H "Content-Type: application/json" \
  -d '{
    "content": "hello",
    "from_lang": "en",
    "to_lang": "el"
  }'
```

**Expected Response:**
```json
{
  "translation": "γεια σας",
  "cached": false,
  "cache_key": "abc123..."
}
```

### 2. Generate Same Translation Again (Cache Hit)

Make the same request again - should return cached result:

```bash
curl -X POST http://localhost:8080/flashcards/translate \
  -H "Content-Type: application/json" \
  -d '{
    "content": "hello",
    "from_lang": "en",
    "to_lang": "el"
  }'
```

**Expected Response:**
```json
{
  "translation": "γεια σας",
  "cached": true,
  "cache_key": "abc123..."
}
```

Note: `cached` should now be `true`

### 3. Create Flashcard with AI Translation

Now create a flashcard using the translation:

```bash
curl -X POST http://localhost:8080/flashcards \
  -H "Content-Type: application/json" \
  -d '{
    "question": "hello",
    "answer": "γεια σας",
    "question_lang": "en",
    "answer_lang": "el",
    "ai_translated_answer": true
  }'
```

**Expected Response:**
```json
{
  "id": 1,
  "question": "hello",
  "answer": "γεια σας",
  "question_lang": "en",
  "answer_lang": "el",
  "ai_translated_question": false,
  "ai_translated_answer": true,
  "created_at": "2025-11-02T...",
  "updated_at": "2025-11-02T..."
}
```

### 4. Create Flashcard Without AI (Manual Entry)

Create a flashcard where user typed both sides manually:

```bash
curl -X POST http://localhost:8080/flashcards \
  -H "Content-Type: application/json" \
  -d '{
    "question": "goodbye",
    "answer": "αντίο",
    "question_lang": "en",
    "answer_lang": "el"
  }'
```

**Expected Response:**
```json
{
  "id": 2,
  "question": "goodbye",
  "answer": "αντίο",
  "question_lang": "en",
  "answer_lang": "el",
  "ai_translated_question": false,
  "ai_translated_answer": false,
  "created_at": "2025-11-02T...",
  "updated_at": "2025-11-02T..."
}
```

### 5. Get All Flashcards (Verify AI Flags)

Retrieve all flashcards to see AI flags:

```bash
curl http://localhost:8080/flashcards
```

**Expected Response:**
```json
[
  {
    "id": 2,
    "question": "goodbye",
    "answer": "αντίο",
    "question_lang": "en",
    "answer_lang": "el",
    "ai_translated_question": false,
    "ai_translated_answer": false,
    "created_at": "...",
    "updated_at": "..."
  },
  {
    "id": 1,
    "question": "hello",
    "answer": "γεια σας",
    "question_lang": "en",
    "answer_lang": "el",
    "ai_translated_question": false,
    "ai_translated_answer": true,
    "created_at": "...",
    "updated_at": "..."
  }
]
```

### 6. Get Flashcard by ID

```bash
curl http://localhost:8080/flashcards/1
```

### 7. Update Flashcard (Manual Edit)

```bash
curl -X PUT http://localhost:8080/flashcards/1 \
  -H "Content-Type: application/json" \
  -d '{
    "answer": "γειά σου"
  }'
```

Note: AI flags remain unchanged when updating

### 8. Delete Flashcard

```bash
curl -X DELETE http://localhost:8080/flashcards/2
```

### 9. Test Random Flashcard with AI Hint

```bash
curl "http://localhost:8080/flashcards/random?lang=el"
```

## Error Cases

### Missing Content

```bash
curl -X POST http://localhost:8080/flashcards/translate \
  -H "Content-Type: application/json" \
  -d '{
    "content": "",
    "from_lang": "en",
    "to_lang": "el"
  }'
```

**Expected:** 400 Bad Request with error: "Content cannot be empty"

### Missing Languages

```bash
curl -X POST http://localhost:8080/flashcards/translate \
  -H "Content-Type: application/json" \
  -d '{
    "content": "hello",
    "from_lang": "",
    "to_lang": "el"
  }'
```

**Expected:** 400 Bad Request with error: "Both from_lang and to_lang must be provided"

### Create Flashcard with Empty Fields

```bash
curl -X POST http://localhost:8080/flashcards \
  -H "Content-Type: application/json" \
  -d '{
    "question": "",
    "answer": "test"
  }'
```

**Expected:** 400 Bad Request with error: "Both question and answer must be provided"

## Cache Behavior Testing

### Test Different Language Pairs

```bash
# English to French
curl -X POST http://localhost:8080/flashcards/translate \
  -H "Content-Type: application/json" \
  -d '{
    "content": "hello",
    "from_lang": "en",
    "to_lang": "fr"
  }'

# French to English
curl -X POST http://localhost:8080/flashcards/translate \
  -H "Content-Type: application/json" \
  -d '{
    "content": "bonjour",
    "from_lang": "fr",
    "to_lang": "en"
  }'
```

Each language pair should have its own cache entry.

### Test Case Sensitivity

```bash
# Lowercase
curl -X POST http://localhost:8080/flashcards/translate \
  -H "Content-Type: application/json" \
  -d '{
    "content": "hello",
    "from_lang": "en",
    "to_lang": "el"
  }'

# Uppercase (different cache key)
curl -X POST http://localhost:8080/flashcards/translate \
  -H "Content-Type: application/json" \
  -d '{
    "content": "HELLO",
    "from_lang": "en",
    "to_lang": "el"
  }'
```

Different content = different cache keys.

## Complete User Workflow Example

Simulating a user creating a flashcard with AI assistance:

```bash
# Step 1: User types "hello" in English field
# Step 2: User clicks "Generate Translation" button
curl -X POST http://localhost:8080/flashcards/translate \
  -H "Content-Type: application/json" \
  -d '{"content": "hello", "from_lang": "en", "to_lang": "el"}'

# Response: {"translation": "γεια σας", "cached": false, "cache_key": "..."}

# Step 3: User sees translation, maybe modifies it or keeps it
# Step 4: User clicks "Save Flashcard"
curl -X POST http://localhost:8080/flashcards \
  -H "Content-Type: application/json" \
  -d '{
    "question": "hello",
    "answer": "γεια σας",
    "question_lang": "en",
    "answer_lang": "el",
    "ai_translated_answer": true
  }'

# Flashcard saved with AI flag!
```

## Verification Queries

### Check Database Directly

If you want to verify the data in the database:

```sql
-- See all flashcards with AI tracking
SELECT id, question, answer, question_lang, answer_lang, 
       ai_translated_question, ai_translated_answer,
       created_at
FROM flashcards
ORDER BY created_at DESC;

-- Count flashcards by AI usage
SELECT 
  COUNT(*) as total_cards,
  SUM(CASE WHEN ai_translated_question THEN 1 ELSE 0 END) as ai_questions,
  SUM(CASE WHEN ai_translated_answer THEN 1 ELSE 0 END) as ai_answers
FROM flashcards;
```

## Performance Testing

### Cache Hit Rate

Generate the same translation multiple times and measure response time:

```bash
# First call (cache miss - slower, calls OpenAI)
time curl -X POST http://localhost:8080/flashcards/translate \
  -H "Content-Type: application/json" \
  -d '{"content": "hello", "from_lang": "en", "to_lang": "el"}'

# Second call (cache hit - much faster)
time curl -X POST http://localhost:8080/flashcards/translate \
  -H "Content-Type: application/json" \
  -d '{"content": "hello", "from_lang": "en", "to_lang": "el"}'
```

Expected: Second call should be significantly faster.

## Notes

- Cache TTL: 1 hour by default (configurable via `TRANSLATION_CACHE_TTL`)
- Cache cleanup runs every 10 minutes (configurable via `TRANSLATION_CACHE_CLEANUP_INTERVAL`)
- If OpenAI API key is not set, translation endpoints will return error
- Language codes should be ISO 639-1 (2 lowercase letters)
- Cache is in-memory, so it's lost on server restart

## Health Check

```bash
curl http://localhost:8080/health
```

**Expected:** `{"status": "healthy"}`


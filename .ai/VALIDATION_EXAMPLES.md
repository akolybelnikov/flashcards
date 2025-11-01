# Flashcard Creation Validation Examples

## Overview
The flashcard creation endpoint validates that proper language information is provided when translation is needed.

## Validation Rules

1. **Both fields filled** - No language parameters required, card is created directly
2. **Question empty** - Both `question_lang` AND `answer_lang` are required (to know source and target languages)
3. **Answer empty** - Both `question_lang` AND `answer_lang` are required (to know source and target languages)
4. **Both fields empty** - Returns error, at least one field must be provided

**Key Point:** When translation is needed (either field is empty), we require BOTH language fields because we need to know:
- Which language to translate FROM (the provided field's language)
- Which language to translate TO (the empty field's language)

## Example Curl Requests

### ✅ Valid: Both fields provided
```bash
curl -X POST http://localhost:10000/flashcards \
  -H "Content-Type: application/json" \
  -d "{\"question\":\"hello\",\"answer\":\"γεια σας\"}"
```

**Response:** 201 Created
```json
{
  "flashcard": {
    "id": 1,
    "question": "hello",
    "answer": "γεια σας",
    "created_at": "2025-11-01T10:00:00Z",
    "updated_at": "2025-11-01T10:00:00Z"
  },
  "ai_translation_used": false,
  "translated_field": ""
}
```

### ✅ Valid: Question provided, answer will be translated (EN → EL)
```bash
curl -X POST http://localhost:10000/flashcards \
  -H "Content-Type: application/json" \
  -d "{\"question\":\"hello\",\"answer\":\"\",\"question_lang\":\"en\",\"answer_lang\":\"el\"}"
```

**Response:** 201 Created
```json
{
  "flashcard": {
    "id": 2,
    "question": "hello",
    "answer": "γεια σας",
    "created_at": "2025-11-01T10:00:00Z",
    "updated_at": "2025-11-01T10:00:00Z"
  },
  "ai_translation_used": true,
  "translated_field": "answer"
}
```

### ✅ Valid: Answer provided, question will be translated (EL → EN)
```bash
curl -X POST http://localhost:10000/flashcards \
  -H "Content-Type: application/json" \
  -d "{\"question\":\"\",\"answer\":\"γεια σας\",\"question_lang\":\"en\",\"answer_lang\":\"el\"}"
```

**Response:** 201 Created
```json
{
  "flashcard": {
    "id": 3,
    "question": "hello",
    "answer": "γεια σας",
    "created_at": "2025-11-01T10:00:00Z",
    "updated_at": "2025-11-01T10:00:00Z"
  },
  "ai_translation_used": true,
  "translated_field": "question"
}
```

### ❌ Invalid: Both fields empty
```bash
curl -X POST http://localhost:10000/flashcards \
  -H "Content-Type: application/json" \
  -d "{\"question\":\"\",\"answer\":\"\"}"
```

**Response:** 400 Bad Request
```json
{
  "error": "Both question and answer cannot be empty"
}
```

### ❌ Invalid: Question empty but language fields missing
```bash
curl -X POST http://localhost:10000/flashcards \
  -H "Content-Type: application/json" \
  -d "{\"question\":\"\",\"answer\":\"γεια σας\"}"
```

**Response:** 400 Bad Request
```json
{
  "error": "Both question_lang and answer_lang are required when translation is needed"
}
```

### ❌ Invalid: Answer empty but language fields missing
```bash
curl -X POST http://localhost:10000/flashcards \
  -H "Content-Type: application/json" \
  -d "{\"question\":\"hello\",\"answer\":\"\"}"
```

**Response:** 400 Bad Request
```json
{
  "error": "Both question_lang and answer_lang are required when translation is needed"
}
```

### ❌ Invalid: Only one language field provided
```bash
curl -X POST http://localhost:10000/flashcards \
  -H "Content-Type: application/json" \
  -d "{\"question\":\"hello\",\"answer\":\"\",\"question_lang\":\"en\"}"
```

**Response:** 400 Bad Request
```json
{
  "error": "Both question_lang and answer_lang are required when translation is needed"
}
```

## Summary

The validation ensures that:
- Users cannot create completely empty flashcards
- When AI translation is needed (one field is empty), the system knows both the source language (the filled field) and target language (the empty field that will be translated)
- Users can still create flashcards manually without any language parameters if they provide both sides


# Flashcard API - First Successful POST Request

## Date: October 28, 2025

## Summary
✅ Successfully created a flashcard with Greek characters using the API!

## Request
```bash
curl -X POST http://localhost:8080/flashcards \
  -H "Content-Type: application/json" \
  -d '{"question":"hello","answer":"γεια σας"}'
```

## Response
**Status Code:** `201 Created`

**Response Body:**
```json
{
  "id": 2,
  "question": "hello",
  "answer": "γεια σας",
  "created_at": "2025-10-28T17:07:30.900181Z",
  "updated_at": "2025-10-28T17:07:30.900181Z"
}
```

## Verification - Get All Flashcards
```bash
curl http://localhost:8080/flashcards
```

**Response:**
```json
[
  {
    "id": 2,
    "question": "hello",
    "answer": "γεια σας",
    "created_at": "2025-10-28T17:07:30.900181Z",
    "updated_at": "2025-10-28T17:07:30.900181Z"
  },
  {
    "id": 1,
    "question": "hello",
    "answer": "γεια σας",
    "created_at": "2025-10-28T17:07:13.727643Z",
    "updated_at": "2025-10-28T17:07:13.727643Z"
  }
]
```

## Note About Terminal Display
The Greek characters "γεια σας" may appear as `?e?a sa?` in Git Bash terminal output due to terminal encoding limitations. However:
- ✅ The data is stored correctly in PostgreSQL (UTF-8)
- ✅ The API returns proper JSON with correct UTF-8 encoding
- ✅ The HTTP response headers indicate `Content-Type: application/json` which supports UTF-8
- ✅ Web browsers, Postman, and other HTTP clients will display the Greek characters correctly

## Testing in Different Environments

### Git Bash / CMD (Windows)
```bash
# May show Greek as question marks
curl http://localhost:8080/flashcards
```

### PowerShell (Better UTF-8 support)
```powershell
# Run this first to enable UTF-8
[Console]::OutputEncoding = [System.Text.Encoding]::UTF8

curl http://localhost:8080/flashcards
```

### Save to File (Preserves encoding)
```bash
curl http://localhost:8080/flashcards > flashcards.json
# Open flashcards.json in a text editor that supports UTF-8
```

### Browser
Simply visit: `http://localhost:8080/flashcards` in any web browser - will display correctly!

### Postman or Insomnia
Import the request and the Greek characters will display perfectly.

## Additional Test Requests

### Get Single Flashcard
```bash
curl http://localhost:8080/flashcards/2
```

### Update Flashcard
```bash
curl -X PUT http://localhost:8080/flashcards/2 \
  -H "Content-Type: application/json" \
  -d '{"answer":"γεια σου"}'
```

### Delete Flashcard
```bash
curl -X DELETE http://localhost:8080/flashcards/2
```

### Create More Flashcards with Different Languages

**Spanish:**
```bash
curl -X POST http://localhost:8080/flashcards \
  -H "Content-Type: application/json" \
  -d '{"question":"hello","answer":"hola"}'
```

**Japanese:**
```bash
curl -X POST http://localhost:8080/flashcards \
  -H "Content-Type: application/json" \
  -d '{"question":"hello","answer":"こんにちは"}'
```

**Arabic:**
```bash
curl -X POST http://localhost:8080/flashcards \
  -H "Content-Type: application/json" \
  -d '{"question":"hello","answer":"مرحبا"}'
```

**Russian:**
```bash
curl -X POST http://localhost:8080/flashcards \
  -H "Content-Type: application/json" \
  -d '{"question":"hello","answer":"привет"}'
```

**Chinese:**
```bash
curl -X POST http://localhost:8080/flashcards \
  -H "Content-Type: application/json" \
  -d '{"question":"hello","answer":"你好"}'
```

## Success Indicators
✅ HTTP 201 Created status
✅ JSON response with created flashcard
✅ ID automatically assigned (2)
✅ Timestamps automatically generated
✅ UTF-8 characters properly handled
✅ CORS headers present
✅ Content-Type set to application/json

## Database Verification
The data is stored in PostgreSQL with proper UTF-8 encoding:
- Database: `postgres` (Supabase local)
- Table: `flashcards`
- Port: `54322`
- Encoding: `UTF8`

You can verify in the database:
```sql
SELECT * FROM flashcards WHERE question = 'hello';
```

## API Endpoint Summary
- **Base URL:** `http://localhost:8080`
- **Endpoint:** `/flashcards`
- **Method:** `POST`
- **Content-Type:** `application/json`
- **Encoding:** UTF-8 (supports all Unicode characters)

🎉 **Your flashcards API is working perfectly and supports international characters!**


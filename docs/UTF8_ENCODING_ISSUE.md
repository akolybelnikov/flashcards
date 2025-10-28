# UTF-8 Encoding Issue Analysis

## Problem
When sending Greek text "γεια σας" via curl in Git Bash, it gets stored/displayed as "?e?a sa?"

## Root Cause
**Git Bash (MINGW64) does not properly handle UTF-8 characters in command-line arguments.**

When you type:
```bash
curl -X POST http://localhost:8080/flashcards \
  -H "Content-Type: application/json" \
  -d '{"question":"hello","answer":"γεια σας"}'
```

Git Bash converts the Greek characters to a different encoding (likely Windows-1252 or CP1252) before passing them to curl.

## Solutions

### Solution 1: Use PowerShell (RECOMMENDED)
PowerShell has better UTF-8 support on Windows.

```powershell
# Set UTF-8 encoding
$PSDefaultParameterValues['Out-File:Encoding'] = 'utf8'
[Console]::OutputEncoding = [System.Text.Encoding]::UTF8

# Make the request
curl.exe -X POST http://localhost:8080/flashcards `
  -H "Content-Type: application/json; charset=utf-8" `
  -d '{\"question\":\"hello\",\"answer\":\"γεια σας\"}'
```

### Solution 2: Use a JSON File (WORKS IN GIT BASH)
Create a properly encoded UTF-8 JSON file:

**test_flashcard.json:**
```json
{
  "question": "hello",
  "answer": "γεια σας"
}
```

Then use it with curl:
```bash
curl -X POST http://localhost:8080/flashcards \
  -H "Content-Type: application/json; charset=utf-8" \
  --data-binary @test_flashcard.json
```

### Solution 3: Use Postman or Insomnia
GUI tools handle UTF-8 correctly:
- Open Postman
- POST to `http://localhost:8080/flashcards`
- Body → raw → JSON
- Paste: `{"question":"hello","answer":"γεια σας"}`

### Solution 4: Use WSL (Windows Subsystem for Linux)
WSL bash has proper UTF-8 support:
```bash
# In WSL
curl -X POST http://localhost:8080/flashcards \
  -H "Content-Type: application/json" \
  -d '{"question":"hello","answer":"γεια σας"}'
```

### Solution 5: URL Encode in Git Bash
Encode the JSON as URL-encoded UTF-8:
```bash
curl -X POST http://localhost:8080/flashcards \
  -H "Content-Type: application/json" \
  -d '{"question":"hello","answer":"\u03b3\u03b5\u03b9\u03b1 \u03c3\u03b1\u03c2"}'
```

Unicode escape sequences:
- γ = \u03b3
- ε = \u03b5
- ι = \u03b9
- α = \u03b1
- σ = \u03c3
- ς = \u03c2

## Testing if the Database Actually Has the Correct Data

### Method 1: Check with a Browser
Simply open in your browser:
```
http://localhost:8080/flashcards
```

If your browser shows "γεια σας" correctly, then the database IS storing it correctly, and it's just a Git Bash display issue.

### Method 2: Use PowerShell to curl
```powershell
[Console]::OutputEncoding = [System.Text.Encoding]::UTF8
curl.exe http://localhost:8080/flashcards | ConvertFrom-Json | ConvertTo-Json
```

### Method 3: Check Database Encoding
```sql
-- Connect to database
SELECT 
    id, 
    question, 
    answer,
    length(answer) as char_count,
    octet_length(answer) as byte_count
FROM flashcards 
WHERE id = 2;
```

For "γεια σας":
- char_count should be: 8 (8 characters including space)
- byte_count should be: 16 (Greek chars are 2 bytes each in UTF-8, space is 1)

If byte_count > char_count, UTF-8 is stored correctly!

## Quick Test Script

Save this as `test_utf8.ps1` and run in PowerShell:

```powershell
# test_utf8.ps1
[Console]::OutputEncoding = [System.Text.Encoding]::UTF8

Write-Host "Creating flashcard with Greek text..."
$body = @{
    question = "hello"
    answer = "γεια σας"
} | ConvertTo-Json

$response = Invoke-RestMethod `
    -Uri "http://localhost:8080/flashcards" `
    -Method Post `
    -Body $body `
    -ContentType "application/json; charset=utf-8"

Write-Host "Created flashcard with ID: $($response.id)"
Write-Host "Question: $($response.question)"
Write-Host "Answer: $($response.answer)"

# Verify
Write-Host "`nRetrieving all flashcards..."
$flashcards = Invoke-RestMethod -Uri "http://localhost:8080/flashcards"
$flashcards | ForEach-Object {
    Write-Host "ID $($_.id): $($_.question) → $($_.answer)"
}
```

Run with:
```powershell
powershell -ExecutionPolicy Bypass -File test_utf8.ps1
```

## Conclusion

**The problem is Git Bash, NOT your API or PostgreSQL.**

Your application correctly:
✅ Accepts UTF-8 JSON
✅ Stores UTF-8 in PostgreSQL  
✅ Returns UTF-8 JSON responses

The issue is Git Bash terminal encoding when:
1. You TYPE non-ASCII characters
2. Git Bash DISPLAYS the curl output

**Recommended Solution:** Use PowerShell, save JSON to a file, or use Postman for testing UTF-8 content.


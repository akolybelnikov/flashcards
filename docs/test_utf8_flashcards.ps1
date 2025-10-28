# test_utf8_flashcards.ps1
# PowerShell script to test UTF-8 encoding with flashcards API

# Ensure UTF-8 encoding
$PSDefaultParameterValues['Out-File:Encoding'] = 'utf8'
[Console]::OutputEncoding = [System.Text.Encoding]::UTF8
$OutputEncoding = [System.Text.Encoding]::UTF8

Write-Host "=== Flashcards UTF-8 Test ===" -ForegroundColor Cyan
Write-Host ""

# Test 1: Create flashcard with Greek
Write-Host "Test 1: Creating flashcard with Greek text..." -ForegroundColor Yellow
$greekBody = @{
    question = "hello"
    answer = "γεια σας"
} | ConvertTo-Json -Compress

try {
    $response1 = Invoke-RestMethod `
        -Uri "http://localhost:8080/flashcards" `
        -Method Post `
        -Body ([System.Text.Encoding]::UTF8.GetBytes($greekBody)) `
        -ContentType "application/json; charset=utf-8"

    Write-Host "✓ Created flashcard ID: $($response1.id)" -ForegroundColor Green
    Write-Host "  Question: $($response1.question)" -ForegroundColor White
    Write-Host "  Answer: $($response1.answer)" -ForegroundColor White
    Write-Host "  Answer bytes: $([System.Text.Encoding]::UTF8.GetByteCount($response1.answer))" -ForegroundColor Gray
} catch {
    Write-Host "✗ Failed: $_" -ForegroundColor Red
}

Write-Host ""

# Test 2: Create flashcard with Japanese
Write-Host "Test 2: Creating flashcard with Japanese text..." -ForegroundColor Yellow
$japaneseBody = @{
    question = "hello"
    answer = "こんにちは"
} | ConvertTo-Json -Compress

try {
    $response2 = Invoke-RestMethod `
        -Uri "http://localhost:8080/flashcards" `
        -Method Post `
        -Body ([System.Text.Encoding]::UTF8.GetBytes($japaneseBody)) `
        -ContentType "application/json; charset=utf-8"

    Write-Host "✓ Created flashcard ID: $($response2.id)" -ForegroundColor Green
    Write-Host "  Question: $($response2.question)" -ForegroundColor White
    Write-Host "  Answer: $($response2.answer)" -ForegroundColor White
} catch {
    Write-Host "✗ Failed: $_" -ForegroundColor Red
}

Write-Host ""

# Test 3: Retrieve all flashcards
Write-Host "Test 3: Retrieving all flashcards..." -ForegroundColor Yellow
try {
    $flashcards = Invoke-RestMethod -Uri "http://localhost:8080/flashcards"

    Write-Host "✓ Retrieved $($flashcards.Count) flashcards:" -ForegroundColor Green
    $flashcards | ForEach-Object {
        Write-Host "  [$($_.id)] $($_.question) → $($_.answer)" -ForegroundColor White
    }
} catch {
    Write-Host "✗ Failed: $_" -ForegroundColor Red
}

Write-Host ""
Write-Host "=== Test Complete ===" -ForegroundColor Cyan
Write-Host ""
Write-Host "If you see Greek/Japanese characters correctly above, UTF-8 is working!" -ForegroundColor Green
Write-Host "If you see question marks, the issue is in Git Bash terminal encoding." -ForegroundColor Yellow


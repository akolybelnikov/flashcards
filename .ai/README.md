# Documentation & Test Files

This folder contains documentation, guides, and test files for the Flashcards API.

## Documentation Files

- **FLASHCARDS_IMPLEMENTATION.md** - Complete implementation guide for the flashcards feature
- **TODO_REMOVAL_SUMMARY.md** - Summary of removing todo functionality from the codebase
- **UTF8_ENCODING_ISSUE.md** - Detailed explanation of UTF-8 encoding issues with Git Bash
- **GREEK_ENCODING_SOLVED.md** - Quick reference guide for solving Greek/UTF-8 character issues
- **FIRST_API_TEST.md** - First successful API test documentation

## Test Files

- **test_utf8_flashcards.ps1** - PowerShell script for testing UTF-8 encoding
- **greek_unicode.json** - Test JSON file with Unicode escape sequences
- **test_greek.json** - Test file with Greek characters
- **test_greek_proper.json** - Alternative Greek test file

## Quick Links

### For UTF-8 Testing Issues
See: `GREEK_ENCODING_SOLVED.md` for the quick solution

### For API Implementation Details
See: `FLASHCARDS_IMPLEMENTATION.md` for architecture and endpoints

### For Testing with PowerShell
Run: `powershell -ExecutionPolicy Bypass -File test_utf8_flashcards.ps1`

### For Testing with curl in Git Bash
Use: `curl -X POST http://localhost:8080/flashcards -H "Content-Type: application/json" -d @greek_unicode.json`


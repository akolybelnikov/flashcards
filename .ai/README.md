# .ai Directory - AI-Generated Documentation & Artifacts

This directory contains all AI-generated documentation, implementation plans, testing guides, and other artifacts created during the development of the flashcard application.

## Purpose

The `.ai` folder serves as a centralized location for:
- Implementation plans and design decisions
- Testing guides and examples
- Feature documentation
- Troubleshooting guides
- Development notes and context

## 📚 Current Documentation

### 🎯 Start Here
- **`COMPLETE_SUMMARY.md`** ⭐⭐⭐ - **READ THIS FIRST** - Complete implementation summary (Nov 2, 2025)

### Core Implementation Docs
- **`IMPLEMENTATION_SUMMARY.md`** ⭐ - Complete overview of stateful translation feature (Nov 2, 2025)
- **`STATEFUL_TRANSLATION_PLAN.md`** ⭐ - Detailed implementation plan with architecture
- **`FLASHCARDS_IMPLEMENTATION.md`** - Original flashcard implementation notes

### Testing & Quality
- **`MOCK_GENERATION_PLAN.md`** ⭐ - Comprehensive plan for replacing handwritten mocks with gomock (Nov 2, 2025)
- **`MOCK_VISUAL_GUIDE.md`** ⭐ - Visual guide and comparison of mocking approaches (Nov 2, 2025)
- **`TESTING_GUIDE.md`** ⭐ - Comprehensive curl commands for testing all endpoints
- **`TESTS_FIXED.md`** ⭐ - Documentation of test fixes for v2.0 (Nov 2, 2025)
- **`VALIDATION_EXAMPLES.md`** - Validation examples and error cases
- **`FIRST_API_TEST.md`** - Initial API testing documentation

### API Documentation
- **`OPENAPI_UPDATE.md`** ⭐ - OpenAPI v2.0.0 changes and migration guide (Nov 2, 2025)
- **`OPENAPI_COMPLETE.md`** ⭐ - OpenAPI update completion summary (Nov 2, 2025)

### Deployment & Operations
- **`DEPLOY_QUICK_START.md`** ⭐ - Quick 3-step deployment guide (Nov 2, 2025)
- **`DEPLOY_MIGRATIONS.md`** ⭐ - Comprehensive migration deployment guide (Nov 2, 2025)
- **`REMOTE_RESET_GUIDE.md`** ⭐ - Reset remote database guide (Nov 2, 2025)
- **`RESET_READY.md`** - Quick reference for database reset
- **`MAKEFILE_FIX.md`** ⭐ - Windows compatibility fixes (Nov 2, 2025)

### Technical Decisions
- **`LANGUAGE_DECISION.md`** ⭐ - Language handling strategy and ISO 639-1 usage
- **`SUPPORTED_LANGUAGES.md`** ⭐ - Language codes reference for frontend development

### Project Organization
- **`DIRECTORY_MIGRATION.md`** ⭐ - Documentation about ai/ → .ai/ migration (Nov 2, 2025)

### Troubleshooting
- **`GREEK_ENCODING_SOLVED.md`** - Greek character encoding issue resolution
- **`UTF8_ENCODING_ISSUE.md`** - UTF-8 encoding troubleshooting

### Test Data & Scripts
- `greek_unicode.json` - Greek character test data
- `test_greek.json` - Test flashcards with Greek
- `test_greek_proper.json` - Properly encoded Greek test data
- `test_utf8_flashcards.ps1` - PowerShell script for UTF-8 testing

⭐ = Recently added/updated (Nov 2, 2025)  
⭐⭐⭐ = **Essential reading**

## 🚀 Quick Start

### To understand the latest feature:
1. Read `IMPLEMENTATION_SUMMARY.md` for overview
2. Check `TESTING_GUIDE.md` for testing instructions
3. Review `STATEFUL_TRANSLATION_PLAN.md` for detailed design

### To test the API:
```bash
# See TESTING_GUIDE.md for full examples

# Test translation endpoint
curl -X POST http://localhost:8080/flashcards/translate \
  -H "Content-Type: application/json" \
  -d '{"content": "hello", "from_lang": "en", "to_lang": "el"}'

# Create flashcard
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

### For UTF-8 Issues:
See `GREEK_ENCODING_SOLVED.md` for quick solution

## 📁 Maintenance

All future AI-generated artifacts should be placed in this directory to:
1. Keep the repository root clean
2. Maintain a single source of truth for documentation
3. Make it easy to find implementation context
4. Track the evolution of features

## 📝 Convention

When creating new documentation:
- Use descriptive UPPERCASE filenames for markdown docs (e.g., `FEATURE_NAME.md`)
- Include creation date in the document header
- Reference related documents for cross-referencing
- Keep test data files lowercase with underscores (e.g., `test_data.json`)
- Add ⭐ to README for newly added docs

## 🔍 Git Tracking

This directory IS tracked in git to:
- Preserve implementation context
- Share knowledge across team members
- Document design decisions
- Maintain testing procedures

**Note:** While this is a hidden directory (starts with `.`), it's intentionally tracked for documentation purposes.

## 📅 Recent Updates

**November 2, 2025:**
- ✅ Implemented stateful translation workflow
- ✅ Added translation caching with TTL
- ✅ New endpoint: `POST /flashcards/translate`
- ✅ Database migration: AI tracking columns
- ✅ Updated language handling (ISO 639-1)
- 📚 Created comprehensive testing guide

## 🔗 Related Files

- Main README: `../README.md`
- OpenAPI Spec: `../openapi.yaml`
- Migrations: `../supabase/migrations/`
- Source Code: `../cmd/`, `../handlers/`, `../services/`, `../models/`


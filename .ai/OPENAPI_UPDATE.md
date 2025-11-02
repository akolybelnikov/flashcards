# OpenAPI Specification Update - v2.0.0 ✅

**Date:** November 2, 2025  
**File:** `openapi.yaml`  
**Status:** Updated to reflect stateful translation workflow

## Changes Summary

### Version Update
- **From:** v1.0.0
- **To:** v2.0.0 (breaking changes)

### New Endpoint Added

#### `POST /flashcards/translate`
Generate AI-powered translations with caching.

**Request:**
```json
{
  "content": "hello",
  "from_lang": "en",
  "to_lang": "el"
}
```

**Response:**
```json
{
  "translation": "γεια σας",
  "cached": false,
  "cache_key": "a1b2c3d4..."
}
```

**Features:**
- ✅ Translation caching (1 hour TTL)
- ✅ Cache hit/miss indicator
- ✅ ISO 639-1 language codes
- ✅ Validation for language format

### Breaking Changes

#### `POST /flashcards`
**Before (v1.0):**
- Could provide only one field (question OR answer)
- Automatic translation if one field empty
- Required `question_lang` and `answer_lang` for translation

**After (v2.0):**
- ⚠️ **Both `question` and `answer` required**
- No automatic translation
- Optional `ai_translated_question` and `ai_translated_answer` flags
- Optional `question_lang` and `answer_lang` for tracking

**Migration:**
```javascript
// OLD (v1.0) - No longer works
POST /flashcards
{
  "question": "hello",
  "question_lang": "en",
  "answer_lang": "el"
}

// NEW (v2.0) - Required workflow
// Step 1: Generate translation
POST /flashcards/translate
{
  "content": "hello",
  "from_lang": "en",
  "to_lang": "el"
}
// Returns: {"translation": "γεια σας", ...}

// Step 2: Create flashcard
POST /flashcards
{
  "question": "hello",
  "answer": "γεια σας",
  "question_lang": "en",
  "answer_lang": "el",
  "ai_translated_answer": true
}
```

### Schema Updates

#### `Flashcard` Schema (Updated)
**Added fields:**
- `question_lang` (string, nullable) - Language code for question
- `answer_lang` (string, nullable) - Language code for answer
- `ai_translated_question` (boolean) - AI generation flag
- `ai_translated_answer` (boolean) - AI generation flag

**Example:**
```json
{
  "id": 1,
  "question": "hello",
  "answer": "γεια σας",
  "question_lang": "en",
  "answer_lang": "el",
  "ai_translated_question": false,
  "ai_translated_answer": true,
  "created_at": "2025-11-02T10:00:00Z",
  "updated_at": "2025-11-02T10:00:00Z"
}
```

#### `CreateFlashcardRequest` Schema (Updated)
**Changed:**
- `question` - Now **required** (was optional)
- `answer` - Now **required** (was optional)
- `question_lang` - Now **optional** (was conditionally required)
- `answer_lang` - Now **optional** (was conditionally required)

**Added:**
- `ai_translated_question` (boolean, optional, default: false)
- `ai_translated_answer` (boolean, optional, default: false)

#### New Schemas Added

**`GenerateTranslationRequest`:**
```yaml
required:
  - content
  - from_lang
  - to_lang
properties:
  content: string (min 1 char)
  from_lang: string (pattern: ^[a-z]{2}$)
  to_lang: string (pattern: ^[a-z]{2}$)
```

**`GenerateTranslationResponse`:**
```yaml
required:
  - translation
  - cached
  - cache_key
properties:
  translation: string
  cached: boolean
  cache_key: string (SHA-256 hash)
```

#### Removed Schemas

**`CreateFlashcardResponse`** - Removed (now returns `Flashcard` directly)

### Tag Updates

**Added:**
- `Translation` - AI-powered translation generation with caching

**Updated:**
- `Flashcards` - CRUD operations for managing flashcards (updated description)
- `Health` - Health check endpoints (no change)

### Documentation Improvements

#### Updated Descriptions
- ✅ API info section - New workflow explanation
- ✅ Key features list - Reflects v2.0 capabilities
- ✅ Endpoint descriptions - Detailed with examples
- ✅ Language codes - ISO 639-1 reference links

#### Added Examples
- ✅ Translation cache hit/miss scenarios
- ✅ Manual card creation (no AI)
- ✅ AI-assisted card creation
- ✅ Multiple language pair examples
- ✅ Error responses for all endpoints

#### Improved Error Documentation
- ✅ Validation errors with examples
- ✅ Missing required fields
- ✅ Invalid language codes
- ✅ AI unavailable scenarios

### Language Support

**Before:** Limited to `en` and `el` (enum)
**After:** Any ISO 639-1 code (pattern validation)

**Validation:**
- Pattern: `^[a-z]{2}$`
- Must be 2 lowercase letters
- Examples: en, el, fr, de, es, it, etc.

### Response Format Changes

#### `POST /flashcards`
**Before:**
```json
{
  "flashcard": { ... },
  "ai_translation_used": true,
  "translated_field": "answer"
}
```

**After:**
```json
{
  "id": 1,
  "question": "hello",
  "answer": "γεια σας",
  "ai_translated_question": false,
  "ai_translated_answer": true,
  ...
}
```

Direct `Flashcard` object (simpler, more consistent).

### Health Check Update

**Before:**
```json
{
  "status": "ok",
  "timestamp": "2025-11-01T10:00:00Z"
}
```

**After:**
```json
{
  "status": "healthy"
}
```

Matches actual implementation.

## Validation & Testing

### OpenAPI Validation
```bash
# Validate the spec
npx @apidevtools/swagger-cli validate openapi.yaml

# Or use online validator
# https://editor.swagger.io/
```

### Generate Client Code
```bash
# Generate TypeScript client
npx @openapitools/openapi-generator-cli generate \
  -i openapi.yaml \
  -g typescript-fetch \
  -o ./client

# Generate Go client
openapi-generator-cli generate \
  -i openapi.yaml \
  -g go \
  -o ./go-client
```

### Interactive Documentation
```bash
# Serve with Swagger UI
docker run -p 8081:8080 \
  -e SWAGGER_JSON=/openapi.yaml \
  -v $(pwd):/usr/share/nginx/html \
  swaggerapi/swagger-ui
```

## API Documentation URLs

### Production
- **Swagger UI:** https://flashcards-hqwc.onrender.com/docs (if configured)
- **OpenAPI JSON:** https://flashcards-hqwc.onrender.com/openapi.yaml (if served)

### Local Development
- **Swagger UI:** http://localhost:8080/docs (if configured)
- **OpenAPI File:** `./openapi.yaml`

## Client Integration Notes

### For Frontend Developers

1. **Update API client** to v2.0.0
2. **Implement new workflow:**
   - Call `/flashcards/translate` first
   - Show translation to user
   - Call `/flashcards` to save
3. **Update types/interfaces:**
   - Add `ai_translated_question` and `ai_translated_answer` fields
   - Add `question_lang` and `answer_lang` (nullable)
4. **Remove old logic:**
   - Don't send empty question/answer to `/flashcards`
   - Don't expect `CreateFlashcardResponse` wrapper

### Breaking Change Checklist

- [ ] Update API client library
- [ ] Implement translation endpoint calls
- [ ] Update flashcard creation flow
- [ ] Update TypeScript/interface types
- [ ] Test with new response format
- [ ] Handle cache indicators in UI
- [ ] Update error handling

## Files Modified

- ✅ `openapi.yaml` - Complete rewrite for v2.0.0

## Related Documentation

- **Implementation:** `.ai/IMPLEMENTATION_SUMMARY.md`
- **Testing:** `.ai/TESTING_GUIDE.md`
- **API Design:** `.ai/STATEFUL_TRANSLATION_PLAN.md`
- **Languages:** `.ai/SUPPORTED_LANGUAGES.md`

## Next Steps

1. ✅ OpenAPI spec updated
2. ⏳ Deploy to production
3. ⏳ Update API documentation site
4. ⏳ Generate client libraries
5. ⏳ Update frontend to use v2.0 API
6. ⏳ Version old API endpoints (if needed)

---

**Note:** This is a **breaking change**. Consider:
- API versioning strategy (e.g., `/v2/flashcards`)
- Deprecation timeline for v1.0
- Client migration guide
- Backward compatibility layer (if needed)


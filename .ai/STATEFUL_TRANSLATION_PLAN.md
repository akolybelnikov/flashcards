# Stateful Flashcard Creation with AI Translation Caching - Implementation Plan

## Overview
Separate the AI translation generation from flashcard database persistence, introducing a stateful workflow where users can generate, review, accept/reject translations before saving cards.

## Key Requirements
1. **Generate Translation Endpoint**: New endpoint to generate AI translation without saving to DB
2. **Translation Caching**: Store generated translations temporarily to prevent redundant API calls
3. **One Translation Per Session**: Prevent users from generating translations indefinitely for same content
4. **Separation of Concerns**: Translation generation ≠ Card creation
5. **Future-proof**: Support marking cards as "AI-translated" for update scenarios

---

## Session Management Strategy Options

### Option 1: Redis Cache (RECOMMENDED)
**Pros:**
- Industry standard for session/cache management
- Built-in TTL (time-to-live) for automatic cleanup
- High performance
- Scalable across multiple server instances
- Simple key-value storage

**Cons:**
- Additional infrastructure dependency
- Requires Redis server/container

**Implementation:**
- Key format: `translation:{hash(question+answer+from+to)}`
- TTL: 1 hour (configurable)
- Store: translation result + timestamp + attempt count

### Option 2: In-Memory Map with Sync (Simple, for MVP)
**Pros:**
- No external dependencies
- Simple implementation
- Good for single-instance deployments
- Fast access

**Cons:**
- Lost on server restart
- Not shared across multiple instances
- Requires manual cleanup/TTL management
- Memory grows unbounded without cleanup

**Implementation:**
- Use `sync.Map` for thread-safe access
- Background goroutine for cleanup (every 10 mins)
- Key: hash of (content + languages)

### Option 3: Temporary Database Table
**Pros:**
- Persists across restarts
- Uses existing infrastructure
- Can query/audit translation history

**Cons:**
- Slower than in-memory
- Requires migration
- More complex cleanup

**Recommendation:** Start with **Option 2 (In-Memory)** for MVP, design with interface for easy swap to Redis later.

---

## Database Schema Changes

### Migration: Add AI Translation Tracking
**File:** `supabase/migrations/20251102000000_add_translation_tracking.sql`

```sql
-- Add columns to track AI translation usage
ALTER TABLE flashcards 
ADD COLUMN ai_translated_question BOOLEAN DEFAULT FALSE,
ADD COLUMN ai_translated_answer BOOLEAN DEFAULT FALSE,
ADD COLUMN question_lang VARCHAR(10),
ADD COLUMN answer_lang VARCHAR(10);

-- Add index for querying by translation status
CREATE INDEX idx_flashcards_ai_translated 
ON flashcards(ai_translated_question, ai_translated_answer);

-- Add comments for clarity
COMMENT ON COLUMN flashcards.ai_translated_question IS 'Indicates if the question was AI-generated';
COMMENT ON COLUMN flashcards.ai_translated_answer IS 'Indicates if the answer was AI-generated';
COMMENT ON COLUMN flashcards.question_lang IS 'Language code for question (e.g., en, el)';
COMMENT ON COLUMN flashcards.answer_lang IS 'Language code for answer (e.g., en, el)';
```

---

## Implementation Steps

### Step 1: Create Translation Cache Service
**File:** `services/translation_cache.go`

**Components:**
- `TranslationCache` interface
- `InMemoryTranslationCache` implementation
- Cache entry structure with metadata
- Hash function for generating cache keys
- Background cleanup routine
- TTL management

**Key Methods:**
- `Get(key string) (*CachedTranslation, bool)`
- `Set(key string, translation *CachedTranslation)`
- `GenerateKey(content, fromLang, toLang string) string`
- `Cleanup()` - background task

### Step 2: Update Models
**File:** `models/flashcard.go`

**Add:**
1. `GenerateTranslationRequest` struct
   - `Content string` (the text to translate)
   - `FromLang string` (source language)
   - `ToLang string` (target language)
   
2. `GenerateTranslationResponse` struct
   - `Translation string`
   - `Cached bool` (was this from cache?)
   - `CacheKey string` (for client to send back if accepted)
   
3. Update `CreateFlashcardRequest`
   - Add `AITranslatedQuestion *bool`
   - Add `AITranslatedAnswer *bool`
   - Keep existing `QuestionLang` and `AnswerLang`

4. Update `Flashcard` model
   - Add `AITranslatedQuestion bool`
   - Add `AITranslatedAnswer bool`
   - Add `QuestionLang string`
   - Add `AnswerLang string`

5. Update `UpdateFlashcardRequest`
   - Consider if translation regeneration allowed on update
   - For now: updates are manual only (no AI on update)

### Step 3: Update Service Layer
**File:** `services/flashcard_service.go`

**Modify:**
1. Add `translationCache TranslationCache` field to `FlashcardService`
2. Update constructor to inject cache

**New Method:**
- `GenerateTranslation(req *GenerateTranslationRequest) (*GenerateTranslationResponse, error)`
  - Check cache first using generated key
  - If cached, return immediately with `Cached: true`
  - If not cached, call LLM
  - Store in cache before returning
  - Return with `Cached: false` and cache key

**Update:**
- `CreateFlashcard` - Remove automatic translation logic
  - Expect both `Question` and `Answer` to be filled
  - Accept `AITranslatedQuestion` and `AITranslatedAnswer` flags
  - Validate that at least one field has content
  - Store language codes if provided

### Step 4: Update Repository Layer
**File:** `db/flashcard_repository.go`

**Changes:**
1. Update `Create` method signature/implementation
   - Accept `AITranslatedQuestion`, `AITranslatedAnswer`, `QuestionLang`, `AnswerLang`
   - Update INSERT query to include new columns

2. Update query methods to retrieve new fields

3. Update `Update` method if needed (probably no changes)

### Step 5: Add New Handler
**File:** `handlers/flashcard_handler.go`

**New Method:**
- `GenerateTranslation(w http.ResponseWriter, r *http.Request)`
  - Parse `GenerateTranslationRequest`
  - Validate: content, fromLang, toLang all required
  - Call service.GenerateTranslation
  - Return `GenerateTranslationResponse`

**Update:**
- `CreateFlashcard` - Simplify validation
  - Both question and answer must be provided
  - Language fields optional but recommended
  - AI flags optional (default false)
  - Remove translation logic
  
- `RegisterRoutes` - Add new endpoint
  - `POST /flashcards/translate` -> GenerateTranslation

### Step 6: Update Service Interface
**File:** `services/flashcard_service.go`

**Update `FlashcardServiceInterface`:**
- Remove translation return values from `CreateFlashcard` signature
- Add `GenerateTranslation(req *models.GenerateTranslationRequest) (*models.GenerateTranslationResponse, error)`

### Step 7: Update Tests
**Files:** 
- `handlers/flashcard_handler_test.go`
- `services/flashcard_service_test.go`

**Changes:**
1. Add mock for `TranslationCache`
2. Add tests for `GenerateTranslation` endpoint
   - First call generates translation
   - Second call returns cached result
   - Cache key consistency
3. Update existing `CreateFlashcard` tests
   - Remove translation expectations
   - Add tests with AI flags set
4. Add integration test for full workflow:
   - Generate translation
   - Accept and create card
   - Verify AI flags in database

### Step 8: Update Configuration
**File:** `config/config.go`

**Add:**
- `TranslationCacheTTL time.Duration` (default: 1 hour)
- `TranslationCacheCleanupInterval time.Duration` (default: 10 minutes)

### Step 9: Update Main Application
**File:** `cmd/main.go`

**Changes:**
1. Initialize translation cache
2. Pass cache to service constructor
3. Start cache cleanup goroutine (if in-memory implementation)
4. Graceful shutdown: stop cleanup goroutine

---

## API Endpoint Changes

### New Endpoint: Generate Translation
```
POST /flashcards/translate
Content-Type: application/json

Request:
{
  "content": "hello",
  "from_lang": "en",
  "to_lang": "el"
}

Response: 200 OK
{
  "translation": "γεια σας",
  "cached": false,
  "cache_key": "abc123def456..."
}
```

### Updated Endpoint: Create Flashcard
```
POST /flashcards
Content-Type: application/json

Request:
{
  "question": "hello",
  "answer": "γεια σας",
  "question_lang": "en",
  "answer_lang": "el",
  "ai_translated_question": false,
  "ai_translated_answer": true
}

Response: 201 Created
{
  "flashcard": {
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
}
```

---

## Update Scenarios

### Preventing Redundant Translation on Updates

**Strategy:**
When a user updates a flashcard that already has `ai_translated_question` or `ai_translated_answer` set to `true`, we have options:

1. **Prevent translation regeneration** (RECOMMENDED for MVP)
   - Don't allow calling `/flashcards/translate` with existing card content
   - User must manually edit if they want different translation
   - Keeps it simple and predictable

2. **Allow regeneration with warning**
   - Check if content exists in DB with AI flag
   - Return warning in response: "This content was previously translated"
   - Still allow but log/track

3. **Block completely**
   - `/flashcards/translate` checks if content exists in DB
   - Returns error if already translated

**Recommendation:** Option 1 - Keep translation endpoint stateless, don't check DB. Update endpoint only accepts manual edits.

---

## Workflow Example

### User Creates Card with AI Translation

1. **Frontend:** User types "hello" in English field
2. **Frontend:** User clicks "Generate Greek Translation" button
3. **Frontend → Backend:** `POST /flashcards/translate` with `{content: "hello", from_lang: "en", to_lang: "el"}`
4. **Backend:** Checks cache → miss → calls OpenAI → stores in cache → returns "γεια σας"
5. **Frontend:** Shows translation to user
6. **User:** Clicks "Discard", manually types "γειά σου"
7. **User:** Clicks "Generate" again (wants to see formal version)
8. **Frontend → Backend:** `POST /flashcards/translate` (same request)
9. **Backend:** Checks cache → hit → returns "γεια σας" (cached)
10. **User:** Accepts this time
11. **Frontend → Backend:** `POST /flashcards` with both question and answer, `ai_translated_answer: true`
12. **Backend:** Saves to database with AI flag

### User Edits Existing Card

1. **Frontend:** User opens card #5 for editing
2. **Frontend:** Shows current question/answer with AI flags (grayed out or marked)
3. **User:** Manually edits answer
4. **Frontend → Backend:** `PUT /flashcards/5` with updated answer
5. **Backend:** Updates answer only, doesn't touch AI flags (they remain from creation)

*Note: For MVP, no re-translation on updates. Future enhancement could reset AI flag when field is manually edited.*

---

## Testing Strategy

### Unit Tests
1. **Translation Cache:**
   - Set and get operations
   - Key generation consistency
   - TTL expiration
   - Cleanup routine

2. **Service Layer:**
   - Generate translation (cache miss)
   - Generate translation (cache hit)
   - Create flashcard with AI flags
   - Invalid inputs

3. **Handler Layer:**
   - Valid translation request
   - Invalid translation request
   - Create with AI flags
   - Response format validation

### Integration Tests
1. Full workflow: generate → create
2. Cache behavior across multiple requests
3. Database persistence of AI flags
4. Concurrent requests (race conditions)

---

## Migration Rollout

### Database Migration Steps
1. Run migration to add new columns
2. Existing records get `FALSE` for AI flags (default)
3. Existing records have `NULL` for language codes (acceptable)
4. Deploy new code
5. New cards will have proper flags and languages

### Backward Compatibility
- Existing API behavior changes:
  - `POST /flashcards` no longer auto-translates
  - New endpoint required for translation
  - This is a **breaking change** for any existing clients

**Mitigation:**
- Version API if needed: `/v2/flashcards` vs `/v1/flashcards`
- Or: Document breaking change clearly, this is early development

---

## Future Enhancements

### Phase 2 (Post-MVP)
1. **Redis Integration**
   - Replace in-memory cache with Redis
   - Add Redis configuration
   - Update cache interface implementation

2. **Advanced Session Tracking**
   - User authentication/sessions
   - Per-user translation limits
   - Translation history per user

3. **Smarter Update Handling**
   - Reset AI flags when user manually edits
   - Optional re-translation on update
   - Version history for translations

4. **Analytics**
   - Track cache hit rate
   - Monitor AI translation usage
   - Cost tracking per translation

5. **Multi-language Support**
   - Remove hardcoded en/el
   - Support any language pair
   - Language detection

6. **Batch Translation**
   - Generate multiple cards from list
   - Bulk import with translation

---

## Security Considerations

1. **Rate Limiting:** Prevent abuse of translation endpoint
2. **Cost Control:** Monitor OpenAI API usage
3. **Cache Poisoning:** Validate inputs before caching
4. **Input Sanitization:** Prevent SQL injection with new fields

---

## Success Metrics

1. **Cache Hit Rate:** Target 30%+ (users regenerating translations)
2. **API Cost Reduction:** Measured vs current auto-translate
3. **User Workflow:** Average time from translate → save
4. **Error Rate:** Monitor translation failures

---

## Implementation Priority

**High Priority (MVP):**
- ✅ Step 1: Translation cache (in-memory)
- ✅ Step 2: Update models
- ✅ Step 3: Update service layer
- ✅ Step 4: Update repository
- ✅ Step 5: Add handler
- ✅ Step 6: Update interface
- ✅ Step 9: Wire up in main.go

**Medium Priority (Pre-Launch):**
- ✅ Step 7: Tests
- ✅ Step 8: Configuration
- ✅ Database migration

**Low Priority (Post-Launch):**
- Redis integration
- Rate limiting
- Analytics

---

## Questions to Resolve

1. **Cache TTL:** 1 hour reasonable? Or should it be session-based?
2. **Cache Key:** Include user ID in key when we add auth?
3. **API Versioning:** Break existing API or add v2?
4. **Update Behavior:** Should updates allow re-translation?
5. **Language Codes:** ✅ Use ISO 639-1 codes directly (en, el, fr, etc.) - OpenAI understands them natively. No hardcoded language mappings needed. Frontend will provide dropdown with supported languages.

---

## Estimated Effort

- **Translation Cache Service:** 2-3 hours
- **Models Update:** 1 hour
- **Service Layer:** 2 hours
- **Repository Update:** 1-2 hours
- **Handler Update:** 1-2 hours
- **Tests:** 3-4 hours
- **Migration:** 1 hour
- **Integration & Bug Fixes:** 2-3 hours

**Total:** ~15-20 hours for complete implementation and testing

---

## Conclusion

This plan separates translation generation from card persistence, introducing a stateful workflow that:
- Reduces unnecessary AI API calls through caching
- Gives users control over accepting/rejecting translations
- Tracks which fields were AI-generated for future features
- Maintains clean separation of concerns
- Scales to Redis when needed

The in-memory cache provides a simple MVP solution that can be upgraded to Redis later without changing the service interface.


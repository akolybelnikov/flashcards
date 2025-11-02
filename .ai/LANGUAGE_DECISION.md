# Pre-Implementation Summary - Language Handling Decision

## Decision Made: Remove Hardcoded Language Mappings ✅

### What Changed
- ✅ **Removed:** Hardcoded `langMap` with "en" → "English", "el" → "Greek"
- ✅ **Now:** Use ISO 639-1 codes directly in OpenAI prompts
- ✅ **Validation:** Backend checks format only (2 lowercase letters)
- ✅ **Flexibility:** Any ISO 639-1 language code accepted

### Why This Approach?
1. **OpenAI understands ISO codes** - No need to translate "en" to "English"
2. **Frontend controls languages** - Dropdown will have curated list
3. **Backend validates format** - Security: ensure proper input format
4. **No maintenance** - Add languages to frontend without backend changes
5. **Future-proof** - Easy to expand language support

### Backend Validation Strategy
```go
// Backend checks FORMAT only, not content
func validateLanguageCode(code string) error {
    // Must be 2 lowercase letters
    if len(code) != 2 { error }
    if not [a-z] { error }
    return nil
}
```

**What backend DOES validate:**
- ✅ Length is exactly 2 characters
- ✅ Characters are lowercase letters

**What backend DOES NOT validate:**
- ❌ Not checking against whitelist of "allowed" languages
- ❌ Not verifying if language exists
- ❌ OpenAI will return error if language unsupported (rare)

### Frontend Implementation (Future)
- Create dropdown with ~10-20 common languages
- Map language codes to display names (e.g., "en" → "English")
- Send only the 2-letter code to backend
- Can expand language list anytime

### Testing Approach
For now (pre-frontend):
- Test with "en" and "el" (our main use case)
- Test with other valid codes like "fr", "de", "es"
- Test validation: reject "EN", "eng", "e", "123", ""

### Files Modified
1. **services/llm_client.go** - Removed langMap, added validation
2. **ai/SUPPORTED_LANGUAGES.md** - Created language reference for frontend team

### Ready to Proceed
All questions resolved. Ready to implement the stateful translation workflow!

**Next Step:** Begin implementation following `STATEFUL_TRANSLATION_PLAN.md`


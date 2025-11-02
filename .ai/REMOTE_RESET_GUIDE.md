# Reset Remote Database & Run Migrations

**Date:** November 2, 2025  
**Task:** Reset remote Supabase database and apply all migrations from scratch

---

## ⚠️ WARNING

**This will DELETE ALL DATA in your remote database!**

Use this when:
- ✅ You want a clean slate in production
- ✅ You're okay losing all existing data
- ✅ You need to fix migration conflicts
- ✅ You're in development/testing phase

**DO NOT use this if:**
- ❌ You have production data you need to keep
- ❌ Users are actively using the application
- ❌ You haven't backed up important data

---

## Quick Solution

### Option 1: Using Makefile (Recommended)

```bash
# This will prompt for confirmation
make db-remote-reset
```

**What happens:**
1. Shows warning message
2. Prompts: "Press Ctrl+C to cancel or any key to continue"
3. Drops all tables in remote database
4. Re-applies ALL migrations from scratch
5. Your remote database is now clean with the new schema ✅

### Option 2: Direct Supabase CLI

```bash
supabase db reset --linked
```

---

## Complete Workflow

### Step 1: Ensure You're Linked to Remote

```bash
make db-link
# Or:
supabase link
```

### Step 2: Check Current Migration Status

```bash
make db-status
```

You might see something like:
```
    LOCAL      │ REMOTE │ TIME (UTC)          │ NAME
  ─────────────┼────────┼─────────────────────┼──────────────────────────────────
   20251028... │   ✓    │ 2025-10-28 00:00:00 │ create_flashcards
   20251102... │        │ 2025-11-02 18:44:09 │ add_translation_tracking
```

### Step 3: Reset Remote Database

```bash
make db-remote-reset
```

**You'll see:**
```
========================================
REMOTE DATABASE RESET
========================================
This will:
  1. Drop all tables in remote database
  2. Re-apply all migrations from scratch
  3. All data will be LOST!

WARNING: This affects PRODUCTION database!
========================================

Press Ctrl+C now to cancel, or
Press any key to continue . . .
```

**Type any key to proceed** (or Ctrl+C to cancel)

**Output:**
```
Resetting remote database...
Dropping all objects...
Applying migration 20251028000000_create_flashcards.sql...
Applying migration 20251102184409_add_translation_tracking.sql...
✓ Reset complete
```

### Step 4: Verify

```bash
make db-status
```

**Now all migrations should be applied:**
```
    LOCAL      │ REMOTE │ TIME (UTC)          │ NAME
  ─────────────┼────────┼─────────────────────┼──────────────────────────────────
   20251028... │   ✓    │ 2025-10-28 00:00:00 │ create_flashcards
   20251102... │   ✓    │ 2025-11-02 18:44:09 │ add_translation_tracking
```

✅ **All migrations applied!**

---

## Verify in Supabase Dashboard

1. Go to https://supabase.com/dashboard
2. Select your project
3. Click **Table Editor**
4. Check `flashcards` table has:
   - ✅ `id`
   - ✅ `question`
   - ✅ `answer`
   - ✅ `question_lang` (NEW)
   - ✅ `answer_lang` (NEW)
   - ✅ `ai_translated_question` (NEW)
   - ✅ `ai_translated_answer` (NEW)
   - ✅ `created_at`
   - ✅ `updated_at`

---

## Test Your API

```bash
# Test health endpoint
curl https://flashcards-hqwc.onrender.com/health

# Create a test flashcard
curl -X POST https://flashcards-hqwc.onrender.com/flashcards \
  -H "Content-Type: application/json" \
  -d '{
    "question": "hello",
    "answer": "γεια σας",
    "question_lang": "en",
    "answer_lang": "el",
    "ai_translated_answer": true
  }'

# Get all flashcards
curl https://flashcards-hqwc.onrender.com/flashcards
```

---

## When to Use Each Command

| Command | When to Use |
|---------|-------------|
| `make db-push` | Normal deployment - just applies pending migrations |
| `make db-remote-reset` | Fresh start - drops everything and reapplies all migrations |
| `make db-pull` | Sync local with remote changes made by others |
| `make db-status` | Check what's deployed vs pending |

---

## Alternative: Skip Confirmation

If you're absolutely sure and want to skip the confirmation prompt:

```bash
# BE VERY CAREFUL WITH THIS!
supabase db reset --linked --no-confirm
```

---

## Backup Before Reset (Recommended)

### Export Data Before Reset

```bash
# Export all data from flashcards table
curl https://flashcards-hqwc.onrender.com/flashcards > backup.json
```

### Restore After Reset

```bash
# Restore data after reset (if needed)
# You'd need to write a script to POST each flashcard back
```

---

## Troubleshooting

### "Error: Project not linked"

**Solution:**
```bash
make db-link
```

### "Error: Cannot connect to remote database"

**Solution:**
1. Check internet connection
2. Verify Supabase project is active
3. Check database URL is correct

### Reset Hangs or Times Out

**Solution:**
```bash
# Cancel with Ctrl+C
# Wait a few seconds
# Try again
make db-remote-reset
```

### "Migration history is out of sync"

**Solution:**
```bash
# Reset will fix this - it drops everything and reapplies
make db-remote-reset
```

---

## What Gets Reset

**Dropped:**
- ❌ All tables
- ❌ All data
- ❌ All indexes
- ❌ All functions/triggers
- ❌ Migration history

**Re-created:**
- ✅ All tables (from migrations)
- ✅ All indexes (from migrations)
- ✅ All comments (from migrations)
- ✅ Migration history (clean slate)

**NOT affected:**
- ✅ Authentication users (separate schema)
- ✅ Storage buckets (separate)
- ✅ Edge functions (separate)

---

## Example Session

```bash
# 1. Link to remote
$ make db-link
Enter your project ref: bgarvmskqdrbmzoejikn
✓ Linked

# 2. Check status
$ make db-status
LOCAL: 2 migrations
REMOTE: 1 migration (1 pending)

# 3. Reset and apply all
$ make db-remote-reset
WARNING: This will delete all data!
Press any key to continue...
✓ Reset complete

# 4. Verify
$ make db-status
LOCAL: 2 migrations  ✓
REMOTE: 2 migrations ✓

# 5. Test API
$ curl https://flashcards-hqwc.onrender.com/health
{"status":"healthy"}

✅ Done!
```

---

## Summary

**To reset remote database and apply all migrations:**

```bash
make db-remote-reset
```

This is the nuclear option - use it when you need a clean slate! 🚀

---

## Files

- **Makefile** - Updated with `db-remote-reset` command
- **This Guide** - `.ai/REMOTE_RESET_GUIDE.md`

---

**Status:** ✅ Ready to use  
**Last Updated:** November 2, 2025


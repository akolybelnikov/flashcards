# 🚀 Deploy Migrations to Remote Supabase - Quick Guide

## TL;DR

```bash
# 1. Link to remote Supabase (one-time)
supabase link

# 2. Check migration status
supabase migration list

# 3. Push migrations to production
supabase db push

# Done! ✅
```

---

## Step-by-Step Instructions

### Step 1: Link Your Project (First Time Only)

```bash
cd /c/Users/akoly/GolandProjects/flashcards
supabase link
```

**You'll be asked for:**
1. **Project Reference ID** - Find it in your Supabase dashboard URL
   - Example: `https://supabase.com/dashboard/project/bgarvmskqdrbmzoejikn`
   - The part after `/project/` is your project ref: `bgarvmskqdrbmzoejikn`

2. **Database Password** - Your Supabase project database password
   - Found in: Supabase Dashboard → Settings → Database → Password

**Example interaction:**
```
$ supabase link
Enter your project ref: bgarvmskqdrbmzoejikn
Enter your database password: ••••••••••••
✓ Linked to project: flashcards
```

### Step 2: Check What Needs Deploying

```bash
supabase migration list
```

**You should see:**
```
    LOCAL      │ REMOTE │ TIME (UTC)          │ NAME
  ─────────────┼────────┼─────────────────────┼──────────────────────────────────
   20251028... │   ✓    │ 2025-10-28 00:00:00 │ create_flashcards
   20251102... │        │ 2025-11-02 18:44:09 │ add_translation_tracking
```

The second row is **missing the checkmark** in REMOTE column - that's what we need to deploy!

### Step 3: Push to Production

```bash
supabase db push
```

**Expected output:**
```
Applying migration 20251102184409_add_translation_tracking.sql...
ALTER TABLE
CREATE INDEX
COMMENT
✓ Finished supabase db push.
```

### Step 4: Verify Deployment

```bash
supabase migration list
```

**Now both should have checkmarks:**
```
    LOCAL      │ REMOTE │ TIME (UTC)          │ NAME
  ─────────────┼────────┼─────────────────────┼──────────────────────────────────
   20251028... │   ✓    │ 2025-10-28 00:00:00 │ create_flashcards
   20251102... │   ✓    │ 2025-11-02 18:44:09 │ add_translation_tracking
```

✅ **Done!** Your production database now has the new columns.

---

## Verify in Supabase Dashboard

1. Open https://supabase.com/dashboard
2. Go to your project
3. Click **Table Editor** in sidebar
4. Select `flashcards` table
5. You should see new columns:
   - ✅ `ai_translated_question`
   - ✅ `ai_translated_answer`
   - ✅ `question_lang`
   - ✅ `answer_lang`

---

## Test Your Production API

```bash
# Test that Render can now access new columns
curl https://flashcards-hqwc.onrender.com/flashcards

# Create a flashcard with AI flags
curl -X POST https://flashcards-hqwc.onrender.com/flashcards \
  -H "Content-Type: application/json" \
  -d '{
    "question": "hello",
    "answer": "γεια σας",
    "ai_translated_answer": true
  }'
```

---

## Using the Makefile (Alternative)

If you've updated your Makefile:

```bash
# Link (first time)
make db-link

# Check status
make db-status

# Push migrations
make db-push

# Verify
make db-status
```

---

## Common Issues

### "Error: Project not linked"
**Solution:** Run `supabase link` first

### "Error: permission denied"
**Solution:** Check your database password is correct

### "Migration already applied"
**Solution:** Run `supabase migration list` to see current status

### Dashboard not showing new columns
**Solution:** Hard refresh (Ctrl+Shift+R) or wait 30 seconds

---

## Summary

| Command | Purpose |
|---------|---------|
| `supabase link` | Connect local to remote (once) |
| `supabase migration list` | Check what's deployed |
| `supabase db push` | Deploy migrations |

That's it! Your remote database is now updated. 🎉


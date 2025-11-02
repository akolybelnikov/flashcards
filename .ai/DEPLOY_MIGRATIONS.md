# Deploying Migrations to Remote Supabase

**Date:** November 2, 2025  
**Task:** Apply database migrations to production Supabase

---

## Problem

You've created migrations locally, but they haven't been applied to your remote Supabase database. The remote schema is missing the new columns:
- `ai_translated_question`
- `ai_translated_answer`
- `question_lang`
- `answer_lang`

---

## Solution: Deploy Migrations to Remote

### Step 1: Link to Remote Project (One-time setup)

If you haven't linked your project yet:

```bash
make db-link
# Or directly:
supabase link
```

You'll need:
- **Project Reference ID** - Found in your Supabase dashboard URL: `https://supabase.com/dashboard/project/YOUR_PROJECT_REF`
- **Database Password** - Your Supabase project's database password

**Example:**
```bash
$ supabase link
Enter your project ref: bgarvmskqdrbmzoejikn
Enter your database password: ****
Linked to project: flashcards
```

### Step 2: Check Migration Status

See which migrations are pending:

```bash
make db-status
# Or directly:
supabase migration list
```

**Expected output:**
```
    LOCAL      │ REMOTE │ TIME (UTC)          │ NAME
  ─────────────┼────────┼─────────────────────┼──────────────────────────────────
   20251028... │   ✓    │ 2025-10-28 00:00:00 │ create_flashcards
   20251102... │        │ 2025-11-02 18:44:09 │ add_translation_tracking
```

The missing checkmark (✓) in REMOTE column means the migration hasn't been applied yet.

### Step 3: Push Migrations to Remote

**⚠️ WARNING:** This will apply migrations to your **production database**!

```bash
make db-push
# Or directly:
supabase db push
```

The Makefile will ask for confirmation before proceeding.

**Expected output:**
```
Pushing migrations to remote Supabase...
WARNING: This will apply migrations to production!
Continue? [y/N] y
Applying migration 20251102184409_add_translation_tracking.sql...
✓ Migration applied successfully!
```

### Step 4: Verify Migration Applied

Check the status again:

```bash
make db-status
```

**Expected output (all migrations should have ✓):**
```
    LOCAL      │ REMOTE │ TIME (UTC)          │ NAME
  ─────────────┼────────┼─────────────────────┼──────────────────────────────────
   20251028... │   ✓    │ 2025-10-28 00:00:00 │ create_flashcards
   20251102... │   ✓    │ 2025-11-02 18:44:09 │ add_translation_tracking
```

---

## Verify in Supabase Dashboard

1. Go to https://supabase.com/dashboard
2. Select your project
3. Navigate to **Table Editor** → `flashcards` table
4. Check that new columns exist:
   - `ai_translated_question` (boolean)
   - `ai_translated_answer` (boolean)
   - `question_lang` (varchar)
   - `answer_lang` (varchar)

---

## Troubleshooting

### Error: "Project not linked"

**Problem:** You haven't linked your local project to remote Supabase.

**Solution:**
```bash
make db-link
```

### Error: "Migration already exists remotely"

**Problem:** Migration was partially applied or manually run.

**Solution:**
```bash
# Check migration history
supabase migration list

# If needed, repair migration history
supabase migration repair
```

### Error: "Network unreachable"

**Problem:** Can't connect to remote Supabase.

**Solution:**
1. Check your internet connection
2. Verify Supabase is accessible
3. Check firewall settings

### Error: "Authentication failed"

**Problem:** Database password is incorrect.

**Solution:**
```bash
# Re-link with correct password
supabase link --project-ref YOUR_PROJECT_REF
```

### Migration Applied but Schema Unchanged

**Problem:** Cached schema or browser issue.

**Solution:**
1. Hard refresh Supabase dashboard (Ctrl+Shift+R)
2. Wait 30 seconds for cache to clear
3. Check using SQL editor:
   ```sql
   SELECT column_name, data_type 
   FROM information_schema.columns 
   WHERE table_name = 'flashcards';
   ```

---

## Rolling Back (If Needed)

If something goes wrong, you can rollback:

```bash
# Rollback last migration on remote
supabase db reset --db-url YOUR_REMOTE_DB_URL

# Or manually in Supabase SQL Editor:
ALTER TABLE flashcards DROP COLUMN IF EXISTS ai_translated_question;
ALTER TABLE flashcards DROP COLUMN IF EXISTS ai_translated_answer;
ALTER TABLE flashcards DROP COLUMN IF EXISTS question_lang;
ALTER TABLE flashcards DROP COLUMN IF EXISTS answer_lang;
```

---

## Complete Workflow Summary

```bash
# 1. Link to remote (first time only)
make db-link

# 2. Check what needs to be deployed
make db-status

# 3. Push migrations to production
make db-push

# 4. Verify deployment
make db-status

# 5. Test your API against production database
curl https://flashcards-hqwc.onrender.com/flashcards
```

---

## Environment Variables for Remote

Make sure your Render deployment has the correct `DB_URL`:

**In Render Dashboard:**
1. Go to your service
2. Navigate to **Environment** tab
3. Verify `DB_URL` points to remote Supabase:
   ```
   DB_URL=postgresql://postgres.bgarvmskqdrbmzoejikn:PASSWORD@aws-0-eu-west-1.pooler.supabase.com:6543/postgres
   ```

**Note:** Use the **connection pooler** URL (port 6543), not direct connection (port 5432).

---

## Local vs Remote Commands

| Task | Local (Docker) | Remote (Production) |
|------|----------------|---------------------|
| Start DB | `make db-start` | N/A (always running) |
| Stop DB | `make db-stop` | N/A |
| Apply migrations | `make db-up` | `make db-push` |
| Check status | `make db-list` | `make db-status` |
| Rollback | `make db-down` | Manual SQL or reset |

---

## Best Practices

### Before Pushing to Production

1. ✅ Test migration locally first: `make db-reset`
2. ✅ Verify application works with new schema
3. ✅ Backup production database (Supabase auto-backups daily)
4. ✅ Push during low-traffic hours if possible
5. ✅ Have rollback plan ready

### After Pushing to Production

1. ✅ Verify migration applied: `make db-status`
2. ✅ Check Supabase dashboard for new columns
3. ✅ Test API endpoints
4. ✅ Monitor logs for errors
5. ✅ Verify existing data is intact

---

## Quick Reference

```bash
# Most common workflow:
make db-link          # First time only
make db-status        # Check what's pending
make db-push          # Deploy to production
make db-status        # Verify deployment

# Create new migration:
make db-new           # Creates timestamped file

# Pull remote changes (if team members made changes):
make db-pull          # Syncs remote schema to local
```

---

## File Locations

- **Local migrations:** `supabase/migrations/`
- **Current pending migration:** `20251102184409_add_translation_tracking.sql`
- **Makefile:** Updated with remote commands
- **Supabase config:** `supabase/config.toml`

---

## Next Steps

1. Run `make db-push` to apply migrations to production
2. Verify in Supabase dashboard
3. Test your Render deployment
4. Check that API returns new fields

---

## Support

If you encounter issues:
1. Check Supabase status: https://status.supabase.com/
2. View migration files: `ls supabase/migrations/`
3. Check Supabase logs in dashboard
4. Review migration SQL for syntax errors

---

**TL;DR:** Run `make db-push` to deploy your local migrations to remote Supabase production database.


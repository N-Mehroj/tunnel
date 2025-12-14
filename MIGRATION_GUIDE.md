# 🗂️ Migration System Quick Reference

## ⚡ Quick Commands

```bash
# Create a new migration
make migration NAME=create_posts_table
go run database/migrate.go create create_posts_table

# Run all pending migrations
make migrate-up
go run database/migrate.go up

# Run only 3 pending migrations
make migrate-up-n N=3
go run database/migrate.go up 3

# Rollback last migration
make migrate-down
go run database/migrate.go down

# Rollback last 2 migrations
make migrate-down-n N=2
go run database/migrate.go down 2

# Check migration status
make migrate-status
go run database/migrate.go status
```

## 📁 File Structure

```
database/
├── migrate.go              ← Migration CLI tool
├── migrations/             ← All migration files
│   ├── 20240101000000_create_users_table.up.sql
│   ├── 20240101000000_create_users_table.down.sql
│   ├── 20251214192414_create_auth_table.up.sql
│   └── 20251214192414_create_auth_table.down.sql
└── MIGRATIONS.md           ← Full documentation
```

## ✏️ Create Migration

**Step 1:** Generate migration files
```bash
make migration NAME=add_email_to_posts
```

**Step 2:** Edit the UP migration (`database/migrations/TIMESTAMP_add_email_to_posts.up.sql`)
```sql
BEGIN;

ALTER TABLE posts ADD COLUMN email VARCHAR(255);

COMMIT;
```

**Step 3:** Edit the DOWN migration (`database/migrations/TIMESTAMP_add_email_to_posts.down.sql`)
```sql
BEGIN;

ALTER TABLE posts DROP COLUMN email;

COMMIT;
```

**Step 4:** Run the migration
```bash
make migrate-up
```

## 📊 Migration Tracking

- Migrations are tracked in the `migrations` table
- Each applied migration is recorded with name and timestamp
- DOWN migrations automatically remove the record

## 🔍 Check Status

```bash
make migrate-status
```

Output example:
```
Migration Status:
------------------------------------------------------------
20240101000000_create_users_table              ● Applied
20251214192414_create_auth_table               ○ Pending
```

## 🛑 Common Patterns

### Create a Table
```bash
make migration NAME=create_comments_table
```

### Add Column
```bash
make migration NAME=add_status_to_users
```

### Create Index
```bash
make migration NAME=create_user_email_index
```

### Add Foreign Key
```bash
make migration NAME=add_user_id_to_posts
```

## ⚠️ Important Notes

✅ **DO:**
- Always test DOWN migrations
- Keep migrations small and focused
- Use meaningful names
- Wrap in transactions (BEGIN/COMMIT)
- Test in development first

❌ **DON'T:**
- Mix multiple changes in one migration
- Leave migrations without rollback procedures
- Modify production directly
- Hardcode values that change per environment

## 🐛 Troubleshooting

**Problem:** Migration fails
- Check SQL syntax in the migration file
- Verify database connection in `.env`
- Check current status: `make migrate-status`

**Problem:** Cannot connect to database
- Verify `.env` has correct DB credentials
- Check database is running
- Test connection manually

## 📚 More Information

See `database/MIGRATIONS.md` for complete documentation.

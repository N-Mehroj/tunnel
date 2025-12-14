# 🗄️ Database Migration System

Complete, production-ready database migration system for your Go application.

## 🎯 Features

✅ **Create migrations** - Generate `.up.sql` and `.down.sql` file pairs
✅ **Apply migrations** - Run pending migrations with `migrate-up`
✅ **Rollback migrations** - Undo changes with `migrate-down`
✅ **Track migrations** - Database tracks which migrations have been applied
✅ **Batch operations** - Run/rollback multiple migrations at once
✅ **Status checking** - See which migrations are applied/pending
✅ **Transaction safety** - All migrations wrapped in BEGIN/COMMIT
✅ **Multi-database** - PostgreSQL and MySQL support

## 📋 Commands

### Create a Migration
```bash
# Create migration files
make migration NAME=create_users_table

# Or run directly
go run database/migrate.go create create_users_table
```

### Run Migrations (UP)
```bash
# Run all pending migrations
make migrate-up

# Run only N migrations
make migrate-up-n N=3
```

### Rollback Migrations (DOWN)
```bash
# Rollback last migration
make migrate-down

# Rollback N migrations
make migrate-down-n N=2
```

### Check Status
```bash
# Show migration status
make migrate-status
```

## 📂 Directory Structure

```
database/
├── migrate.go                          ← Migration CLI tool
├── migrations/                         ← All migration files
│   ├── 20240101000000_create_users_table.up.sql
│   ├── 20240101000000_create_users_table.down.sql
│   ├── 20251214192414_create_auth_table.up.sql
│   └── 20251214192414_create_auth_table.down.sql
└── MIGRATIONS.md                       ← Detailed documentation
```

## 🚀 Getting Started

### 1. Setup Environment
Create `.env` file:
```env
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=your_database
```

### 2. Create Your First Migration
```bash
make migration NAME=create_users_table
```

### 3. Edit the Migration Files
Edit `database/migrations/TIMESTAMP_create_users_table.up.sql`:
```sql
BEGIN;
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
COMMIT;
```

Edit `database/migrations/TIMESTAMP_create_users_table.down.sql`:
```sql
BEGIN;
DROP TABLE users;
COMMIT;
```

### 4. Run the Migration
```bash
make migrate-up
```

### 5. Verify
```bash
make migrate-status
```

## 📝 Migration File Format

**Naming**: `{YYYYMMDDHHMMSS}_{name}.{up|down}.sql`

**Content template**:
```sql
-- Migration: migration_name (UP)
-- Created: 2025-12-14T19:24:14+05:00

BEGIN;

-- Your SQL here

COMMIT;
```

## 🔄 Complete Workflow Example

```bash
# 1. Create migration
$ make migration NAME=create_posts_table
✓ Migration created successfully

# 2. Edit the files (see EXAMPLES.md)
# vim database/migrations/TIMESTAMP_create_posts_table.up.sql

# 3. Run migration
$ make migrate-up
✓ Applied: TIMESTAMP_create_posts_table
✓ Migration UP completed

# 4. Check status
$ make migrate-status
TIMESTAMP_create_posts_table              ● Applied

# 5. If something is wrong, rollback
$ make migrate-down
✓ Rolled back: TIMESTAMP_create_posts_table

# 6. Fix and try again
```

## 📚 Documentation Files

- **MIGRATION_GUIDE.md** - Quick reference and common patterns
- **EXAMPLES.md** - Real-world migration examples
- **database/MIGRATIONS.md** - Complete detailed documentation
- **MIGRATION_SETUP_COMPLETE.md** - Setup completion checklist

## ⚙️ How It Works

1. **Create Phase**: Generates migration pair with timestamp
2. **Track Phase**: `migrations` table stores applied migrations
3. **Up Phase**: Executes `.up.sql` files in alphabetical order
4. **Down Phase**: Executes `.down.sql` files in reverse order
5. **Record Phase**: Automatically tracks applied/rolled-back migrations

## 🛡️ Safety Features

✅ **Transactions** - All changes wrapped in BEGIN/COMMIT
✅ **Tracking** - Prevents re-running migrations
✅ **Versioning** - Timestamp ensures consistent ordering
✅ **Reversibility** - DOWN migrations for rollback capability
✅ **Testing** - Can test UP and DOWN in development first

## 🎓 Common Patterns

### Create Table with Constraints
```bash
make migration NAME=create_orders_table
```

### Add Column to Existing Table
```bash
make migration NAME=add_verified_to_users
```

### Create Index
```bash
make migration NAME=create_email_index
```

### Add Foreign Key
```bash
make migration NAME=add_user_id_to_posts
```

### Modify Column
```bash
make migration NAME=change_user_status_type
```

See **EXAMPLES.md** for complete SQL examples.

## 🐛 Troubleshooting

| Problem | Solution |
|---------|----------|
| Database connection error | Check `.env` credentials and DB is running |
| Migration syntax error | Review `.up.sql` file, test SQL in database client |
| Cannot rollback | Check `.down.sql` file, ensure migration was applied |
| Duplicate migration applied | Delete from `migrations` table manually if needed |

## 📊 Check Applied Migrations

```bash
# Using CLI
make migrate-status

# Direct SQL query (PostgreSQL)
SELECT * FROM migrations ORDER BY name;

# Direct SQL query (MySQL)
SELECT * FROM migrations ORDER BY name;
```

## 🔐 Best Practices

✓ **Keep migrations small** - One logical change per migration
✓ **Test DOWN migrations** - Always rollback and re-apply
✓ **Use meaningful names** - Describe what the migration does
✓ **Review before applying** - Check migration files before running
✓ **Test in development** - Don't run untested migrations on production
✓ **Create backups** - Back up database before applying new migrations
✓ **Document complex changes** - Add comments explaining why

## ❌ Avoid

❌ Modifying migration files after they're applied
❌ Large migrations with multiple unrelated changes
❌ Migrations without rollback procedures
❌ Hardcoding environment-specific values
❌ Skipping migrations in sequence
❌ Manual database changes outside migrations

## 📞 Need Help?

1. Check **EXAMPLES.md** for common scenarios
2. Review **database/MIGRATIONS.md** for full documentation
3. Look at example migrations in `database/migrations/`
4. Test your SQL in database client first

## 🎉 You're Ready!

Your migration system is fully functional. Start creating migrations:

```bash
make migration NAME=your_first_table
# Edit the .up.sql and .down.sql files
make migrate-up
make migrate-status
```

Happy migrating! 🚀

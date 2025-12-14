# Database Migrations Guide

This directory contains all database migration files for the application.

## Directory Structure

```
database/
├── migrate.go          # Migration CLI tool
├── migrations/         # Migration files directory
│   ├── 20240101120000_create_users_table.up.sql
│   └── 20240101120000_create_users_table.down.sql
└── MIGRATIONS.md       # This file
```

## Migration Files Naming

Migrations are named with a timestamp and descriptive name:
- **Format**: `YYYYMMDDHHMMSS_migration_name.{up|down}.sql`
- **Example**: `20240101120000_create_users_table.up.sql`

## Using the Migration CLI

### 1. Create a New Migration

```bash
# Using Make
make migration NAME=create_users_table

# Or directly
go run database/migrate.go create create_users_table
```

This creates two files:
- `{timestamp}_create_users_table.up.sql` - Forward migration
- `{timestamp}_create_users_table.down.sql` - Rollback migration

### 2. Run Pending Migrations

```bash
# Run all pending migrations
make migrate-up
go run database/migrate.go up

# Run only N pending migrations
make migrate-up-n N=3
go run database/migrate.go up 3
```

### 3. Rollback Migrations

```bash
# Rollback the last migration
make migrate-down
go run database/migrate.go down

# Rollback N migrations
make migrate-down-n N=2
go run database/migrate.go down 2
```

### 4. Check Migration Status

```bash
# Show which migrations are applied
make migrate-status
go run database/migrate.go status
```

## Migration File Format

Each migration file should follow this format:

**UP Migration** (`20240101120000_create_users_table.up.sql`):
```sql
-- Migration: create_users_table (UP)
-- Created: 2024-01-01T12:00:00Z

BEGIN;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMIT;
```

**DOWN Migration** (`20240101120000_create_users_table.down.sql`):
```sql
-- Migration: create_users_table (DOWN)
-- Created: 2024-01-01T12:00:00Z

BEGIN;

DROP TABLE users;

COMMIT;
```

## Important Notes

1. **Transactions**: Always wrap migrations in `BEGIN;` and `COMMIT;` for safety
2. **Idempotency**: UP migrations should be idempotent (safe to run multiple times)
3. **Rollback Testing**: Always test your DOWN migrations
4. **Database Tracking**: Applied migrations are tracked in the `migrations` table
5. **Filename Format**: Use lowercase with underscores for migration names

## Database Support

- **PostgreSQL**: Fully supported
- **MySQL**: Fully supported
- Other databases may work but require `.env` configuration

## Environment Variables

Set these in your `.env` file:

```env
DB_DRIVER=postgres          # or mysql
DB_HOST=localhost
DB_PORT=5432               # 3306 for MySQL
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=your_database
```

## Common Scenarios

### Create a Table
```bash
make migration NAME=create_posts_table
```

### Add a Column
```bash
make migration NAME=add_status_to_users
```

Edit the `.up.sql` file:
```sql
BEGIN;
ALTER TABLE users ADD COLUMN status VARCHAR(50) DEFAULT 'active';
COMMIT;
```

Edit the `.down.sql` file:
```sql
BEGIN;
ALTER TABLE users DROP COLUMN status;
COMMIT;
```

### Create an Index
```bash
make migration NAME=create_email_index
```

### Add Foreign Key
```bash
make migration NAME=add_user_id_to_posts
```

## Troubleshooting

### Migration Failed
- Check the migration file for syntax errors
- Verify database connectivity
- Check `.env` configuration
- Run `make migrate-status` to see current state

### Locked Migration Table
If the database connection was interrupted:
```sql
-- Manually check the migrations table
SELECT * FROM migrations;

-- If needed, manually fix the record
DELETE FROM migrations WHERE name = 'migration_name';
```

## Best Practices

✓ Write small, focused migrations
✓ Test both UP and DOWN migrations
✓ Use meaningful migration names
✓ Keep SQL clean and readable
✓ Document complex changes with comments
✓ Avoid using application logic in migrations
✓ Test in development environment first

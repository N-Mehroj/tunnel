# 🚀 Quick Start - Database Migrations

## 1️⃣ First Time Setup

### Configure Your Database
Create `.env` file in project root:
```env
DB_DRIVER=postgres          # or mysql
DB_HOST=localhost
DB_PORT=5432               # 3306 for MySQL
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=your_database
```

### Verify Setup Works
```bash
make migrate-status
```

You should see pending migrations.

## 2️⃣ Create Your First Migration

```bash
make migration NAME=create_users_table
```

Output:
```
✓ Migration created successfully:
  UP:   database/migrations/20251214192414_create_users_table.up.sql
  DOWN: database/migrations/20251214192414_create_users_table.down.sql
```

## 3️⃣ Edit Migration Files

### Open UP migration file
`database/migrations/{timestamp}_create_users_table.up.sql`

Add your SQL:
```sql
-- Migration: create_users_table (UP)
-- Created: ...

BEGIN;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);

COMMIT;
```

### Open DOWN migration file
`database/migrations/{timestamp}_create_users_table.down.sql`

Add your rollback SQL:
```sql
-- Migration: create_users_table (DOWN)
-- Created: ...

BEGIN;

DROP INDEX IF EXISTS idx_users_email;
DROP TABLE users;

COMMIT;
```

## 4️⃣ Run the Migration

```bash
make migrate-up
```

Expected output:
```
Running 1 migration(s) UP:

✓ Applied: 20251214192414_create_users_table

✓ Migration UP completed
```

## 5️⃣ Verify It Worked

```bash
make migrate-status
```

Expected output:
```
Migration Status:
------------------------------------------------------------
20251214192414_create_users_table              ● Applied
------------------------------------------------------------
```

## 6️⃣ Continue with More Migrations

Repeat steps 2-5 for each migration:

```bash
make migration NAME=create_posts_table
# Edit UP and DOWN files
make migrate-up
make migrate-status

make migration NAME=add_status_to_users
# Edit UP and DOWN files
make migrate-up
make migrate-status
```

## 🔄 Testing Rollbacks (Recommended)

Always test your DOWN migrations!

```bash
# 1. Run a migration
make migrate-up

# 2. Verify it's applied
make migrate-status

# 3. Test the rollback
make migrate-down

# 4. Verify it's rolled back
make migrate-status

# 5. Re-apply for real
make migrate-up
```

## 📊 Common Operations

### Run Only 1 Migration
```bash
make migrate-up-n N=1
```

### Run Only 3 Pending Migrations
```bash
make migrate-up-n N=3
```

### Rollback Only 1 Migration
```bash
make migrate-down
```

### Rollback Only 2 Migrations
```bash
make migrate-down-n N=2
```

### Rollback All (Careful!)
```bash
# Check how many first
make migrate-status

# Then rollback all
make migrate-down-n N=5  # If 5 are applied
```

## 🆘 Common Issues

### Issue: "Cannot connect to database"
**Solution**: Check `.env` file has correct credentials
```bash
# Test manually
psql -U postgres -d your_database  # PostgreSQL
mysql -u user -p                   # MySQL
```

### Issue: "Migration syntax error"
**Solution**: Review the SQL in `.up.sql` file
```bash
# Test SQL in database client first
# Then fix the migration file
# Then run again
```

### Issue: "Migration applied multiple times"
**Solution**: Migrations are tracked, this shouldn't happen
```bash
# If it did, check the migrations table
SELECT * FROM migrations;
```

## 📚 Learn More

- **EXAMPLES.md** - Real-world examples
- **MIGRATION_GUIDE.md** - All commands and patterns
- **README_MIGRATIONS.md** - Complete documentation

## ✅ Checklist

Before running migrations in production:

- [ ] Test all migrations in development
- [ ] Test all rollbacks in development
- [ ] Create database backup
- [ ] Run migrations on staging first
- [ ] Verify no data loss
- [ ] Get approval if required
- [ ] Then run on production

## 🎉 You're All Set!

Start building your database:
```bash
make migration NAME=your_first_table
```

Happy coding! 🚀

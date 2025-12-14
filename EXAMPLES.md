# 📖 Migration Examples

## Example 1: Create Users Table

### Step 1: Create Migration
```bash
$ make migration NAME=create_users_table
✓ Migration created successfully:
  UP:   database/migrations/20251214192414_create_users_table.up.sql
  DOWN: database/migrations/20251214192414_create_users_table.down.sql
```

### Step 2: Edit UP Migration
File: `database/migrations/20251214192414_create_users_table.up.sql`

```sql
-- Migration: create_users_table (UP)
-- Created: 2025-12-14T19:24:14+05:00

BEGIN;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);

COMMIT;
```

### Step 3: Edit DOWN Migration
File: `database/migrations/20251214192414_create_users_table.down.sql`

```sql
-- Migration: create_users_table (DOWN)
-- Created: 2025-12-14T19:24:14+05:00

BEGIN;

DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;

COMMIT;
```

### Step 4: Run Migration
```bash
$ make migrate-up
Running 1 migration(s) UP:

✓ Applied: 20251214192414_create_users_table

✓ Migration UP completed
```

### Step 5: Verify
```bash
$ make migrate-status
Migration Status:
------------------------------------------------------------
20251214192414_create_users_table              ● Applied
------------------------------------------------------------
```

---

## Example 2: Add Posts Table with Foreign Key

### Step 1: Create Migration
```bash
$ make migration NAME=create_posts_table
```

### Step 2: Edit UP Migration
```sql
-- Migration: create_posts_table (UP)
-- Created: 2025-12-14T19:30:00+05:00

BEGIN;

CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    status VARCHAR(50) DEFAULT 'draft',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_posts_user_id ON posts(user_id);
CREATE INDEX idx_posts_status ON posts(status);

COMMIT;
```

### Step 3: Edit DOWN Migration
```sql
-- Migration: create_posts_table (DOWN)
-- Created: 2025-12-14T19:30:00+05:00

BEGIN;

DROP INDEX IF EXISTS idx_posts_status;
DROP INDEX IF EXISTS idx_posts_user_id;
DROP TABLE IF EXISTS posts;

COMMIT;
```

### Step 4: Run Migration
```bash
$ make migrate-up
Running 1 migration(s) UP:

✓ Applied: 20251214192414_create_posts_table

✓ Migration UP completed
```

---

## Example 3: Add Column to Existing Table

### Step 1: Create Migration
```bash
$ make migration NAME=add_bio_to_users
```

### Step 2: Edit UP Migration
```sql
-- Migration: add_bio_to_users (UP)
-- Created: 2025-12-14T19:45:00+05:00

BEGIN;

ALTER TABLE users ADD COLUMN bio TEXT;
ALTER TABLE users ADD COLUMN website VARCHAR(255);

COMMIT;
```

### Step 3: Edit DOWN Migration
```sql
-- Migration: add_bio_to_users (DOWN)
-- Created: 2025-12-14T19:45:00+05:00

BEGIN;

ALTER TABLE users DROP COLUMN website;
ALTER TABLE users DROP COLUMN bio;

COMMIT;
```

### Step 4: Run
```bash
$ make migrate-up
✓ Applied: 20251214192414_add_bio_to_users
✓ Migration UP completed
```

---

## Example 4: Create Index

### Step 1: Create Migration
```bash
$ make migration NAME=create_post_title_index
```

### Step 2: Edit UP Migration
```sql
-- Migration: create_post_title_index (UP)

BEGIN;

CREATE INDEX idx_posts_title ON posts(title);
CREATE FULLTEXT INDEX idx_posts_content ON posts(content);

COMMIT;
```

### Step 3: Edit DOWN Migration
```sql
-- Migration: create_post_title_index (DOWN)

BEGIN;

DROP INDEX IF EXISTS idx_posts_content ON posts;
DROP INDEX IF EXISTS idx_posts_title;

COMMIT;
```

---

## Example 5: Rollback Operations

### Rollback Last 1 Migration
```bash
$ make migrate-down
Running 1 migration(s) DOWN:

✓ Rolled back: 20251214192414_create_posts_table

✓ Migration DOWN completed
```

### Rollback Last 3 Migrations
```bash
$ make migrate-down-n N=3
Running 3 migration(s) DOWN:

✓ Rolled back: 20251214192414_add_bio_to_users
✓ Rolled back: 20251214192414_create_posts_table
✓ Rolled back: 20251214192414_create_users_table

✓ Migration DOWN completed
```

### Check Status After Rollback
```bash
$ make migrate-status
Migration Status:
------------------------------------------------------------
20251214192414_create_users_table              ○ Pending
20251214192414_create_posts_table              ○ Pending
20251214192414_add_bio_to_users                ○ Pending
------------------------------------------------------------
```

---

## Example 6: Multiple Pending Migrations

### Create Multiple Migrations
```bash
$ make migration NAME=create_comments_table
✓ Migration created: database/migrations/20251214200000_create_comments_table.up.sql

$ make migration NAME=create_tags_table
✓ Migration created: database/migrations/20251214200100_create_tags_table.up.sql

$ make migration NAME=add_tag_column_to_posts
✓ Migration created: database/migrations/20251214200200_add_tag_column_to_posts.up.sql
```

### Check Status
```bash
$ make migrate-status
Migration Status:
------------------------------------------------------------
20251214192414_create_users_table              ● Applied
20251214192414_create_posts_table              ● Applied
20251214200000_create_comments_table           ○ Pending
20251214200100_create_tags_table               ○ Pending
20251214200200_add_tag_column_to_posts         ○ Pending
------------------------------------------------------------
```

### Run Only 2 of 3 Pending
```bash
$ make migrate-up-n N=2
Running 2 migration(s) UP:

✓ Applied: 20251214200000_create_comments_table
✓ Applied: 20251214200100_create_tags_table

✓ Migration UP completed
```

### Check Status Again
```bash
$ make migrate-status
Migration Status:
------------------------------------------------------------
20251214192414_create_users_table              ● Applied
20251214192414_create_posts_table              ● Applied
20251214200000_create_comments_table           ● Applied
20251214200100_create_tags_table               ● Applied
20251214200200_add_tag_column_to_posts         ○ Pending
------------------------------------------------------------
```

### Run Remaining
```bash
$ make migrate-up
Running 1 migration(s) UP:

✓ Applied: 20251214200200_add_tag_column_to_posts

✓ Migration UP completed
```

---

## Common SQL Patterns

### Create Table with All Constraints
```sql
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    order_number VARCHAR(50) UNIQUE NOT NULL,
    total_amount DECIMAL(10, 2) NOT NULL DEFAULT 0,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CHECK (total_amount >= 0)
);
```

### Add Constraint
```sql
ALTER TABLE posts ADD CONSTRAINT fk_user 
FOREIGN KEY (user_id) REFERENCES users(id);
```

### Modify Column Type
```sql
ALTER TABLE users ALTER COLUMN bio TYPE VARCHAR(1000);
```

### Set Default Value
```sql
ALTER TABLE posts ALTER COLUMN status SET DEFAULT 'published';
```

### Drop Constraint
```sql
ALTER TABLE posts DROP CONSTRAINT fk_user;
```

---

## Tips & Tricks

### Always Test DOWN First
Before running migrations in production, test rollback:
```bash
# In staging
make migrate-up
make migrate-down
make migrate-up
# If this works, safe for production
```

### Use Comments
```sql
BEGIN;

-- Drop old index before modifying column
DROP INDEX IF EXISTS idx_users_status;

-- Expand status field to support new values
ALTER TABLE users ALTER COLUMN status VARCHAR(50);

-- Create new more specific index
CREATE INDEX idx_users_active ON users(status) 
WHERE status = 'active';

COMMIT;
```

### Meaningful Names
```
✅ Good names:
- create_users_table
- add_email_verification_token_to_users
- create_index_on_posts_user_id
- add_foreign_key_user_to_posts

❌ Bad names:
- update_database
- fix_stuff
- temp_migration
- migration1
```

---

For more information, see:
- `database/MIGRATIONS.md` - Complete documentation
- `MIGRATION_GUIDE.md` - Quick reference

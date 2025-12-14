#!/bin/bash

# Recovery script to fix the migration records in the database
# The tables were created but not recorded due to a SQL placeholder bug (now fixed)

# This script helps you manually record the migrations that were applied
# You need to run these commands directly in your database

echo "╔═══════════════════════════════════════════════════════════════╗"
echo "║     MIGRATION RECORDS RECOVERY SCRIPT                         ║"
echo "╚═══════════════════════════════════════════════════════════════╝"
echo ""
echo "The tables were created but migration records weren't stored."
echo "Follow these steps to fix it:"
echo ""
echo "─────────────────────────────────────────────────────────────────"
echo "Option 1: FOR POSTGRESQL"
echo "─────────────────────────────────────────────────────────────────"
echo ""
echo "Run these commands in psql:"
echo ""
echo "  psql -U \$DB_USER -d \$DB_NAME"
echo ""
echo "Then execute:"
echo ""
cat << 'PGEOF'
  INSERT INTO migrations (name) VALUES ('20240101000000_create_users_table') 
    ON CONFLICT (name) DO NOTHING;
  INSERT INTO migrations (name) VALUES ('20251214192414_create_auth_table') 
    ON CONFLICT (name) DO NOTHING;
  SELECT * FROM migrations;
PGEOF
echo ""
echo "─────────────────────────────────────────────────────────────────"
echo "Option 2: FOR MYSQL"
echo "─────────────────────────────────────────────────────────────────"
echo ""
echo "Run these commands in mysql:"
echo ""
echo "  mysql -u \$DB_USER -p"
echo "  USE \$DB_NAME;"
echo ""
echo "Then execute:"
echo ""
cat << 'MYSQLEOF'
  INSERT IGNORE INTO migrations (name) VALUES ('20240101000000_create_users_table');
  INSERT IGNORE INTO migrations (name) VALUES ('20251214192414_create_auth_table');
  SELECT * FROM migrations;
MYSQLEOF
echo ""
echo "─────────────────────────────────────────────────────────────────"
echo "After recording the migrations, verify with:"
echo "  make migrate-status"
echo ""
echo "The bug in database/migrate.go has been fixed!"
echo "─────────────────────────────────────────────────────────────────"

#!/bin/bash
# Fix migration records that failed to be recorded due to SQL placeholder bug

# This script helps manually record the migrations that were applied but not recorded
# in the migrations table due to the SQL placeholder syntax error that was fixed.

echo "Fixing migration records..."

# For PostgreSQL (update the connection details as needed)
# psql -U postgres -d your_database << EOF
# INSERT INTO migrations (name) VALUES ('20240101000000_create_users_table');
# INSERT INTO migrations (name) VALUES ('20251214192414_create_auth_table');
# EOF

# For MySQL (update the connection details as needed)
# mysql -u user -p database << EOF
# INSERT INTO migrations (name) VALUES ('20240101000000_create_users_table');
# INSERT INTO migrations (name) VALUES ('20251214192414_create_auth_table');
# EOF

echo ""
echo "Option 1: Using psql (PostgreSQL)"
echo "  psql -U \$DB_USER -d \$DB_NAME << EOF"
echo "  INSERT INTO migrations (name) VALUES ('20240101000000_create_users_table');"
echo "  INSERT INTO migrations (name) VALUES ('20251214192414_create_auth_table');"
echo "  EOF"
echo ""

echo "Option 2: Using mysql (MySQL)"
echo "  mysql -u \$DB_USER -p \$DB_NAME << EOF"
echo "  INSERT INTO migrations (name) VALUES ('20240101000000_create_users_table');"
echo "  INSERT INTO migrations (name) VALUES ('20251214192414_create_auth_table');"
echo "  EOF"
echo ""

echo "The bug has been fixed in database/migrate.go"
echo "Future migrations will record correctly!"

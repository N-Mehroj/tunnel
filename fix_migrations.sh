#!/bin/bash

# Auto-recovery script for migration records
# This script attempts to automatically fix the migration records

set -e

echo "╔════════════════════════════════════════════════════════════════╗"
echo "║     Automatic Migration Records Recovery Script                ║"
echo "╚════════════════════════════════════════════════════════════════╝"
echo ""

# Load environment
if [ ! -f .env ]; then
    echo "❌ Error: .env file not found!"
    echo "Please create .env with your database configuration"
    exit 1
fi

source .env

echo "Database configuration loaded from .env"
echo "Driver: $DB_DRIVER"
echo "Host: $DB_HOST"
echo "Database: $DB_NAME"
echo ""

# Function to fix PostgreSQL
fix_postgres() {
    echo "🔧 Fixing PostgreSQL migration records..."
    
    PGPASSWORD=$DB_PASSWORD psql \
        -h $DB_HOST \
        -p ${DB_PORT:-5432} \
        -U $DB_USER \
        -d $DB_NAME \
        << EOF
INSERT INTO migrations (name) VALUES ('20240101000000_create_users_table') 
    ON CONFLICT (name) DO NOTHING;

INSERT INTO migrations (name) VALUES ('20251214192414_create_auth_table') 
    ON CONFLICT (name) DO NOTHING;

-- Verify
SELECT '✓ Current migration records:' as status;
SELECT * FROM migrations ORDER BY name;
EOF
    
    if [ $? -eq 0 ]; then
        echo "✅ PostgreSQL migration records fixed!"
        return 0
    else
        echo "❌ Failed to fix migration records"
        return 1
    fi
}

# Function to fix MySQL
fix_mysql() {
    echo "🔧 Fixing MySQL migration records..."
    
    mysql -h $DB_HOST \
        -P ${DB_PORT:-3306} \
        -u $DB_USER \
        -p$DB_PASSWORD \
        $DB_NAME \
        << EOF
INSERT IGNORE INTO migrations (name) VALUES ('20240101000000_create_users_table');
INSERT IGNORE INTO migrations (name) VALUES ('20251214192414_create_auth_table');

-- Verify
SELECT '✓ Current migration records:' as status;
SELECT * FROM migrations ORDER BY name;
EOF
    
    if [ $? -eq 0 ]; then
        echo "✅ MySQL migration records fixed!"
        return 0
    else
        echo "❌ Failed to fix migration records"
        return 1
    fi
}

# Main logic
case $DB_DRIVER in
    postgres)
        fix_postgres
        ;;
    mysql)
        fix_mysql
        ;;
    *)
        echo "❌ Unknown database driver: $DB_DRIVER"
        echo "Supported: postgres, mysql"
        exit 1
        ;;
esac

if [ $? -eq 0 ]; then
    echo ""
    echo "╔════════════════════════════════════════════════════════════════╗"
    echo "║                    ✅ RECOVERY COMPLETE!                        ║"
    echo "╚════════════════════════════════════════════════════════════════╝"
    echo ""
    echo "Verify with: make migrate-status"
    exit 0
else
    echo ""
    echo "❌ Recovery failed. Please refer to RECOVERY_GUIDE.md"
    exit 1
fi

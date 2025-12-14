package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const migrationsDir = "database/migrations"

func main() {
	args := os.Args
	if len(args) < 2 {
		printUsage()
		return
	}

	action := args[1]

	switch action {
	case "create":
		if len(args) < 3 {
			fmt.Println("Usage: go run database/migrate.go create <migration_name>")
			return
		}
		name := args[2]
		createMigration(name)
	case "up":
		count := 0
		if len(args) > 2 {
			fmt.Sscanf(args[2], "%d", &count)
		}
		migrateUp(count)
	case "down":
		count := 1
		if len(args) > 2 {
			fmt.Sscanf(args[2], "%d", &count)
		}
		migrateDown(count)
	case "status":
		migrateStatus()
	default:
		fmt.Println("Unknown action:", action)
		printUsage()
	}
}

func printUsage() {
	fmt.Println(`
			Migration CLI - Manage database migrations

			Usage:
			go run database/migrate.go create <migration_name>   Create a new migration
			go run database/migrate.go up [number]               Run pending migrations (or specify count)
			go run database/migrate.go down [number]             Rollback migrations (default 1)
			go run database/migrate.go status                    Show migration status

			Examples:
			go run database/migrate.go create create_users_table
			go run database/migrate.go up
			go run database/migrate.go up 3
			go run database/migrate.go down
			go run database/migrate.go down 2
			go run database/migrate.go status
	`)
}

func createMigration(name string) {
	if err := os.MkdirAll(migrationsDir, 0755); err != nil {
		log.Fatalf("Failed to create migrations directory: %v", err)
	}

	timestamp := time.Now().Format("20060102150405")
	upFileName := filepath.Join(migrationsDir, fmt.Sprintf("%s_%s.up.sql", timestamp, name))
	downFileName := filepath.Join(migrationsDir, fmt.Sprintf("%s_%s.down.sql", timestamp, name))

	upTemplate := fmt.Sprintf(`-- Migration: %s (UP)
-- Created: %s

-- Write your UP migration SQL here

`, name, time.Now().Format(time.RFC3339))

	downTemplate := fmt.Sprintf(`-- Migration: %s (DOWN)
-- Created: %s

-- Write your DOWN migration SQL here

`, name, time.Now().Format(time.RFC3339))

	if err := os.WriteFile(upFileName, []byte(upTemplate), 0644); err != nil {
		log.Fatalf("Failed to create UP migration file: %v", err)
	}

	if err := os.WriteFile(downFileName, []byte(downTemplate), 0644); err != nil {
		log.Fatalf("Failed to create DOWN migration file: %v", err)
	}

	fmt.Printf("✓ Migration created successfully:\n")
	fmt.Printf("  UP:   %s\n", upFileName)
	fmt.Printf("  DOWN: %s\n", downFileName)
}

func getDB() *sql.DB {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	dbDriver := os.Getenv("DB_DRIVER")
	if dbDriver == "" {
		dbDriver = "postgres"
	}

	var db *sql.DB
	var err error

	switch dbDriver {
	case "postgres":
		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
		)
		db, err = sql.Open("postgres", dsn)
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
		)
		db, err = sql.Open("mysql", dsn)
	default:
		log.Fatalf("Unknown DB_DRIVER: %s", dbDriver)
	}

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	return db
}

func initMigrationsTable(db *sql.DB) error {
	dbDriver := os.Getenv("DB_DRIVER")
	if dbDriver == "" {
		dbDriver = "postgres"
	}

	var createTableSQL string
	if dbDriver == "postgres" {
		createTableSQL = `
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL UNIQUE,
			executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		`
	} else if dbDriver == "mysql" {
		createTableSQL = `
		CREATE TABLE IF NOT EXISTS migrations (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL UNIQUE,
			executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		`
	}

	_, err := db.Exec(createTableSQL)
	return err
}

func getMigrationFiles() []string {
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}
		}
		log.Fatalf("Failed to read migrations directory: %v", err)
	}

	var migrations []string
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".up.sql") {
			migrations = append(migrations, strings.TrimSuffix(f.Name(), ".up.sql"))
		}
	}

	sort.Strings(migrations)
	return migrations
}

func getExecutedMigrations(db *sql.DB) map[string]bool {
	executed := make(map[string]bool)
	rows, err := db.Query("SELECT name FROM migrations ORDER BY name")
	if err != nil {
		return executed
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Println("Warning: Failed to scan migration:", err)
			continue
		}
		executed[name] = true
	}

	return executed
}

func migrateUp(count int) {
	db := getDB()
	defer db.Close()

	if err := initMigrationsTable(db); err != nil {
		log.Fatalf("Failed to initialize migrations table: %v", err)
	}

	allMigrations := getMigrationFiles()
	executed := getExecutedMigrations(db)

	var pending []string
	for _, m := range allMigrations {
		if !executed[m] {
			pending = append(pending, m)
		}
	}

	if len(pending) == 0 {
		fmt.Println("✓ No pending migrations")
		return
	}

	if count > 0 && count < len(pending) {
		pending = pending[:count]
	}

	fmt.Printf("Running %d migration(s) UP:\n\n", len(pending))

	dbDriver := os.Getenv("DB_DRIVER")
	if dbDriver == "" {
		dbDriver = "postgres"
	}

	for _, migrationName := range pending {
		upFile := filepath.Join(migrationsDir, migrationName+".up.sql")
		content, err := os.ReadFile(upFile)
		if err != nil {
			log.Printf("✗ Failed to read %s: %v\n", upFile, err)
			continue
		}

		tx, err := db.Begin()
		if err != nil {
			log.Printf("✗ Failed to begin transaction for %s: %v\n", migrationName, err)
			continue
		}

		if _, err := tx.Exec(string(content)); err != nil {
			tx.Rollback()
			log.Printf("✗ Failed to execute %s: %v\n", migrationName, err)
			continue
		}

		var insertSQL string
		if dbDriver == "postgres" {
			insertSQL = "INSERT INTO migrations (name) VALUES ($1)"
		} else {
			insertSQL = "INSERT INTO migrations (name) VALUES (?)"
		}

		if _, err := tx.Exec(insertSQL, migrationName); err != nil {
			tx.Rollback()
			log.Printf("✗ Failed to record migration %s: %v\n", migrationName, err)
			continue
		}

		if err := tx.Commit(); err != nil {
			log.Printf("✗ Failed to commit migration %s: %v\n", migrationName, err)
			continue
		}

		fmt.Printf("✓ Applied: %s\n", migrationName)
	}

	fmt.Println("\n✓ Migration UP completed")
}

func migrateDown(count int) {
	db := getDB()
	defer db.Close()

	executed := getExecutedMigrations(db)

	var execOrder []string
	for name := range executed {
		execOrder = append(execOrder, name)
	}

	sort.Sort(sort.Reverse(sort.StringSlice(execOrder)))

	if len(execOrder) == 0 {
		fmt.Println("✓ No migrations to rollback")
		return
	}

	if count > len(execOrder) {
		count = len(execOrder)
	}

	toRollback := execOrder[:count]

	fmt.Printf("Rolling back %d migration(s) DOWN:\n\n", len(toRollback))

	dbDriver := os.Getenv("DB_DRIVER")
	if dbDriver == "" {
		dbDriver = "postgres"
	}

	for _, migrationName := range toRollback {
		downFile := filepath.Join(migrationsDir, migrationName+".down.sql")
		content, err := os.ReadFile(downFile)
		if err != nil {
			log.Printf("✗ Failed to read %s: %v\n", downFile, err)
			continue
		}

		tx, err := db.Begin()
		if err != nil {
			log.Printf("✗ Failed to begin transaction for %s: %v\n", migrationName, err)
			continue
		}

		if _, err := tx.Exec(string(content)); err != nil {
			tx.Rollback()
			log.Printf("✗ Failed to execute rollback for %s: %v\n", migrationName, err)
			continue
		}

		var deleteSQL string
		if dbDriver == "postgres" {
			deleteSQL = "DELETE FROM migrations WHERE name = $1"
		} else {
			deleteSQL = "DELETE FROM migrations WHERE name = ?"
		}

		if _, err := tx.Exec(deleteSQL, migrationName); err != nil {
			tx.Rollback()
			log.Printf("✗ Failed to remove migration record %s: %v\n", migrationName, err)
			continue
		}

		if err := tx.Commit(); err != nil {
			log.Printf("✗ Failed to commit rollback for %s: %v\n", migrationName, err)
			continue
		}

		fmt.Printf("✓ Rolled back: %s\n", migrationName)
	}

	fmt.Println("\n✓ Migration DOWN completed")
}

func migrateStatus() {
	db := getDB()
	defer db.Close()

	if err := initMigrationsTable(db); err != nil {
		log.Printf("Warning: Could not initialize migrations table: %v\n", err)
	}

	allMigrations := getMigrationFiles()
	executed := getExecutedMigrations(db)

	fmt.Println("\nMigration Status:")
	fmt.Println(strings.Repeat("-", 60))

	if len(allMigrations) == 0 {
		fmt.Println("No migrations found")
		return
	}

	for _, m := range allMigrations {
		status := "○ Pending"
		if executed[m] {
			status = "● Applied"
		}
		fmt.Printf("%-40s %s\n", m, status)
	}

	fmt.Println(strings.Repeat("-", 60))
}

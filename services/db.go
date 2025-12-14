package services

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    dbDriver := os.Getenv("DB_DRIVER")
    dbName := os.Getenv("DB_NAME")
    user := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASSWORD")
    host := os.Getenv("DB_HOST")
    port := os.Getenv("DB_PORT")
    
    var defaultDBName string
    var driverName string
    var dialect gorm.Dialector

    switch dbDriver {
    case "postgres":
        driverName = "postgres"
        defaultDBName = "postgres" 
    case "mysql":
        driverName = "mysql"
        defaultDBName = "mysql" 
    default:
        log.Fatalf("Unknown DB_DRIVER: %s. Only 'postgres' or 'mysql' are accepted.", dbDriver)
    }

    var dsn string
    
    if dbDriver == "postgres" {
        dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
            host, user, password, defaultDBName, port)
    } else if dbDriver == "mysql" {
        dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
            user, password, host, port, defaultDBName)
    }
    
    sqlDB, err := sql.Open(driverName, dsn)
    if err != nil {
        log.Fatalf("Error connecting to the database server: %v", err)
    }
    defer sqlDB.Close()

    _, err = sqlDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
    if err == nil {
        fmt.Println("Database created:", dbName)
    }
    time.Sleep(1 * time.Second)

    var gormDSN string

    if dbDriver == "postgres" {
        gormDSN = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tashkent",
            host, user, password, dbName, port)
        dialect = postgres.Open(gormDSN)
    } else if dbDriver == "mysql" {
        gormDSN = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
            user, password, host, port, dbName)
        dialect = mysql.Open(gormDSN)
    }

    database, err := gorm.Open(dialect, &gorm.Config{})
    if err != nil {
        log.Fatalf("Error connecting with GORM: %v", err)
    }

    DB = database
    // fmt.Printf("%s connected to the database successfully.\n", dbDriver)
}
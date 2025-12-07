package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	var err error
	var connStr string
	
	// Check if DATABASE_URL is set (used by Railway and other PaaS)
	databaseURL := os.Getenv("DATABASE_URL")
	
	if databaseURL != "" {
		// Use DATABASE_URL directly
		connStr = databaseURL
	} else {
		// Get database configuration from environment
		dbHost := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")
		dbUser := os.Getenv("DB_USER")
		dbPassword := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")
		dbSSLMode := os.Getenv("DB_SSLMODE")

		// Default values
		if dbHost == "" {
			dbHost = "localhost"
		}
		if dbPort == "" {
			dbPort = "5432"
		}
		if dbSSLMode == "" {
			dbSSLMode = "disable"
		}

		// Create connection string
		connStr = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode,
		)
	}

	// Open database connection
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Test connection
	if err = DB.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Database connected successfully")
}

func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("Database connection closed")
	}
}

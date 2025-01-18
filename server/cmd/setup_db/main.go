package main

import (
	"database/sql"
	"log"
	"os"
	"server/crud"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// Get database connection details from environment variables
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL environment variable is not set")
	}

	// Connect to PostgreSQL using the URL from .env
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	log.Println("Successfully connected to database")

	// Create tables
	if err := crud.CreateTables(db); err != nil {
		log.Fatal("Failed to create tables:", err)
	}
	log.Println("Successfully created tables")

	// Insert all user data
	if err := crud.InsertAllUserData(db); err != nil {
		log.Fatal("Failed to insert user data:", err)
	}
	log.Println("Successfully inserted all user data")

	// Verify data
	for _, accountID := range []string{"1234567890", "1234567891", "1234567892", "1234567893"} {
		transactions, err := crud.GetTransactions(db, accountID)
		if err != nil {
			log.Printf("Failed to get transactions for account %s: %v", accountID, err)
			continue
		}
		log.Printf("Found %d transactions for account %s", len(transactions), accountID)
	}
} 
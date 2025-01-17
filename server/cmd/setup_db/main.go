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

	// Insert Jane's data
	if err := crud.InsertJaneData(db); err != nil {
		log.Fatal("Failed to insert Jane's data:", err)
	}
	log.Println("Successfully inserted Jane's data")

	// Verify data
	transactions, err := crud.GetTransactions(db, "1234567891")
	if err != nil {
		log.Fatal("Failed to get transactions:", err)
	}
	log.Printf("Found %d transactions for Jane", len(transactions))
} 
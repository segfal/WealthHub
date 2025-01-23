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
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// Get database URL from environment
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL environment variable is not set")
	}

	// Connect to the database using the URL from .env
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

	// Get all transactions to verify setup
	transactions, err := crud.GetTransactions(db, "test123")
	if err != nil {
		log.Fatal("Failed to get transactions:", err)
	}
	log.Printf("Found %d transactions", len(transactions))
} 
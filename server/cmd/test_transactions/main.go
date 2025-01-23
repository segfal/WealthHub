package main

import (
	"database/sql"
	"log"
	"os"
	"server/crud"
	"server/types"
	"time"

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

	// Create test user
	testUser := &types.User{
		AccountID:     "test123",
		AccountName:   "Test Account",
		AccountType:   "checking",
		AccountNumber: "1234567890",
		Balance: types.UserBalance{
			Current:   1000.00,
			Available: 900.00,
			Currency:  "USD",
		},
		OwnerName: "Test User",
		BankDetails: types.UserBankDetails{
			BankName:      "Test Bank",
			RoutingNumber: "987654321",
			Branch:        "Test Branch",
		},
	}

	if err := crud.CreateUser(db, testUser); err != nil {
		log.Fatal("Failed to create test user:", err)
	}

	// Create test transaction
	testTransaction := &types.Transaction{
		TransactionID: "txn123",
		AccountID:     "test123",
		Date:         time.Now(),
		Amount:       100.00,
		Category:     "groceries",
		Merchant:     "Test Store",
		Location:     "Test Location",
	}

	if err := crud.InsertTransaction(db, testTransaction); err != nil {
		log.Fatal("Failed to insert test transaction:", err)
	}

	// Get transactions to verify
	transactions, err := crud.GetTransactions(db, "test123")
	if err != nil {
		log.Fatal("Failed to get transactions:", err)
	}
	log.Printf("Found %d test transactions", len(transactions))
} 
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"server/crud"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
	}

	// Connect to database
	db, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test accounts
	accounts := []string{
		"1234567890", // John
		"1234567891", // Jane
		"1234567892", // Jake
		"1234567893", // Jill
	}

	// Get and display sample transactions for each account
	for _, accountID := range accounts {
		transactions, err := crud.GetTransactions(db, accountID)
		if err != nil {
			log.Printf("Error getting transactions for account %s: %v", accountID, err)
			continue
		}

		fmt.Printf("\nAccount %s Transactions:\n", accountID)
		fmt.Printf("Total transactions: %d\n", len(transactions))
		
		// Display first 5 transactions
		for i := 0; i < 5 && i < len(transactions); i++ {
			t := transactions[i]
			fmt.Printf("Transaction %d:\n", i+1)
			fmt.Printf("  ID: %s (Original ID: %s, User: %s)\n", 
				t.UserPrefix+"_"+t.TransactionID, t.TransactionID, t.UserPrefix)
			fmt.Printf("  Date: %s\n", t.Date.Format("2006-01-02"))
			fmt.Printf("  Amount: $%.2f\n", t.Amount)
			fmt.Printf("  Category: %s\n", t.Category)
			fmt.Printf("  Merchant: %s\n", t.Merchant)
			fmt.Printf("  Location: %s\n", t.Location)
			fmt.Println()
		}
	}
} 
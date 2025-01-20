package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"server/types"
	"time"

	_ "github.com/lib/pq"
)

func createTables(db *sql.DB) error {
	// Drop existing tables
	dropTables := `
		DROP TABLE IF EXISTS transactions;
		DROP TABLE IF EXISTS bank_details;
		DROP TABLE IF EXISTS users;`
	
	if _, err := db.Exec(dropTables); err != nil {
		return fmt.Errorf("failed to drop tables: %w", err)
	}

	// Create users table
	createUsers := `
		CREATE TABLE users (
			account_id VARCHAR(20) PRIMARY KEY,
			account_name VARCHAR(50),
			account_type VARCHAR(20),
			account_number VARCHAR(20),
			balance_current DECIMAL(10, 2),
			balance_available DECIMAL(10, 2),
			balance_currency VARCHAR(3),
			owner_name VARCHAR(50)
		)`
	
	if _, err := db.Exec(createUsers); err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	// Create bank_details table
	createBankDetails := `
		CREATE TABLE bank_details (
			account_id VARCHAR(20) PRIMARY KEY REFERENCES users(account_id),
			bank_name VARCHAR(50),
			routing_number VARCHAR(20),
			branch VARCHAR(100)
		)`
	
	if _, err := db.Exec(createBankDetails); err != nil {
		return fmt.Errorf("failed to create bank_details table: %w", err)
	}

	// Create transactions table
	createTransactions := `
		CREATE TABLE transactions (
			transaction_id VARCHAR(20) PRIMARY KEY,
			account_id VARCHAR(20) REFERENCES users(account_id),
			date TIMESTAMP,
			amount DECIMAL(10, 2),
			category VARCHAR(50),
			merchant VARCHAR(50),
			location VARCHAR(100)
		)`
	
	if _, err := db.Exec(createTransactions); err != nil {
		return fmt.Errorf("failed to create transactions table: %w", err)
	}

	return nil
}

// InsertJaneData reads JaneDoe.json and inserts the data into the database
func InsertJaneData(db *sql.DB) error {
	// Read JSON file
	data, err := os.ReadFile("JaneDoe.json")
	if err != nil {
		return fmt.Errorf("failed to read JSON file: %w", err)
	}

	// Parse JSON
	var jsonData struct {
		Account struct {
			AccountID     int    `json:"account_id"`
			AccountName   string `json:"account_name"`
			AccountType   string `json:"account_type"`
			AccountNumber string `json:"account_number"`
			Balance       struct {
				Current   float64 `json:"current"`
				Available float64 `json:"available"`
				Currency  string  `json:"currency"`
			} `json:"balance"`
			OwnerName   string `json:"owner_name"`
			BankDetails struct {
				BankName      string `json:"bank_name"`
				RoutingNumber string `json:"routing_number"`
				Branch        string `json:"branch"`
			} `json:"bank_details"`
			Transactions []struct {
				TransactionID string  `json:"transaction_id"`
				Date         string  `json:"date"`
				Amount       float64 `json:"amount"`
				Category     string  `json:"category"`
				Merchant     string  `json:"merchant"`
				Location     string  `json:"location"`
			} `json:"transactions"`
		} `json:"account"`
	}

	if err := json.Unmarshal(data, &jsonData); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert user
	user := &types.User{
		AccountID:     fmt.Sprintf("%d", jsonData.Account.AccountID),
		AccountName:   jsonData.Account.AccountName,
		AccountType:   jsonData.Account.AccountType,
		AccountNumber: jsonData.Account.AccountNumber,
		Balance: types.UserBalance{
			Current:   jsonData.Account.Balance.Current,
			Available: jsonData.Account.Balance.Available,
			Currency:  jsonData.Account.Balance.Currency,
		},
		OwnerName: jsonData.Account.OwnerName,
		BankDetails: types.UserBankDetails{
			BankName:      jsonData.Account.BankDetails.BankName,
			RoutingNumber: jsonData.Account.BankDetails.RoutingNumber,
			Branch:        jsonData.Account.BankDetails.Branch,
		},
	}

	if err := createUserTx(tx, user); err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	// Insert transactions
	for _, t := range jsonData.Account.Transactions {
		date, err := time.Parse("2006-01-02", t.Date)
		if err != nil {
			return fmt.Errorf("failed to parse date %s: %w", t.Date, err)
		}

		transaction := &types.Transaction{
			TransactionID: t.TransactionID,
			AccountID:    fmt.Sprintf("%d", jsonData.Account.AccountID),
			Date:        date,
			Amount:      t.Amount,
			Category:    t.Category,
			Merchant:    t.Merchant,
			Location:    t.Location,
		}

		if err := insertTransactionTx(tx, transaction); err != nil {
			return fmt.Errorf("failed to insert transaction %s: %w", t.TransactionID, err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// createUserTx creates a user within a transaction
func createUserTx(tx *sql.Tx, user *types.User) error {
	query := `
		INSERT INTO users (
			account_id, account_name, account_type, account_number, 
			owner_name, balance_current, balance_available, balance_currency
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := tx.Exec(query,
		user.AccountID,
		user.AccountName,
		user.AccountType,
		user.AccountNumber,
		user.OwnerName,
		user.Balance.Current,
		user.Balance.Available,
		user.Balance.Currency,
	)
	if err != nil {
		return err
	}

	bankQuery := `
		INSERT INTO bank_details (
			account_id, bank_name, routing_number, branch
		) VALUES ($1, $2, $3, $4)`

	_, err = tx.Exec(bankQuery,
		user.AccountID,
		user.BankDetails.BankName,
		user.BankDetails.RoutingNumber,
		user.BankDetails.Branch,
	)
	return err
}

// insertTransactionTx inserts a transaction within a transaction
func insertTransactionTx(tx *sql.Tx, transaction *types.Transaction) error {
	query := `
		INSERT INTO transactions (
			transaction_id, account_id, date, amount, category, merchant, location
		) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	
	_, err := tx.Exec(query,
		transaction.TransactionID,
		transaction.AccountID,
		transaction.Date,
		transaction.Amount,
		transaction.Category,
		transaction.Merchant,
		transaction.Location,
	)
	return err
}

func createUser(db *sql.DB, user *types.User) error {
	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if err := createUserTx(tx, user); err != nil {
		return err
	}

	return tx.Commit()
}

func getTransactions(db *sql.DB, accountID string) ([]types.Transaction, error) { 
	query := ` 
		SELECT transaction_id, account_id, date, amount, category, merchant, location 
		FROM transactions 
		WHERE account_id = $1
		ORDER BY date DESC`
	
	rows, err := db.Query(query, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []types.Transaction
	for rows.Next() {
		var t types.Transaction
		if err := rows.Scan(
			&t.TransactionID,
			&t.AccountID,
			&t.Date,
			&t.Amount,
			&t.Category,
			&t.Merchant,
			&t.Location,
		); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

func insertTransaction(db *sql.DB, transaction *types.Transaction) error {
	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if err := insertTransactionTx(tx, transaction); err != nil {
		return err
	}

	return tx.Commit()
}

// SetupDatabase initializes the database and loads initial data
func SetupDatabase(db *sql.DB) error {
	// Create tables
	if err := createTables(db); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}
	log.Println("Successfully created tables")

	// Insert Jane's data
	if err := InsertJaneData(db); err != nil {
		return fmt.Errorf("failed to insert Jane's data: %w", err)
	}
	log.Println("Successfully inserted Jane's data")

	// Get transactions to verify
	transactions, err := getTransactions(db, "1234567891")
	if err != nil {
		return fmt.Errorf("failed to get transactions: %w", err)
	}
	log.Printf("Found %d transactions for Jane", len(transactions))

	return nil
}
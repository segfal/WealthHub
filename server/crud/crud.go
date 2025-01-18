package crud

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"server/types"
	"strings"
	"time"
)

// CreateTables creates all necessary database tables
func CreateTables(db *sql.DB) error {
	// Drop existing tables in correct order (dependent tables first)
	dropTables := `
		DROP TABLE IF EXISTS transactions CASCADE;
		DROP TABLE IF EXISTS users CASCADE;`
	
	if _, err := db.Exec(dropTables); err != nil {
		return fmt.Errorf("failed to drop tables: %w", err)
	}

	// Create users table with bank details included
	createUsers := `
		CREATE TABLE users (
			account_id VARCHAR(20) PRIMARY KEY,
			account_name VARCHAR(50),
			account_type VARCHAR(20),
			account_number VARCHAR(20),
			balance_current DECIMAL(10, 2),
			balance_available DECIMAL(10, 2),
			balance_currency VARCHAR(3),
			owner_name VARCHAR(50),
			bank_name VARCHAR(50),
			routing_number VARCHAR(20),
			branch VARCHAR(100)
		)`
	
	if err := db.QueryRow(createUsers).Err(); err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
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
	
	if err := db.QueryRow(createTransactions).Err(); err != nil {
		return fmt.Errorf("failed to create transactions table: %w", err)
	}

	return nil
}

// InsertUserData reads a user's JSON file and inserts the data into the database
func InsertUserData(db *sql.DB, filename string) error {
	// Read JSON file
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read JSON file %s: %w", filename, err)
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
				AccountID    int     `json:"account_id"`
				Date        string  `json:"date"`
				Amount      float64 `json:"amount"`
				Category    string  `json:"category"`
				Merchant    string  `json:"merchant"`
				Location    string  `json:"location"`
				Type       string  `json:"type"`
				Status     string  `json:"status"`
				Timestamp  string  `json:"timestamp"`
			} `json:"transactions"`
		} `json:"account"`
	}

	if err := json.Unmarshal(data, &jsonData); err != nil {
		return fmt.Errorf("failed to parse JSON from %s: %w", filename, err)
	}

	fmt.Printf("Processing %s: Found %d transactions\n", filename, len(jsonData.Account.Transactions))

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert user with bank details using ON CONFLICT DO NOTHING
	userQuery := `
		INSERT INTO users (
			account_id, account_name, account_type, account_number, 
			owner_name, balance_current, balance_available, balance_currency,
			bank_name, routing_number, branch
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (account_id) DO UPDATE SET
			account_name = EXCLUDED.account_name,
			balance_current = EXCLUDED.balance_current,
			balance_available = EXCLUDED.balance_available`

	_, err = tx.Exec(userQuery,
		fmt.Sprintf("%d", jsonData.Account.AccountID),
		jsonData.Account.AccountName,
		jsonData.Account.AccountType,
		jsonData.Account.AccountNumber,
		jsonData.Account.OwnerName,
		jsonData.Account.Balance.Current,
		jsonData.Account.Balance.Available,
		jsonData.Account.Balance.Currency,
		jsonData.Account.BankDetails.BankName,
		jsonData.Account.BankDetails.RoutingNumber,
		jsonData.Account.BankDetails.Branch,
	)
	if err != nil {
		return fmt.Errorf("failed to insert user data from %s: %w", filename, err)
	}

	// Get user prefix from filename (e.g., "jane" from "JaneDoe.json")
	userPrefix := strings.ToLower(strings.TrimSuffix(strings.TrimSuffix(filename, ".json"), "Doe"))

	// Batch insert transactions in chunks of 100
	const batchSize = 100
	totalTransactions := len(jsonData.Account.Transactions)
	successfulInserts := 0

	for i := 0; i < totalTransactions; i += batchSize {
		end := i + batchSize
		if end > totalTransactions {
			end = totalTransactions
		}

		batch := jsonData.Account.Transactions[i:end]
		valueStrings := make([]string, 0, len(batch))
		valueArgs := make([]interface{}, 0, len(batch)*7)

		for j, t := range batch {
			date, err := time.Parse("2006-01-02", t.Date)
			if err != nil {
				fmt.Printf("Warning: Skipping transaction with invalid date %s in %s: %v\n", t.Date, filename, err)
				continue
			}
			
			// Create unique transaction ID by prefixing with user identifier
			uniqueTransactionID := fmt.Sprintf("%s_%s", userPrefix, t.TransactionID)
			
			valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d)",
				j*7+1, j*7+2, j*7+3, j*7+4, j*7+5, j*7+6, j*7+7))
			valueArgs = append(valueArgs, 
				uniqueTransactionID,
				fmt.Sprintf("%d", t.AccountID),
				date,
				t.Amount,
				t.Category,
				t.Merchant,
				t.Location)
		}

		if len(valueStrings) > 0 {
			transactionQuery := fmt.Sprintf(`
				INSERT INTO transactions (
					transaction_id, account_id, date, amount, category, merchant, location
				) VALUES %s
				ON CONFLICT (transaction_id) DO NOTHING`, strings.Join(valueStrings, ","))
			
			result, err := tx.Exec(transactionQuery, valueArgs...)
			if err != nil {
				fmt.Printf("Warning: Failed to insert batch from %s: %v\n", filename, err)
				continue
			}

			inserted, _ := result.RowsAffected()
			successfulInserts += int(inserted)
			fmt.Printf("Batch processed for %s: %d/%d transactions inserted\n", filename, inserted, len(batch))
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction for %s: %w", filename, err)
	}

	fmt.Printf("Successfully processed %s: %d/%d transactions inserted\n", filename, successfulInserts, totalTransactions)
	return nil
}

// InsertAllUserData inserts data for all users
func InsertAllUserData(db *sql.DB) error {
	users := []string{
		"JaneDoe.json",
		"JohnDoe.json",
		"JillDoe.json",
		"JakeDoe.json",
	}

	for _, user := range users {
		if err := InsertUserData(db, user); err != nil {
			return fmt.Errorf("failed to insert data for %s: %w", user, err)
		}
	}

	return nil
}

// CreateUser creates a new user in the database
func CreateUser(db *sql.DB, user *types.User) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO users (
			account_id, account_name, account_type, account_number, 
			owner_name, balance_current, balance_available, balance_currency,
			bank_name, routing_number, branch
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err = tx.Exec(query,
		user.AccountID,
		user.AccountName,
		user.AccountType,
		user.AccountNumber,
		user.OwnerName,
		user.Balance.Current,
		user.Balance.Available,
		user.Balance.Currency,
		user.BankDetails.BankName,
		user.BankDetails.RoutingNumber,
		user.BankDetails.Branch,
	)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	return tx.Commit()
}

// GetTransactions retrieves all transactions for a given account
func GetTransactions(db *sql.DB, accountID string) ([]types.Transaction, error) { 
	query := ` 
		SELECT transaction_id, account_id, date, amount, category, merchant, location 
		FROM transactions 
		WHERE account_id = $1
		ORDER BY date DESC`
	
	rows, err := db.Query(query, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to query transactions: %w", err)
	}
	defer rows.Close()

	var transactions []types.Transaction
	for rows.Next() {
		var t types.Transaction
		var prefixedTransactionID string
		if err := rows.Scan(
			&prefixedTransactionID,
			&t.AccountID,
			&t.Date,
			&t.Amount,
			&t.Category,
			&t.Merchant,
			&t.Location,
		); err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		// Extract the original transaction ID by removing the prefix (e.g., "jane_123" -> "123")
		parts := strings.SplitN(prefixedTransactionID, "_", 2)
		if len(parts) == 2 {
			t.TransactionID = parts[1]  // Use the part after the prefix
			t.UserPrefix = parts[0]     // Store the prefix (user identifier) if needed
		} else {
			t.TransactionID = prefixedTransactionID  // Fallback to full ID if no prefix found
		}
		transactions = append(transactions, t)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating transactions: %w", err)
	}

	return transactions, nil
}

// InsertTransaction inserts a new transaction into the database
func InsertTransaction(db *sql.DB, transaction *types.Transaction) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO transactions (
			transaction_id, account_id, date, amount, category, merchant, location
		) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	
	_, err = tx.Exec(query,
		transaction.TransactionID,
		transaction.AccountID,
		transaction.Date,
		transaction.Amount,
		transaction.Category,
		transaction.Merchant,
		transaction.Location,
	)
	if err != nil {
		return fmt.Errorf("failed to insert transaction: %w", err)
	}

	return tx.Commit()
} 
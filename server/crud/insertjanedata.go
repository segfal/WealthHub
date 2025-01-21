package crud

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

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
				AccountID     int     `json:"account_id"`
				Date          string  `json:"date"`
				Amount        float64 `json:"amount"`
				Category      string  `json:"category"`
				Merchant      string  `json:"merchant"`
				Location      string  `json:"location"`
				Type          string  `json:"type"`
				Status        string  `json:"status"`
				Timestamp     string  `json:"timestamp"`
			} `json:"transactions"`
		} `json:"account"`
	}

	if err := json.Unmarshal(data, &jsonData); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	fmt.Printf("Found %d transactions in JSON file\n", len(jsonData.Account.Transactions))

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert user with bank details
	userQuery := `
		INSERT INTO users (
			account_id, account_name, account_type, account_number, 
			owner_name, balance_current, balance_available, balance_currency,
			bank_name, routing_number, branch
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err = tx.Exec(userQuery,
		fmt.Sprintf("%d", jsonData.Account.AccountID), // Convert int to string
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
		return fmt.Errorf("failed to insert user: %w", err)
	}

	// Batch insert transactions
	if len(jsonData.Account.Transactions) > 0 {
		// Create batch insert query
		valueStrings := make([]string, 0, len(jsonData.Account.Transactions))
		valueArgs := make([]interface{}, 0, len(jsonData.Account.Transactions)*7)
		for i, t := range jsonData.Account.Transactions {
			date, err := time.Parse("2006-01-02", t.Date)
			if err != nil {
				return fmt.Errorf("failed to parse date %s: %w", t.Date, err)
			}

			valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d)",
				i*7+1, i*7+2, i*7+3, i*7+4, i*7+5, i*7+6, i*7+7))
			valueArgs = append(valueArgs,
				t.TransactionID,
				fmt.Sprintf("%d", t.AccountID), // Convert int to string
				date,
				t.Amount,
				t.Category,
				t.Merchant,
				t.Location)
		}

		fmt.Printf("Inserting %d transactions\n", len(valueStrings))

		transactionQuery := fmt.Sprintf(`
			INSERT INTO transactions (
				transaction_id, account_id, date, amount, category, merchant, location
			) VALUES %s`, strings.Join(valueStrings, ","))

		_, err = tx.Exec(transactionQuery, valueArgs...)
		if err != nil {
			return fmt.Errorf("failed to insert transactions: %w", err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

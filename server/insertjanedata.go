package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"server/types"
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
			AccountID     string `json:"account_id"`
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
				Date          string  `json:"date"`
				Amount        float64 `json:"amount"`
				Category      string  `json:"category"`
				Merchant      string  `json:"merchant"`
				Location      string  `json:"location"`
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
		AccountID:     jsonData.Account.AccountID,
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
			AccountID:     jsonData.Account.AccountID,
			Date:          date,
			Amount:        t.Amount,
			Category:      t.Category,
			Merchant:      t.Merchant,
			Location:      t.Location,
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

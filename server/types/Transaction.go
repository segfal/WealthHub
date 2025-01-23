package types

import "time"

// Transaction represents a financial transaction as per init.sql schema
type Transaction struct {
	TransactionID string    `json:"transaction_id"` // VARCHAR(20) PRIMARY KEY
	AccountID     string    `json:"account_id"`     // VARCHAR(20) REFERENCES users(account_id)
	Date          time.Time `json:"date"`           // TIMESTAMP
	Amount        float64   `json:"amount"`         // DECIMAL(10, 2)
	Category      string    `json:"category"`       // VARCHAR(50)
	Merchant      string    `json:"merchant"`       // VARCHAR(50)
	Location      string    `json:"location"`       // VARCHAR(100)
	UserPrefix    string    `json:"userPrefix,omitempty"`
}
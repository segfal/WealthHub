package types

import "time"

type Transaction struct {
	TransactionID string    `json:"transactionId"`
	AccountID     string    `json:"accountId"`
	Date          time.Time `json:"date"`
	Amount        float64   `json:"amount"`
	Category      string    `json:"category"`
	Merchant      string    `json:"merchant"`
	Location      string    `json:"location"`
}
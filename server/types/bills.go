package types

import "time"

// RecurringBill represents a bill that occurs regularly
type RecurringBill struct {
	Merchant        string    `json:"merchant"`
	Category        string    `json:"category"`
	MonthsPresent   int       `json:"months_present"`
	AverageAmount   float64   `json:"average_amount"`
	MedianAmount    float64   `json:"median_amount"`
	FirstOccurrence time.Time `json:"first_occurrence"`
	LastOccurrence  time.Time `json:"last_occurrence"`
}

// UpcomingBill represents a predicted future bill payment
type UpcomingBill struct {
	Merchant        string    `json:"merchant"`
	Category        string    `json:"category"`
	ExpectedAmount  float64   `json:"expected_amount"`
	DueDate         time.Time `json:"due_date"`
} 
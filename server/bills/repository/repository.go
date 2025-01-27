package repository

import (
	"context"
	"server/types"
	"time"
)

// Repository defines the interface for bill-related data operations
type Repository interface {
	// GetBillTotals retrieves total bill payments by category for a given time period
	GetBillTotals(ctx context.Context, accountID string, startDate, endDate time.Time) (map[string]float64, error)

	// GetRecurringBills retrieves recurring bill payments for an account
	GetRecurringBills(ctx context.Context, accountID string) ([]types.RecurringBill, error)

	// GetUpcomingBills retrieves upcoming bill payments based on recurring patterns
	GetUpcomingBills(ctx context.Context, accountID string) ([]types.UpcomingBill, error)

	// GetBillHistory retrieves historical bill payments for a specific merchant
	GetBillHistory(ctx context.Context, accountID string, merchantName string) ([]types.Transaction, error)

	// GetBillsByMonth retrieves all bill payments for a specific month
	GetBillsByMonth(ctx context.Context, accountID string, year int, month int) ([]types.Transaction, error)
} 
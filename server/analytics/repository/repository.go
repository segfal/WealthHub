package repository

import (
	"context"
	"server/types"
)

// Repository defines the interface for analytics data operations
type Repository interface {
	// GetTransactions retrieves transactions for analysis
	GetTransactions(ctx context.Context, accountID string, timeRange string) ([]types.Transaction, error)

	// GetCategoryTotals retrieves total spending by category
	GetCategoryTotals(ctx context.Context, accountID string, timeRange string) (map[string]float64, error)

	// GetAccount retrieves account information
	GetAccount(ctx context.Context, accountID string) (*types.Account, error)

	// GetMonthlyIncome retrieves income transactions for a specific month
	GetMonthlyIncome(ctx context.Context, accountID string, year int, month int) ([]types.Transaction, error)

	// GetBillPayments retrieves bill payment transactions for a specific month
	GetBillPayments(ctx context.Context, accountID string, year int, month int) ([]types.Transaction, error)

	// GetDailySpending retrieves daily spending transactions for a specific month
	GetDailySpending(ctx context.Context, accountID string, year int, month int) ([]types.Transaction, error)

	// GetMonthlySpending retrieves monthly spending transactions for a specific month
	GetMonthlySpending(ctx context.Context, accountID string, year int, month int) ([]types.Transaction, error)

	// GetCategoryDiversity retrieves category diversity for a specific month
	GetCategoryDiversity(ctx context.Context, accountID string, year int, month int) (map[string]int, error)
}

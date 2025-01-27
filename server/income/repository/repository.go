package repository

import (
	"context"
	"server/types"
)

// Repository defines the interface for income-related data operations
type Repository interface {
	// GetIncome retrieves all income transactions for an account
	GetIncome(ctx context.Context, accountID string) ([]types.Transaction, error)

	// GetMonthlyIncome retrieves income transactions for a specific month
	GetMonthlyIncome(ctx context.Context, accountID string, year int, month int) ([]types.Transaction, error)
} 
package repository

import (
	"context"
	"server/types"
)

// Repository defines the interface for category-related data operations
type Repository interface {
	// GetCategories retrieves all categories for an account
	GetCategories(ctx context.Context, accountID string) ([]types.Category, error)

	// GetCategoryTotals retrieves total spending by category
	GetCategoryTotals(ctx context.Context, accountID string) (map[string]float64, error)
} 
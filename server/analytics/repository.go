package analytics

import (
	"context"
	"server/types"
)

type Repository interface {
	GetTransactions(ctx context.Context, accountID string, timeRange string) ([]types.Transaction, error)
	GetCategoryTotals(ctx context.Context, accountID string, timeRange string) (map[string]float64, error)
} 
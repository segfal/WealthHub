package analytics

import (
	"context"
	"server/types"
	"time"
)

type UserInfo struct {
	UserID string
	AccountID string
}

// Service defines the interface for analytics operations
type Service interface {
	// AnalyzeSpending analyzes spending patterns for a given account and time range
	AnalyzeSpending(ctx context.Context, accountID string, timeRange string) (*types.SpendingAnalytics, error)
	
	// GetTimePatterns analyzes spending patterns by time of day and day of week
	GetTimePatterns(ctx context.Context, accountID string, startDate, endDate time.Time) ([]types.TimePattern, error)
	
	// PredictSpending generates spending predictions for each category
	PredictSpending(ctx context.Context, accountID string) ([]types.PredictedSpend, error)
}

// Repository defines the interface for analytics data operations
type Repository interface {
	// GetTransactions retrieves transactions for analysis
	GetTransactions(ctx context.Context, accountID string, timeRange string) ([]types.Transaction, error)
	
	// GetCategoryTotals retrieves total spending by category
	GetCategoryTotals(ctx context.Context, accountID string, timeRange string) (map[string]float64, error)
	
	// GetAccount retrieves account information
	GetAccount(ctx context.Context, accountID string) (*types.Account, error) 

	GetBillTotals(ctx context.Context, accountID string, timeRange string) (map[string]float64, error);
} 
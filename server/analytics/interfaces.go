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

	// GetMonthlyIncome retrieves income transactions for a specific month
	GetMonthlyIncome(ctx context.Context, accountID string, year int, month int) ([]types.Transaction, error)

	// GetBillPayments retrieves bill payment transactions for a specific month
	GetBillPayments(ctx context.Context, accountID string, year int, month int) ([]types.Transaction, error)

	// GetDailyPatterns retrieves daily spending patterns for a specific month
	GetDailyPatterns(ctx context.Context, accountID string, year int, month int) ([]types.DailyPattern, error)

	// GetMonthlyPatterns retrieves monthly spending patterns for a specific month
	GetMonthlyPatterns(ctx context.Context, accountID string, year int, month int) ([]types.MonthlyPattern, error)
	
}

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

	// GetRecentSpending retrieves and analyzes spending data for the recent period
	GetRecentSpending(ctx context.Context, accountID string, startDate, endDate time.Time) ([]types.Transaction, error)

	// GetDailySpending retrieves daily spending transactions for a specific month
	GetDailySpending(ctx context.Context, accountID string, year int, month int) ([]types.Transaction, error)

	// GetMonthlySpending retrieves monthly spending transactions for a specific month
	GetMonthlySpending(ctx context.Context, accountID string, year int, month int) ([]types.Transaction, error)

	// GetCategoryDiversity retrieves category diversity for a specific month
	GetCategoryDiversity(ctx context.Context, accountID string, year int, month int) (map[string]int, error)

	
	
} 

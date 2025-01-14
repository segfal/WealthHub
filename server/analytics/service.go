package analytics

import (
	"context"
	"fmt"
	"server/types"
	"sort"
	"strconv"
	"time"
)

type service struct {
	repo Repository
}

// NewService creates a new analytics service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) AnalyzeSpending(ctx context.Context, accountID string, timeRange string) (*types.SpendingAnalytics, error) {
	transactions, err := s.repo.GetTransactions(ctx, accountID, timeRange)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	categoryTotals, err := s.repo.GetCategoryTotals(ctx, accountID, timeRange)
	if err != nil {
		return nil, fmt.Errorf("failed to get category totals: %w", err)
	}

	// Calculate total spent and categories
	var totalSpent float64
	var topCategories []types.CategorySpend
	for category, amount := range categoryTotals {
		totalSpent += amount
		topCategories = append(topCategories, types.CategorySpend{
			Category:   category,
			TotalSpent: fmt.Sprintf("%.2f", amount),
			Percentage: fmt.Sprintf("%.2f", (amount/totalSpent)*100),
		})
	}

	// Sort by actual amounts
	sort.Slice(topCategories, func(i, j int) bool {
		amtI, _ := strconv.ParseFloat(topCategories[i].TotalSpent, 64)
		amtJ, _ := strconv.ParseFloat(topCategories[j].TotalSpent, 64)
		return amtI > amtJ
	})

	if len(topCategories) > 5 {
		topCategories = topCategories[:5]
	}

	patterns := analyzeTimePatterns(transactions)
	predictions := s.predictSpendingInternal(transactions, categoryTotals)

	return &types.SpendingAnalytics{
		TopCategories:     topCategories,
		SpendingPatterns:  patterns,
		PredictedSpending: predictions,
		TotalSpent:        totalSpent,
		MonthlyAverage:    totalSpent / float64(timeRangeToMonths(timeRange)),
	}, nil
}

func (s *service) GetTimePatterns(ctx context.Context, accountID string, startDate, endDate time.Time) ([]types.TimePattern, error) {
	timeRange := fmt.Sprintf("'%s'::timestamp - '%s'::timestamp", endDate.Format(time.RFC3339), startDate.Format(time.RFC3339))
	transactions, err := s.repo.GetTransactions(ctx, accountID, timeRange)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}
	return analyzeTimePatterns(transactions), nil
}

func (s *service) PredictSpending(ctx context.Context, accountID string) ([]types.PredictedSpend, error) {
	transactions, err := s.repo.GetTransactions(ctx, accountID, "6 months")
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	categoryTotals, err := s.repo.GetCategoryTotals(ctx, accountID, "6 months")
	if err != nil {
		return nil, fmt.Errorf("failed to get category totals: %w", err)
	}

	return s.predictSpendingInternal(transactions, categoryTotals), nil
}

func (s *service) predictSpendingInternal(transactions []types.Transaction, categoryTotals map[string]float64) []types.PredictedSpend {
	categoryTransactions := make(map[string][]types.Transaction)
	for _, t := range transactions {
		categoryTransactions[t.Category] = append(categoryTransactions[t.Category], t)
	}

	var predictions []types.PredictedSpend
	for category, txns := range categoryTransactions {
		if len(txns) >= 3 {
			avgTimeBetween := calculateAverageTimeBetween(txns)
			lastTransaction := txns[len(txns)-1]
			predictedDate := lastTransaction.Date.Add(avgTimeBetween)
			
			frequency := float64(len(txns))
			total := categoryTotals[category]
			avgAmount := total / frequency
			
			normalizedFreq := minFloat64(frequency/30.0, 1.0)
			normalizedAmount := minFloat64(avgAmount/1000.0, 1.0)
			likelihood := (normalizedFreq + normalizedAmount) / 2.0
			
			warning := ""
			if likelihood > 0.7 {
				warning = "High likelihood of significant spending in this category soon"
			}

			predictions = append(predictions, types.PredictedSpend{
				Category:      category,
				Likelihood:    likelihood,
				PredictedDate: predictedDate,
				Warning:       warning,
			})
		}
	}
	return predictions
}

func timeRangeToMonths(timeRange string) float64 {
	switch timeRange {
	case "1 month":
		return 1
	case "3 months":
		return 3
	case "6 months":
		return 6
	case "1 year":
		return 12
	default:
		return 1
	}
}

func minFloat64(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
} 
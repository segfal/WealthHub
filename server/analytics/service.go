package analytics

import (
	"context"
	"fmt"
	"math"
	"server/types"
	"sort"
	"strconv"
	"time"
)

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// AnalyzeSpending implements Service.AnalyzeSpending
func (s *service) AnalyzeSpending(ctx context.Context, accountID string, timeRange string) (*types.SpendingAnalytics, error) {
	// First, verify the account exists
	account, err := s.repo.GetAccount(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	categoryTotals, err := s.repo.GetCategoryTotals(ctx, accountID, timeRange)
	if err != nil {
		return nil, fmt.Errorf("failed to get category totals: %w", err)
	}

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

	// Get current month's bill payments
	now := time.Now()
	billPayments, err := s.repo.GetBillPayments(ctx, accountID, now.Year(), int(now.Month()))
	if err != nil {
		return nil, fmt.Errorf("failed to get bill payments: %w", err)
	}

	var totalSpentBills float64
	billTotalsByMerchant := make(map[string]float64)
	
	// Calculate totals by merchant for bill payments
	for _, payment := range billPayments {
		amount := math.Abs(payment.Amount)
		totalSpentBills += amount
		billTotalsByMerchant[payment.Merchant] += amount
	}

	var topBills []types.BillPayment
	for merchant, amount := range billTotalsByMerchant {
		topBills = append(topBills, types.BillPayment{
			Category:   merchant,
			TotalSpent: fmt.Sprintf("%.2f", amount),
			Percentage: fmt.Sprintf("%.2f", (amount/totalSpentBills)*100),
		})
	}

	// Sort by amount spent
	sort.Slice(topBills, func(i, j int) bool {
		amtI, _ := strconv.ParseFloat(topBills[i].TotalSpent, 64)
		amtJ, _ := strconv.ParseFloat(topBills[j].TotalSpent, 64)
		return amtI > amtJ
	})

	// Get top 5 categories
	if len(topCategories) > 5 {
		topCategories = topCategories[:5]
	}

	// Get time patterns for the last month
	endDate := time.Now()
	startDate := endDate.AddDate(0, -1, 0)
	patterns, err := s.GetTimePatterns(ctx, accountID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze time patterns: %w", err)
	}

	// Get spending predictions
	predictions, err := s.PredictSpending(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to predict spending: %w", err)
	}

	return &types.SpendingAnalytics{
		Account:          account,
		TopCategories:    topCategories,
		SpendingPatterns: patterns,
		PredictedSpending: predictions,
		TotalSpent:       totalSpent,
		MonthlyAverage:   totalSpent / float64(timeRangeToMonths(timeRange)),
	}, nil
}

// GetTimePatterns implements Service.GetTimePatterns
func (s *service) GetTimePatterns(ctx context.Context, accountID string, startDate, endDate time.Time) ([]types.TimePattern, error) {
	// First, verify the account exists
	if _, err := s.repo.GetAccount(ctx, accountID); err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	// Convert the date range to a PostgreSQL interval string
	timeRange := "1 month"
	
	transactions, err := s.repo.GetTransactions(ctx, accountID, timeRange)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	// Group transactions by day and hour
	patterns := make(map[string]map[string]struct {
		totalAmount float64
		count      int
	})

	for _, t := range transactions {
		dayOfWeek := t.Date.Format("Monday")
		hourOfDay := t.Date.Format("15:00")

		if _, exists := patterns[dayOfWeek]; !exists {
			patterns[dayOfWeek] = make(map[string]struct {
				totalAmount float64
				count      int
			})
		}

		stats := patterns[dayOfWeek][hourOfDay]
		stats.totalAmount += math.Abs(t.Amount)
		stats.count++
		patterns[dayOfWeek][hourOfDay] = stats
	}

	// Convert to TimePattern slice
	var result []types.TimePattern
	for day, hours := range patterns {
		for hour, stats := range hours {
			result = append(result, types.TimePattern{
				TimeOfDay:    hour,
				DayOfWeek:    day,
				Frequency:    stats.count,
				AverageSpend: stats.totalAmount / float64(stats.count),
			})
		}
	}

	// Sort by frequency and average spend
	sort.Slice(result, func(i, j int) bool {
		if result[i].Frequency == result[j].Frequency {
			return result[i].AverageSpend > result[j].AverageSpend
		}
		return result[i].Frequency > result[j].Frequency
	})

	return result, nil
}

// PredictSpending implements Service.PredictSpending
func (s *service) PredictSpending(ctx context.Context, accountID string) ([]types.PredictedSpend, error) {
	// First, verify the account exists
	if _, err := s.repo.GetAccount(ctx, accountID); err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	// Get last 6 months of transactions for better prediction
	transactions, err := s.repo.GetTransactions(ctx, accountID, "6 months")
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	// Group transactions by category
	categoryTransactions := make(map[string][]types.Transaction)
	for _, t := range transactions {
		categoryTransactions[t.Category] = append(categoryTransactions[t.Category], t)
	}

	var predictions []types.PredictedSpend
	for category, txns := range categoryTransactions {
		if len(txns) < 3 {
			continue // Need at least 3 transactions for prediction
		}

		// Sort transactions by date
		sort.Slice(txns, func(i, j int) bool {
			return txns[i].Date.Before(txns[j].Date)
		})

		// Calculate average time between transactions
		var totalDuration time.Duration
		for i := 1; i < len(txns); i++ {
			totalDuration += txns[i].Date.Sub(txns[i-1].Date)
		}
		avgTimeBetween := totalDuration / time.Duration(len(txns)-1)

		// Calculate frequency and amount metrics
		frequency := float64(len(txns)) / 180 // Normalize by 6 months (180 days)
		var totalAmount float64
		for _, t := range txns {
			totalAmount += math.Abs(t.Amount)
		}
		avgAmount := totalAmount / float64(len(txns))

		// Calculate likelihood score
		normalizedFreq := math.Min(frequency*30, 1.0)  // Normalize to max 1.0 (30 days)
		normalizedAmount := math.Min(avgAmount/1000, 1.0) // Normalize to max 1.0 ($1000)
		likelihood := (normalizedFreq + normalizedAmount) / 2.0

		// Generate prediction
		lastTransaction := txns[len(txns)-1]
		predictedDate := lastTransaction.Date.Add(avgTimeBetween)

		warning := ""
		if likelihood > 0.7 {
			warning = fmt.Sprintf("High likelihood (%.0f%%) of spending in %s category around %s",
				likelihood*100, category, predictedDate.Format("Jan 02"))
		}

		predictions = append(predictions, types.PredictedSpend{
			Category:      category,
			Likelihood:    likelihood,
			PredictedDate: predictedDate,
			Warning:       warning,
		})
	}

	// Sort by likelihood
	sort.Slice(predictions, func(i, j int) bool {
		return predictions[i].Likelihood > predictions[j].Likelihood
	})

	return predictions, nil
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

// GetMonthlyIncome implements Service.GetMonthlyIncome
func (s *service) GetMonthlyIncome(ctx context.Context, accountID string, year int, month int) ([]types.Transaction, error) {
	// First, verify the account exists
	if _, err := s.repo.GetAccount(ctx, accountID); err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	return s.repo.GetMonthlyIncome(ctx, accountID, year, month)
}

// GetBillPayments implements Service.GetBillPayments
func (s *service) GetBillPayments(ctx context.Context, accountID string, year int, month int) ([]types.Transaction, error) {
	// First, verify the account exists
	if _, err := s.repo.GetAccount(ctx, accountID); err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	return s.repo.GetBillPayments(ctx, accountID, year, month)
}

// GetDailyPatterns implements Service.GetDailyPatterns
func (s *service) GetDailyPatterns(ctx context.Context, accountID string, year int, month int) ([]types.DailyPattern, error) {
	// First, verify the account exists
	if _, err := s.repo.GetAccount(ctx, accountID); err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	transactions, err := s.repo.GetDailySpending(ctx, accountID, year, month)
	if err != nil {
		return nil, fmt.Errorf("failed to get daily spending: %w", err)
	}

	// Group transactions by day
	dailyTotals := make(map[string]float64)
	for _, tx := range transactions {
		day := tx.Date.Format("Monday")
		dailyTotals[day] += math.Abs(tx.Amount)
	}

	// Convert to DailyPattern slice
	var patterns []types.DailyPattern
	for day, total := range dailyTotals {
		patterns = append(patterns, types.DailyPattern{
			DayOfWeek: day,
			AverageAmount: total,
		})
	}

	return patterns, nil
}

// GetMonthlyPatterns implements Service.GetMonthlyPatterns
func (s *service) GetMonthlyPatterns(ctx context.Context, accountID string, year int, month int) ([]types.MonthlyPattern, error) {
	// First, verify the account exists
	if _, err := s.repo.GetAccount(ctx, accountID); err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	transactions, err := s.repo.GetMonthlySpending(ctx, accountID, year, month)
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly spending: %w", err)
	}

	// Group transactions by month
	monthlyTotals := make(map[string]float64)
	for _, tx := range transactions {
		monthName := tx.Date.Format("January")
		monthlyTotals[monthName] += math.Abs(tx.Amount)
	}

	// Convert to MonthlyPattern slice
	var patterns []types.MonthlyPattern
	for month, total := range monthlyTotals {
		patterns = append(patterns, types.MonthlyPattern{
			Month: month,
			AverageAmount: total,
		})
	}

	return patterns, nil
} 
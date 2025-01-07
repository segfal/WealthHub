package main

import (
	"database/sql"
	"sort"
	"time" 
	"fmt"
)

type Transaction struct {
	TransactionID string    `json:"transactionId"`
	AccountID     string    `json:"accountId"`
	Date          time.Time `json:"date"`
	Amount        float64   `json:"amount"`
	Category      string    `json:"category"`
	Merchant      string    `json:"merchant"`
	Location      string    `json:"location"`
}

type SpendingAnalytics struct {
	TopCategories     []CategorySpend    `json:"topCategories"`
	SpendingPatterns  []TimePattern      `json:"spendingPatterns"`
	PredictedSpending []PredictedSpend   `json:"predictedSpending"`
	TotalSpent        float64            `json:"totalSpent"`
	MonthlyAverage    float64            `json:"monthlyAverage"`
}

type CategorySpend struct {
	Category    string  `json:"category"`
	TotalSpent  string `json:"totalSpent"`
	Percentage  string `json:"percentage"`
}

type TimePattern struct {
	TimeOfDay   string  `json:"timeOfDay"`
	DayOfWeek   string  `json:"dayOfWeek"`
	Frequency   int     `json:"frequency"`
	AverageSpend float64 `json:"averageSpend"`
}

type PredictedSpend struct {
	Category    string    `json:"category"`
	Likelihood  float64   `json:"likelihood"`
	PredictedDate time.Time `json:"predictedDate"`
	Warning     string    `json:"warning"`
}

func analyzeSpending(db *sql.DB, accountID string, timeRange string) (*SpendingAnalytics, error) {
	// Get transactions for the specified time range
	query := `
		SELECT amount, category, date 
		FROM transactions 
		WHERE account_id = $1 
		AND date >= NOW() - $2::INTERVAL`
	
	rows, err := db.Query(query, accountID, timeRange)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categoryTotals := make(map[string]float64)
	var totalSpent float64
	var transactions []Transaction
	for rows.Next() {
		var t Transaction
		if err := rows.Scan(&t.Amount, &t.Category, &t.Date); err != nil {
			return nil, err
		}

		categoryTotals[t.Category] += t.Amount
		totalSpent += t.Amount
		transactions = append(transactions, t)
	}

	// Calculate top categories
	var topCategories []CategorySpend
	for category, amount := range categoryTotals {
		var percentage float64 = (amount / totalSpent) * 100;
		topCategories = append(topCategories, CategorySpend{
			Category:   category,
			TotalSpent: fmt.Sprintf("%.2f", amount),
			Percentage: fmt.Sprintf("%.2f", percentage),
		})
	}

	// Sort categories by spend
	sort.Slice(topCategories, func(i, j int) bool {
		return topCategories[i].TotalSpent > topCategories[j].TotalSpent
	})

	// Analyze spending patterns
	patterns := analyzeTimePatterns(transactions)

	// Generate predictions
	predictions := predictFutureSpending(transactions, categoryTotals)

	return &SpendingAnalytics{
		TopCategories:     topCategories[:minInt(len(topCategories), 5)],
		SpendingPatterns:  patterns,
		PredictedSpending: predictions,
		TotalSpent:        totalSpent,
		MonthlyAverage:    totalSpent / float64(timeRangeToMonths(timeRange)),
	}, nil
}

func analyzeTimePatterns(transactions []Transaction) []TimePattern {
	timePatterns := make(map[string]map[string][]float64)
	
	for _, t := range transactions {
		hour := t.Date.Format("15:00")
		day := t.Date.Format("Monday")
		
		if _, exists := timePatterns[day]; !exists {
			timePatterns[day] = make(map[string][]float64)
		}
		timePatterns[day][hour] = append(timePatterns[day][hour], t.Amount)
	}

	var patterns []TimePattern
	for day, hours := range timePatterns {
		for hour, amounts := range hours {
			var sum float64
			for _, amount := range amounts {
				sum += amount
			}
			
			patterns = append(patterns, TimePattern{
				TimeOfDay:    hour,
				DayOfWeek:    day,
				Frequency:    len(amounts),
				AverageSpend: sum / float64(len(amounts)),
			})
		}
	}

	return patterns
}

/** @dev
 * @param transactions []Transaction
 * @param categoryTotals map[string]float64
 * @return []PredictedSpend
 */
func predictFutureSpending(transactions []Transaction, categoryTotals map[string]float64) []PredictedSpend {
	var predictions []PredictedSpend

	// Simple prediction based on frequency and amount
	for category, total := range categoryTotals {
		var categoryTransactions []Transaction
		for _, t := range transactions {
			if t.Category == category {
				categoryTransactions = append(categoryTransactions, t)
			}
		}

		if len(categoryTransactions) >= 3 {
			avgTimeBetween := calculateAverageTimeBetween(categoryTransactions)
			lastTransaction := categoryTransactions[len(categoryTransactions)-1]
			predictedDate := lastTransaction.Date.Add(avgTimeBetween)
			
			likelihood := calculateLikelihood(categoryTransactions, total)
			
			warning := ""
			if likelihood > 0.7 {
				warning = "High likelihood of significant spending in this category soon"
			}

			predictions = append(predictions, PredictedSpend{
				Category:      category,
				Likelihood:    likelihood,
				PredictedDate: predictedDate,
				Warning:       warning,
			})
		}
	}

	return predictions
}

func calculateAverageTimeBetween(transactions []Transaction) time.Duration {
	if len(transactions) < 2 {
		return time.Hour * 24 * 7 // Default to weekly if not enough data
	}

	var totalDuration time.Duration
	for i := 1; i < len(transactions); i++ {
		duration := transactions[i].Date.Sub(transactions[i-1].Date)
		totalDuration += duration
	}

	return totalDuration / time.Duration(len(transactions)-1)
}

func calculateLikelihood(transactions []Transaction, totalSpent float64) float64 {
	// Simple likelihood calculation based on frequency and amount
	frequency := float64(len(transactions))
	averageAmount := totalSpent / frequency
	
	// Higher frequency and amount = higher likelihood
	normalizedFreq := minFloat64(frequency/30.0, 1.0) // Normalize to max 1.0
	normalizedAmount := minFloat64(averageAmount/1000.0, 1.0) // Normalize to max 1.0
	
	return (normalizedFreq + normalizedAmount) / 2.0
}

func timeRangeToMonths(timeRange string) float64 {
	// Convert time range string to approximate number of months
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

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
} 
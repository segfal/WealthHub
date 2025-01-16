package analytics

import (
	"server/types"
	"sort"
	"time"
)

// analyzeTimePatterns analyzes spending patterns by time of day and day of week
func analyzeTimePatterns(transactions []types.Transaction) []types.TimePattern {
	if len(transactions) == 0 {
		return nil
	}

	// Use a more efficient data structure with a composite key
	type timeKey struct {
		hour string
		day  string
	}
	patterns := make(map[timeKey]struct {
		sum   float64
		count int
	})

	for _, t := range transactions {
		key := timeKey{
			hour: t.Date.Format("15:00"),
			day:  t.Date.Format("Monday"),
		}
		val := patterns[key]
		val.sum += t.Amount
		val.count++
		patterns[key] = val
	}

	// Convert to slice for return
	result := make([]types.TimePattern, 0, len(patterns))
	for key, val := range patterns {
		result = append(result, types.TimePattern{
			TimeOfDay:    key.hour,
			DayOfWeek:    key.day,
			Frequency:    val.count,
			AverageSpend: val.sum / float64(val.count),
		})
	}

	return result
}

// calculateAverageTimeBetween calculates the average time between transactions
func calculateAverageTimeBetween(transactions []types.Transaction) time.Duration {
	if len(transactions) < 2 {
		return time.Hour * 24 * 7 // Default to weekly if not enough data
	}

	// Pre-sort transactions by date to ensure correct calculation
	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].Date.Before(transactions[j].Date)
	})

	totalDuration := transactions[len(transactions)-1].Date.Sub(transactions[0].Date)
	return totalDuration / time.Duration(len(transactions)-1)
} 
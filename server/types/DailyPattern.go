// DailyPattern represents a daily spending pattern it will get the transactions from the database
package types

type DailyPattern struct {
	DayOfWeek string
	AverageAmount float64
}


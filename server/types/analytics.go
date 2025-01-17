package types

import "time"

type SpendingAnalytics struct {
	TopCategories     []CategorySpend   `json:"topCategories"`
	SpendingPatterns  []TimePattern     `json:"spendingPatterns"`
	PredictedSpending []PredictedSpend  `json:"predictedSpending"`
	TotalSpent        float64           `json:"totalSpent"`
	MonthlyAverage    float64           `json:"monthlyAverage"`
}

type CategorySpend struct {
	Category   string `json:"category"`
	TotalSpent string `json:"totalSpent"`
	Percentage string `json:"percentage"`
}

type TimePattern struct {
	TimeOfDay    string  `json:"timeOfDay"`
	DayOfWeek    string  `json:"dayOfWeek"`
	Frequency    int     `json:"frequency"`
	AverageSpend float64 `json:"averageSpend"`
}

type PredictedSpend struct {
	Category      string    `json:"category"`
	Likelihood    float64   `json:"likelihood"`
	PredictedDate time.Time `json:"predictedDate"`
	Warning       string    `json:"warning,omitempty"`
} 
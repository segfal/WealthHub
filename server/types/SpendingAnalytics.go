package types

// SpendingAnalytics represents the spending analysis for an account
type SpendingAnalytics struct {
	Account          *Account         `json:"account"`
	TopCategories    []CategorySpend  `json:"top_categories"`
	SpendingPatterns []TimePattern    `json:"spending_patterns"`
	PredictedSpending []PredictedSpend `json:"predicted_spending"`
	TotalSpent       float64          `json:"total_spent"`
	MonthlyAverage   float64          `json:"monthly_average"`
}


package types

type SpendingAnalytics struct {
	TopCategories     []CategorySpend    `json:"topCategories"`
	SpendingPatterns  []TimePattern      `json:"spendingPatterns"`
	PredictedSpending []PredictedSpend   `json:"predictedSpending"`
	TotalSpent        float64            `json:"totalSpent"`
	MonthlyAverage    float64            `json:"monthlyAverage"`
}


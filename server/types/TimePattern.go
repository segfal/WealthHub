package types




type TimePattern struct {
	TimeOfDay   string  `json:"timeOfDay"`
	DayOfWeek   string  `json:"dayOfWeek"`
	Frequency   int     `json:"frequency"`
	AverageSpend float64 `json:"averageSpend"`
}
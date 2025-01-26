package types

// Category represents a transaction category
type Category struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	TotalSpent  float64 `json:"total_spent"`
	Count       int     `json:"count"`
} 
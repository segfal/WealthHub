package types

type CategorySpend struct {
	Category    string  `json:"category"`
	TotalSpent  string `json:"totalSpent"`
	Percentage  string `json:"percentage"`
}
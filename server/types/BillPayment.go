package types

type BillPayment struct { 
	Category    string  `json:"category"`
	TotalSpent  string `json:"totalSpent"`
	Percentage  string `json:"percentage"` 
	Time 		string `json:"timestamp"` 
}
package types

// BillDetail represents a bill payment with its details
type BillDetail struct {
	Name              string  `json:"name"`
	Category          string  `json:"category"`
	Amount            float64 `json:"amount"`
	PaycheckPercentage float64 `json:"paycheckPercentage"`
	LastPaidDate      string  `json:"lastPaidDate"`
	Merchant          string  `json:"merchant"`
	Location          string  `json:"location"`
} 
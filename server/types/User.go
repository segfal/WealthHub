package types

type UserBalance struct {
	Current   float64 `json:"current"`
	Available float64 `json:"available"`
	Currency  string  `json:"currency"`
}

type UserBankDetails struct {
	BankName      string `json:"bank_name"`
	RoutingNumber string `json:"routing_number"`
	Branch        string `json:"branch"`
}

type User struct {
	AccountID     string          `json:"account_id"`
	AccountName   string          `json:"account_name"`
	AccountType   string          `json:"account_type"`
	AccountNumber string          `json:"account_number"`
	Balance       UserBalance     `json:"balance"`
	OwnerName     string          `json:"owner_name"`
	BankDetails   UserBankDetails `json:"bank_details"`
} 
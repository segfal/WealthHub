package types

/**

First what we are going to do is add the user account,
the account will contain these values
account_id
account_name
account_type
account_number
balance
	current
	available
	currency
owner_name
bank_details
	bank_name
	routing_number
	branch
*/

// Account represents a user's account information as per init.sql schema
type Account struct {
	AccountID     string  `json:"account_id"`     // VARCHAR(20) PRIMARY KEY
	AccountName   string  `json:"account_name"`   // VARCHAR(50)
	AccountType   string  `json:"account_type"`   // VARCHAR(20)
	AccountNumber string  `json:"account_number"` // VARCHAR(20)
	Balance       Balance `json:"balance"`
	OwnerName     string  `json:"owner_name"`     // VARCHAR(50)
	BankDetails   BankDetails `json:"bank_details"`
}

type Balance struct {
	Current   float64 `json:"current"`    // DECIMAL(10, 2)
	Available float64 `json:"available"`   // DECIMAL(10, 2)
	Currency  string  `json:"currency"`    // VARCHAR(3)
}

type BankDetails struct {
	BankName      string `json:"bank_name"`       // VARCHAR(50)
	RoutingNumber string `json:"routing_number"`  // VARCHAR(20)
	Branch        string `json:"branch"`          // VARCHAR(100)
}
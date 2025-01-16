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


type AccountInformation struct {
	Account_Id     int    `json:"account_id"`
	Account_Name   string `json:"account_name"`
	Account_Type   string `json:"account_type"`
	Account_Number int    `json:"account_number"`
	Owner_Name     string `json:"owner_name"` 
	Balance Balance `json:"balance"` 
	BankDetails BankDetails `json:"bank_details"`
}
 
type Balance struct {  
	Current        float64    `json:"current"`
	Available      float64     `json:"available"`
	Currency       string     `json:"currency"`
}

type BankDetails struct {  
	Bank_Name      string    `json:"bank_name"`
	Routing_Number string    `json:"routing_number"`
	Branch         string `json:"branch"`

}
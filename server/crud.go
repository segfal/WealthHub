package main

import (
	"database/sql"
	"server/types"
)
//create, read, update, delete 

/** @dev	
 * @param db *sql.DB
 * @param accountID string
 * @return []Transaction, error
 */

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


func createUser(db *sql.DB, accountID string){

}



func getTransactions(db *sql.DB, accountID string) ([]Transaction, error) { 
	query := ` 
		SELECT transaction_id, account_id, date, amount, category, merchant, location 
		FROM transactions 
		WHERE account_id = $1
		ORDER BY date DESC` //SQL query to get transactions for a given account ID
	
	rows, err := db.Query(query, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	/** @dev
	 * @param rows *sql.Rows
	 * @return []Transaction, error
	 */

	var transactions []Transaction
	for rows.Next() {
		var t Transaction
		if err := rows.Scan(
			&t.TransactionID,
			&t.AccountID,
			&t.Date,
			&t.Amount,
			&t.Category,
			&t.Merchant,
			&t.Location,
		); err != nil {
			return nil, err
		}
		transactions = append(transactions, t) /// add transaction to transactions slice
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

func insertTransaction(db *sql.DB, transaction *Transaction) error {
	query := `
		INSERT INTO transactions (transaction_id, account_id, date, amount, category, merchant, location)
		VALUES ($1, $2, $3, $4, $5, $6, $7)` // INSERT INTO TABLE(column1, column2, column3, column4, column5, column6, column7) VALUES (value1, value2, value3, value4, value5, value6, value7)
	
	_, err := db.Exec(query,
		transaction.TransactionID,
		transaction.AccountID,
		transaction.Date,
		transaction.Amount,
		transaction.Category,
		transaction.Merchant,
		transaction.Location,
	)
	return err
}


//db *sql.DB: initializes database
func getTransactionsGreaterThan100(db *sql.DB, accountID string) error {
	query := `SELECT * 
	FROM transactions
	WHERE amount >= 100` 
	rows, err := db.Query(query) // talks to the database
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}



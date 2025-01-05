package main

import (
	"database/sql"
)

/** @dev	
 * @param db *sql.DB
 * @param accountID string
 * @return []Transaction, error
 */
func getTransactions(db *sql.DB, accountID string) ([]Transaction, error) {
	query := `
		SELECT transaction_id, account_id, date, amount, category, merchant, location 
		FROM transactions 
		WHERE account_id = $1
		ORDER BY date DESC`
	
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
		transactions = append(transactions, t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

func insertTransaction(db *sql.DB, t *Transaction) error {
	query := `
		INSERT INTO transactions (transaction_id, account_id, date, amount, category, merchant, location)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	
	_, err := db.Exec(query,
		t.TransactionID,
		t.AccountID,
		t.Date,
		t.Amount,
		t.Category,
		t.Merchant,
		t.Location,
	)
	return err
}
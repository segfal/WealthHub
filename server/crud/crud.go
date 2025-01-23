package crud

import (
	"database/sql"
	"fmt"
	"server/types"

	_ "github.com/lib/pq"
)

// CreateTables initializes the database schema
func CreateTables(db *sql.DB) error {
	// Drop existing tables
	dropTables := `
		DROP TABLE IF EXISTS transactions;
		DROP TABLE IF EXISTS bank_details;
		DROP TABLE IF EXISTS users;`
	
	if _, err := db.Exec(dropTables); err != nil {
		return fmt.Errorf("failed to drop tables: %w", err)
	}

	// Create users table
	createUsers := `
		CREATE TABLE users (
			account_id VARCHAR(20) PRIMARY KEY,
			account_name VARCHAR(50),
			account_type VARCHAR(20),
			account_number VARCHAR(20),
			balance_current DECIMAL(10, 2),
			balance_available DECIMAL(10, 2),
			balance_currency VARCHAR(3),
			owner_name VARCHAR(50)
		)`
	
	if _, err := db.Exec(createUsers); err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	// Create bank_details table
	createBankDetails := `
		CREATE TABLE bank_details (
			account_id VARCHAR(20) PRIMARY KEY REFERENCES users(account_id),
			bank_name VARCHAR(50),
			routing_number VARCHAR(20),
			branch VARCHAR(100)
		)`
	
	if _, err := db.Exec(createBankDetails); err != nil {
		return fmt.Errorf("failed to create bank_details table: %w", err)
	}

	// Create transactions table
	createTransactions := `
		CREATE TABLE transactions (
			transaction_id VARCHAR(20) PRIMARY KEY,
			account_id VARCHAR(20) REFERENCES users(account_id),
			date TIMESTAMP,
			amount DECIMAL(10, 2),
			category VARCHAR(50),
			merchant VARCHAR(50),
			location VARCHAR(100)
		)`
	
	if _, err := db.Exec(createTransactions); err != nil {
		return fmt.Errorf("failed to create transactions table: %w", err)
	}

	return nil
}

// createUserTx creates a user within a transaction
func createUserTx(tx *sql.Tx, user *types.User) error {
	query := `
		INSERT INTO users (
			account_id, account_name, account_type, account_number, 
			owner_name, balance_current, balance_available, balance_currency
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := tx.Exec(query,
		user.AccountID,
		user.AccountName,
		user.AccountType,
		user.AccountNumber,
		user.OwnerName,
		user.Balance.Current,
		user.Balance.Available,
		user.Balance.Currency,
	)
	if err != nil {
		return err
	}

	bankQuery := `
		INSERT INTO bank_details (
			account_id, bank_name, routing_number, branch
		) VALUES ($1, $2, $3, $4)`

	_, err = tx.Exec(bankQuery,
		user.AccountID,
		user.BankDetails.BankName,
		user.BankDetails.RoutingNumber,
		user.BankDetails.Branch,
	)
	return err
}

// insertTransactionTx inserts a transaction within a transaction
func insertTransactionTx(tx *sql.Tx, transaction *types.Transaction) error {
	query := `
		INSERT INTO transactions (
			transaction_id, account_id, date, amount, category, merchant, location
		) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	
	_, err := tx.Exec(query,
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

// GetTransactions retrieves all transactions for a given account
func GetTransactions(db *sql.DB, accountID string) ([]types.Transaction, error) {
	query := `
		SELECT transaction_id, account_id, date, amount, category, merchant, location
		FROM transactions
		WHERE account_id = $1
		ORDER BY date DESC`

	rows, err := db.Query(query, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to query transactions: %w", err)
	}
	defer rows.Close()

	var transactions []types.Transaction
	for rows.Next() {
		var t types.Transaction
		if err := rows.Scan(
			&t.TransactionID,
			&t.AccountID,
			&t.Date,
			&t.Amount,
			&t.Category,
			&t.Merchant,
			&t.Location,
		); err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, t)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating transactions: %w", err)
	}

	return transactions, nil
}

// CreateUser creates a new user in the database
func CreateUser(db *sql.DB, user *types.User) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if err := createUserTx(tx, user); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return tx.Commit()
}

// InsertTransaction inserts a new transaction into the database
func InsertTransaction(db *sql.DB, transaction *types.Transaction) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if err := insertTransactionTx(tx, transaction); err != nil {
		return fmt.Errorf("failed to insert transaction: %w", err)
	}

	return tx.Commit()
}

// GetUser retrieves user information by account ID
func GetUser(db *sql.DB, accountID string) (*types.User, error) {
	query := `
		SELECT u.account_id, u.account_name, u.account_type, u.account_number,
			   u.balance_current, u.balance_available, u.balance_currency, u.owner_name,
			   b.bank_name, b.routing_number, b.branch
		FROM users u
		LEFT JOIN bank_details b ON u.account_id = b.account_id
		WHERE u.account_id = $1`

	user := &types.User{}
	err := db.QueryRow(query, accountID).Scan(
		&user.AccountID,
		&user.AccountName,
		&user.AccountType,
		&user.AccountNumber,
		&user.Balance.Current,
		&user.Balance.Available,
		&user.Balance.Currency,
		&user.OwnerName,
		&user.BankDetails.BankName,
		&user.BankDetails.RoutingNumber,
		&user.BankDetails.Branch,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
} 
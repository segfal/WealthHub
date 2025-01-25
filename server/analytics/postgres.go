package analytics

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"server/types"
	"strings"
)

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) Repository {
	if db == nil {
		panic("database connection is required")
	}
	return &postgresRepo{db: db}
}

// GetAccount retrieves account information from the database
func (r *postgresRepo) GetAccount(ctx context.Context, accountID string) (*types.Account, error) {
	if accountID == "" {
		return nil, fmt.Errorf("account ID is required")
	}

	log.Printf("Fetching account information for ID: %s", accountID)

	query := `SELECT account_id, account_name, account_type, account_number, 
	          balance_current, balance_available, balance_currency, owner_name 
	          FROM users WHERE account_id = $1`

	account := &types.Account{
		Balance: types.Balance{},
	}

	err := r.db.QueryRowContext(ctx, query, accountID).Scan(
		&account.AccountID,
		&account.AccountName,
		&account.AccountType,
		&account.AccountNumber,
		&account.Balance.Current,
		&account.Balance.Available,
		&account.Balance.Currency,
		&account.OwnerName,
	)

	if err == sql.ErrNoRows {
		log.Printf("No account found with ID: %s", accountID)
		return nil, fmt.Errorf("account not found")
	}
	if err != nil {
		log.Printf("Error fetching account: %v", err)
		return nil, fmt.Errorf("failed to fetch account: %w", err)
	}

	log.Printf("Successfully retrieved account information for ID: %s", accountID)
	return account, nil
}

func (r *postgresRepo) GetTransactions(ctx context.Context, accountID string, timeRange string) ([]types.Transaction, error) {
	if accountID == "" {
		return nil, fmt.Errorf("account ID is required")
	}

	log.Printf("Fetching transactions for account %s with time range %s", accountID, timeRange)

	query := `
		SELECT transaction_id, account_id, date, amount, category, merchant, location
		FROM transactions 
		WHERE account_id = $1 
		  AND date >= NOW() - $2::INTERVAL
		ORDER BY date DESC`
	
	rows, err := r.db.QueryContext(ctx, query, accountID, timeRange)
	if err != nil {
		log.Printf("Error querying transactions: %v", err)
		return nil, fmt.Errorf("failed to query transactions: %w", err)
	}
	defer rows.Close()

	var transactions []types.Transaction
	for rows.Next() {
		var t types.Transaction
		var prefixedTransactionID string
		if err := rows.Scan(
			&prefixedTransactionID,
			&t.AccountID,
			&t.Date,
			&t.Amount,
			&t.Category,
			&t.Merchant,
			&t.Location,
		); err != nil {
			log.Printf("Error scanning transaction: %v", err)
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}

		// Extract the original transaction ID by removing the prefix
		parts := strings.SplitN(prefixedTransactionID, "_", 2)
		if len(parts) == 2 {
			t.TransactionID = parts[1]
			t.UserPrefix = parts[0]
		} else {
			t.TransactionID = prefixedTransactionID
		}

		transactions = append(transactions, t)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating transactions: %v", err)
		return nil, fmt.Errorf("error iterating transactions: %w", err)
	}

	log.Printf("Found %d transactions for account %s", len(transactions), accountID)
	return transactions, nil
}

func (r *postgresRepo) GetCategoryTotals(ctx context.Context, accountID string, timeRange string) (map[string]float64, error) {
	if accountID == "" {
		return nil, fmt.Errorf("account ID is required")
	}

	log.Printf("Fetching category totals for account %s with time range %s", accountID, timeRange)

	query := `
		SELECT category, COALESCE(SUM(ABS(amount)), 0) as total
		FROM transactions 
		WHERE account_id = $1 
		  AND date >= NOW() - $2::INTERVAL
		GROUP BY category
		ORDER BY total DESC`
	
	rows, err := r.db.QueryContext(ctx, query, accountID, timeRange)
	if err != nil {
		log.Printf("Error querying category totals: %v", err)
		return nil, fmt.Errorf("failed to query category totals: %w", err)
	}
	defer rows.Close()

	categoryTotals := make(map[string]float64)
	for rows.Next() {
		var category string
		var total float64
		if err := rows.Scan(&category, &total); err != nil {
			log.Printf("Error scanning category total: %v", err)
			return nil, fmt.Errorf("failed to scan category total: %w", err)
		}
		categoryTotals[category] = total
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating category totals: %v", err)
		return nil, fmt.Errorf("error iterating category totals: %w", err)
	}

	log.Printf("Found %d categories for account %s", len(categoryTotals), accountID)
	return categoryTotals, nil
}  

//Bills 
func (r *postgresRepo) GetBillTotals(ctx context.Context, accountID string, timeRange string) (map[string]float64, error) {
	if accountID == "" {
		return nil, fmt.Errorf("account ID is required")
	}

	log.Printf("Fetching bill totals for account %s with time range %s", accountID, timeRange)

	query := `
		SELECT category
		FROM transactions 
		WHERE account_id = $1 
		  AND date >= NOW() - $2::INTERVAL AND category = 'Bill Payment'
		`
	
	rows, err := r.db.QueryContext(ctx, query, accountID, timeRange)
	if err != nil {
		log.Printf("Error querying category totals: %v", err)
		return nil, fmt.Errorf("failed to query category totals: %w", err)
	}
	defer rows.Close()

	categoryTotals := make(map[string]float64)
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			log.Printf("Error scanning category total: %v", err)
			return nil, fmt.Errorf("failed to scan category total: %w", err)
		}
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating category totals: %v", err)
		return nil, fmt.Errorf("error iterating category totals: %w", err)
	}

	log.Printf("Found %d categories for account %s", len(categoryTotals), accountID)
	return categoryTotals, nil
} 
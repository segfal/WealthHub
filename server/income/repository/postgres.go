package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"server/types"
	"time"
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

// GetIncome retrieves all income transactions for an account
func (r *postgresRepo) GetIncome(ctx context.Context, accountID string) ([]types.Transaction, error) {
	if accountID == "" {
		return nil, fmt.Errorf("account ID is required")
	}

	log.Printf("Fetching income transactions for account %s", accountID)

	query := `
		SELECT transaction_id, account_id, date, amount, category, merchant, location
		FROM transactions 
		WHERE account_id = $1 
		  AND category = 'Income'
		ORDER BY date DESC`

	rows, err := r.db.QueryContext(ctx, query, accountID)
	if err != nil {
		log.Printf("Error querying income transactions: %v", err)
		return nil, fmt.Errorf("failed to query income transactions: %w", err)
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
			log.Printf("Error scanning income transaction: %v", err)
			return nil, fmt.Errorf("failed to scan income transaction: %w", err)
		}
		transactions = append(transactions, t)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating income transactions: %v", err)
		return nil, fmt.Errorf("error iterating income transactions: %w", err)
	}

	log.Printf("Found %d income transactions for account %s", len(transactions), accountID)
	return transactions, nil
}

// GetMonthlyIncome retrieves income transactions for a specific month
func (r *postgresRepo) GetMonthlyIncome(ctx context.Context, accountID string, year int, month int) ([]types.Transaction, error) {
	if accountID == "" {
		return nil, fmt.Errorf("account ID is required")
	}

	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

	log.Printf("Fetching income for account %s between %s and %s", accountID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))

	query := `
		SELECT transaction_id, account_id, date, amount, category, merchant, location
		FROM transactions 
		WHERE account_id = $1 
		  AND date >= $2
		  AND date <= $3
		  AND category = 'Income'
		ORDER BY date ASC`

	rows, err := r.db.QueryContext(ctx, query, accountID, startDate, endDate)
	if err != nil {
		log.Printf("Error querying monthly income: %v", err)
		return nil, fmt.Errorf("failed to query monthly income: %w", err)
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
			log.Printf("Error scanning income transaction: %v", err)
			return nil, fmt.Errorf("failed to scan income transaction: %w", err)
		}
		transactions = append(transactions, t)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating income transactions: %v", err)
		return nil, fmt.Errorf("error iterating income transactions: %w", err)
	}

	log.Printf("Found %d income transactions for %s", len(transactions), startDate.Format("2006-01"))
	return transactions, nil
} 
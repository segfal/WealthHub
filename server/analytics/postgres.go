package analytics

import (
	"context"
	"database/sql"
	"fmt"
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

func (r *postgresRepo) GetTransactions(ctx context.Context, accountID string, timeRange string) ([]types.Transaction, error) {
	if accountID == "" {
		return nil, fmt.Errorf("account ID is required")
	}

	query := `
		SELECT transaction_id, account_id, date, amount, category, merchant, location
		FROM transactions 
		WHERE account_id = $1 
		  AND date >= NOW() - $2::INTERVAL
		ORDER BY date DESC`
	
	rows, err := r.db.QueryContext(ctx, query, accountID, timeRange)
	if err != nil {
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
		return nil, fmt.Errorf("error iterating transactions: %w", err)
	}

	return transactions, nil
}

func (r *postgresRepo) GetCategoryTotals(ctx context.Context, accountID string, timeRange string) (map[string]float64, error) {
	if accountID == "" {
		return nil, fmt.Errorf("account ID is required")
	}

	query := `
		SELECT category, COALESCE(SUM(ABS(amount)), 0) as total
		FROM transactions 
		WHERE account_id = $1 
		  AND date >= NOW() - $2::INTERVAL
		GROUP BY category
		ORDER BY total DESC`
	
	rows, err := r.db.QueryContext(ctx, query, accountID, timeRange)
	if err != nil {
		return nil, fmt.Errorf("failed to query category totals: %w", err)
	}
	defer rows.Close()

	categoryTotals := make(map[string]float64)
	for rows.Next() {
		var category string
		var total float64
		if err := rows.Scan(&category, &total); err != nil {
			return nil, fmt.Errorf("failed to scan category total: %w", err)
		}
		categoryTotals[category] = total
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating category totals: %w", err)
	}

	return categoryTotals, nil
} 
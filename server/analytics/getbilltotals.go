package analytics

import (
	"context"
	"fmt"
	"log"
	"server/types"
)

// Bills
func (r *postgresRepo) GetBillTotals(ctx context.Context, accountID string, timeRange string) ([]types.Transaction, error) {
	if accountID == "" {
		return nil, fmt.Errorf("account ID is required")
	}

	log.Printf("Fetching bill payments for account %s with time range %s", accountID, timeRange)

	query := `
		SELECT transaction_id, account_id, date, amount, category, merchant, location
		FROM transactions 
		WHERE account_id = $1 
		  AND date >= NOW() - $2::INTERVAL 
		  AND (category = 'Bill Payment' OR category = 'Subscription')
		ORDER BY date DESC`

	rows, err := r.db.QueryContext(ctx, query, accountID, timeRange)
	if err != nil {
		log.Printf("Error querying bill payments: %v", err)
		return nil, fmt.Errorf("failed to query bill payments: %w", err)
	}
	defer rows.Close()

	var billPayments []types.Transaction
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
			log.Printf("Error scanning bill payment: %v", err)
			return nil, fmt.Errorf("failed to scan bill payment: %w", err)
		}
		billPayments = append(billPayments, t)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating bill payments: %v", err)
		return nil, fmt.Errorf("error iterating bill payments: %w", err)
	}

	log.Printf("Found %d bill payments for account %s", len(billPayments), accountID)
	return billPayments, nil
}

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

// GetBillTotals retrieves total bill payments by category for a given time period
func (r *postgresRepo) GetBillTotals(ctx context.Context, accountID string, startDate, endDate time.Time) (map[string]float64, error) {
	if accountID == "" {
		return nil, fmt.Errorf("account ID is required")
	}

	log.Printf("Fetching bill totals for account %s between %s and %s", accountID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))

	query := `
		SELECT merchant, COALESCE(SUM(ABS(amount)), 0) as total
		FROM transactions 
		WHERE account_id = $1 
		  AND date >= $2 
		  AND date <= $3
		  AND (category = 'Bill Payment' OR category = 'Subscription')
		GROUP BY merchant
		ORDER BY total DESC`

	rows, err := r.db.QueryContext(ctx, query, accountID, startDate, endDate)
	if err != nil {
		log.Printf("Error querying bill totals: %v", err)
		return nil, fmt.Errorf("failed to query bill totals: %w", err)
	}
	defer rows.Close()

	billTotals := make(map[string]float64)
	for rows.Next() {
		var merchant string
		var total float64
		if err := rows.Scan(&merchant, &total); err != nil {
			log.Printf("Error scanning bill total: %v", err)
			return nil, fmt.Errorf("failed to scan bill total: %w", err)
		}
		billTotals[merchant] = total
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating bill totals: %v", err)
		return nil, fmt.Errorf("error iterating bill totals: %w", err)
	}

	log.Printf("Found bill totals for %d merchants", len(billTotals))
	return billTotals, nil
}

// GetRecurringBills retrieves recurring bill payments for an account
func (r *postgresRepo) GetRecurringBills(ctx context.Context, accountID string) ([]types.RecurringBill, error) {
	if accountID == "" {
		return nil, fmt.Errorf("account ID is required")
	}

	log.Printf("Fetching recurring bills for account %s", accountID)

	query := `
		WITH monthly_bills AS (
			SELECT 
				merchant,
				category,
				COUNT(DISTINCT DATE_TRUNC('month', date)) as months_present,
				AVG(ABS(amount)) as avg_amount,
				PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY ABS(amount)) as median_amount,
				MIN(date) as first_occurrence,
				MAX(date) as last_occurrence
			FROM transactions
			WHERE account_id = $1
				AND (category = 'Bill Payment' OR category = 'Subscription')
				AND date >= NOW() - INTERVAL '6 months'
			GROUP BY merchant, category
			HAVING COUNT(DISTINCT DATE_TRUNC('month', date)) >= 3
		)
		SELECT 
			merchant,
			category,
			months_present,
			avg_amount,
			median_amount,
			first_occurrence,
			last_occurrence
		FROM monthly_bills
		ORDER BY avg_amount DESC`

	rows, err := r.db.QueryContext(ctx, query, accountID)
	if err != nil {
		log.Printf("Error querying recurring bills: %v", err)
		return nil, fmt.Errorf("failed to query recurring bills: %w", err)
	}
	defer rows.Close()

	var bills []types.RecurringBill
	for rows.Next() {
		var bill types.RecurringBill
		if err := rows.Scan(
			&bill.Merchant,
			&bill.Category,
			&bill.MonthsPresent,
			&bill.AverageAmount,
			&bill.MedianAmount,
			&bill.FirstOccurrence,
			&bill.LastOccurrence,
		); err != nil {
			log.Printf("Error scanning recurring bill: %v", err)
			return nil, fmt.Errorf("failed to scan recurring bill: %w", err)
		}
		bills = append(bills, bill)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating recurring bills: %v", err)
		return nil, fmt.Errorf("error iterating recurring bills: %w", err)
	}

	log.Printf("Found %d recurring bills", len(bills))
	return bills, nil
}

// GetUpcomingBills predicts upcoming bill payments based on recurring patterns
func (r *postgresRepo) GetUpcomingBills(ctx context.Context, accountID string) ([]types.UpcomingBill, error) {
	if accountID == "" {
		return nil, fmt.Errorf("account ID is required")
	}

	log.Printf("Predicting upcoming bills for account %s", accountID)

	query := `
		WITH recurring_bills AS (
			SELECT 
				merchant,
				category,
				AVG(ABS(amount)) as expected_amount,
				MAX(date) as last_payment_date,
				AVG(EXTRACT(DAY FROM date)) as avg_day_of_month
			FROM transactions
			WHERE account_id = $1
				AND (category = 'Bill Payment' OR category = 'Subscription')
				AND date >= NOW() - INTERVAL '6 months'
			GROUP BY merchant, category
			HAVING COUNT(DISTINCT DATE_TRUNC('month', date)) >= 3
		)
		SELECT 
			merchant,
			category,
			expected_amount,
			last_payment_date,
			avg_day_of_month
		FROM recurring_bills
		ORDER BY avg_day_of_month ASC`

	rows, err := r.db.QueryContext(ctx, query, accountID)
	if err != nil {
		log.Printf("Error querying upcoming bills: %v", err)
		return nil, fmt.Errorf("failed to query upcoming bills: %w", err)
	}
	defer rows.Close()

	var bills []types.UpcomingBill
	for rows.Next() {
		var bill types.UpcomingBill
		var lastPaymentDate time.Time
		var avgDayOfMonth float64

		if err := rows.Scan(
			&bill.Merchant,
			&bill.Category,
			&bill.ExpectedAmount,
			&lastPaymentDate,
			&avgDayOfMonth,
		); err != nil {
			log.Printf("Error scanning upcoming bill: %v", err)
			return nil, fmt.Errorf("failed to scan upcoming bill: %w", err)
		}

		// Calculate next due date based on last payment and average day of month
		nextDueDate := lastPaymentDate.AddDate(0, 1, 0)
		bill.DueDate = time.Date(nextDueDate.Year(), nextDueDate.Month(), int(avgDayOfMonth), 0, 0, 0, 0, time.UTC)
		
		bills = append(bills, bill)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating upcoming bills: %v", err)
		return nil, fmt.Errorf("error iterating upcoming bills: %w", err)
	}

	log.Printf("Predicted %d upcoming bills", len(bills))
	return bills, nil
}

// GetBillHistory retrieves historical bill payments for a specific merchant
func (r *postgresRepo) GetBillHistory(ctx context.Context, accountID string, merchantName string) ([]types.Transaction, error) {
	if accountID == "" || merchantName == "" {
		return nil, fmt.Errorf("account ID and merchant name are required")
	}

	log.Printf("Fetching bill history for account %s and merchant %s", accountID, merchantName)

	query := `
		SELECT transaction_id, account_id, date, amount, category, merchant, location
		FROM transactions 
		WHERE account_id = $1 
		  AND merchant = $2
		  AND (category = 'Bill Payment' OR category = 'Subscription')
		ORDER BY date DESC`

	rows, err := r.db.QueryContext(ctx, query, accountID, merchantName)
	if err != nil {
		log.Printf("Error querying bill history: %v", err)
		return nil, fmt.Errorf("failed to query bill history: %w", err)
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
			log.Printf("Error scanning transaction: %v", err)
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, t)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating transactions: %v", err)
		return nil, fmt.Errorf("error iterating transactions: %w", err)
	}

	log.Printf("Found %d bill payments for merchant %s", len(transactions), merchantName)
	return transactions, nil
}

// GetBillsByMonth retrieves all bill payments for a specific month
func (r *postgresRepo) GetBillsByMonth(ctx context.Context, accountID string, year int, month int) ([]types.Transaction, error) {
	if accountID == "" {
		return nil, fmt.Errorf("account ID is required")
	}

	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

	log.Printf("Fetching bills for account %s for %s", accountID, startDate.Format("2006-01"))

	query := `
		SELECT transaction_id, account_id, date, amount, category, merchant, location
		FROM transactions 
		WHERE account_id = $1 
		  AND date >= $2
		  AND date <= $3
		  AND (category = 'Bill Payment' OR category = 'Subscription')
		ORDER BY date ASC`

	rows, err := r.db.QueryContext(ctx, query, accountID, startDate, endDate)
	if err != nil {
		log.Printf("Error querying monthly bills: %v", err)
		return nil, fmt.Errorf("failed to query monthly bills: %w", err)
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
			log.Printf("Error scanning transaction: %v", err)
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, t)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating transactions: %v", err)
		return nil, fmt.Errorf("error iterating transactions: %w", err)
	}

	log.Printf("Found %d bills for %s", len(transactions), startDate.Format("2006-01"))
	return transactions, nil
} 
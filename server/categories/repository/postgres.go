package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"server/types"
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

// GetCategories retrieves all categories for an account
func (r *postgresRepo) GetCategories(ctx context.Context, accountID string) ([]types.Category, error) {
	if accountID == "" {
		return nil, fmt.Errorf("account ID is required")
	}

	log.Printf("Fetching categories for account %s", accountID)

	query := `
		SELECT 
			category as id,
			category as name,
			'' as description,
			COALESCE(SUM(ABS(amount)), 0) as total_spent,
			COUNT(*) as count
		FROM transactions 
		WHERE account_id = $1 
		GROUP BY category
		ORDER BY total_spent DESC`

	rows, err := r.db.QueryContext(ctx, query, accountID)
	if err != nil {
		log.Printf("Error querying categories: %v", err)
		return nil, fmt.Errorf("failed to query categories: %w", err)
	}
	defer rows.Close()

	var categories []types.Category
	for rows.Next() {
		var c types.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.TotalSpent, &c.Count); err != nil {
			log.Printf("Error scanning category: %v", err)
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, c)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating categories: %v", err)
		return nil, fmt.Errorf("error iterating categories: %w", err)
	}

	log.Printf("Found %d categories for account %s", len(categories), accountID)
	return categories, nil
}

// GetCategoryTotals retrieves total spending by category
func (r *postgresRepo) GetCategoryTotals(ctx context.Context, accountID string) (map[string]float64, error) {
	if accountID == "" {
		return nil, fmt.Errorf("account ID is required")
	}

	log.Printf("Fetching category totals for account %s", accountID)

	query := `
		SELECT category, COALESCE(SUM(ABS(amount)), 0) as total
		FROM transactions 
		WHERE account_id = $1 
		GROUP BY category
		ORDER BY total DESC`

	rows, err := r.db.QueryContext(ctx, query, accountID)
	if err != nil {
		log.Printf("Error querying category totals: %v", err)
		return nil, fmt.Errorf("failed to query category totals: %w", err)
	}
	defer rows.Close()

	totals := make(map[string]float64)
	for rows.Next() {
		var category string
		var total float64
		if err := rows.Scan(&category, &total); err != nil {
			log.Printf("Error scanning category total: %v", err)
			return nil, fmt.Errorf("failed to scan category total: %w", err)
		}
		totals[category] = total
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating category totals: %v", err)
		return nil, fmt.Errorf("error iterating category totals: %w", err)
	}

	log.Printf("Found totals for %d categories", len(totals))
	return totals, nil
} 
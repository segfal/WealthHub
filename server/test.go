// test.go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	analyticsRepo "server/analytics/repository"
	"testing"

	_ "github.com/lib/pq"
)


func setupTestDB(t *testing.T) *sql.DB {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		// Fallback to constructing the URL from individual components
		dbURL = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
		)
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}

	return db
}

func TestGetAccount(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := analyticsRepo.NewPostgresRepository(db)
	ctx := context.Background()

	tests := []struct {
		name      string
		accountID string
		wantErr   bool
	}{
		{
			name:      "Valid account",
			accountID: "test_account_1",
			wantErr:   false,
		},
		{
			name:      "Empty account ID",
			accountID: "",
			wantErr:   true,
		},
		{
			name:      "Non-existent account",
			accountID: "non_existent",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account, err := repo.GetAccount(ctx, tt.accountID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && account == nil {
				t.Error("GetAccount() returned nil account for valid ID")
			}
		})
	}
}

func TestGetTransactions(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := analyticsRepo.NewPostgresRepository(db)
	ctx := context.Background()

	tests := []struct {
		name      string
		accountID string
		timeRange string
		wantErr   bool
	}{
		{
			name:      "Valid transactions last month",
			accountID: "test_account_1",
			timeRange: "1 month",
			wantErr:   false,
		},
		{
			name:      "Empty account ID",
			accountID: "",
			timeRange: "1 month",
			wantErr:   true,
		},
		{
			name:      "Invalid time range",
			accountID: "test_account_1",
			timeRange: "invalid",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transactions, err := repo.GetTransactions(ctx, tt.accountID, tt.timeRange)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTransactions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && transactions == nil {
				t.Error("GetTransactions() returned nil transactions for valid parameters")
			}
		})
	}
}

func TestGetCategoryTotals(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := analyticsRepo.NewPostgresRepository(db)
	ctx := context.Background()

	tests := []struct {
		name      string
		accountID string
		timeRange string
		wantErr   bool
	}{
		{
			name:      "Valid category totals",
			accountID: "test_account_1",
			timeRange: "1 month",
			wantErr:   false,
		},
		{
			name:      "Empty account ID",
			accountID: "",
			timeRange: "1 month",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			totals, err := repo.GetCategoryTotals(ctx, tt.accountID, tt.timeRange)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCategoryTotals() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && totals == nil {
				t.Error("GetCategoryTotals() returned nil totals for valid parameters")
			}
		})
	}
}


package crud_test

import (
	"server/crud"
	"server/tests/testutils"
	"server/types"
	"testing"
	"time"
)

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name    string
		user    *types.User
		wantErr bool
	}{
		{
			name: "valid user",
			user: &types.User{
				AccountID:     "test123",
				AccountName:   "Test Account",
				AccountType:   "checking",
				AccountNumber: "1234567890",
				Balance: types.UserBalance{
					Current:   1000.00,
					Available: 900.00,
					Currency:  "USD",
				},
				OwnerName: "Test User",
				BankDetails: types.UserBankDetails{
					BankName:      "Test Bank",
					RoutingNumber: "987654321",
					Branch:        "Test Branch",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := testutils.GetTestDB(t)
			defer db.Close()

			if err := crud.CreateUser(db, tt.user); (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInsertTransaction(t *testing.T) {
	tests := []struct {
		name        string
		transaction *types.Transaction
		wantErr     bool
	}{
		{
			name: "valid transaction",
			transaction: &types.Transaction{
				TransactionID: "txn123",
				AccountID:     "test123",
				Date:         time.Now(),
				Amount:       100.00,
				Category:     "groceries",
				Merchant:     "Test Store",
				Location:     "Test Location",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := testutils.GetTestDB(t)
			defer db.Close()

			if err := crud.InsertTransaction(db, tt.transaction); (err != nil) != tt.wantErr {
				t.Errorf("InsertTransaction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	tests := []struct {
		name      string
		accountID string
		wantErr   bool
	}{
		{
			name:      "existing user",
			accountID: "test123",
			wantErr:   false,
		},
		{
			name:      "non-existing user",
			accountID: "nonexistent",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := testutils.GetTestDB(t)
			defer db.Close()

			got, err := crud.GetUser(db, tt.accountID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Error("GetUser() returned nil user")
			}
		})
	}
} 
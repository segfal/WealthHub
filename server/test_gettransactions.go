package main

import (
	"database/sql"
	"testing"
)

func Test_getTransactions(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		db        *sql.DB
		accountID string
		want      []Transaction
		wantErr   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := getTransactions(tt.db, tt.accountID)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("getTransactions() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("getTransactions() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("getTransactions() = %v, want %v", got, tt.want)
			}
		})
	}
}

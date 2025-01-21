package crud_test

import (
	"database/sql"
	"server/crud"
	"server/types"
	"testing"
)

func TestInsertJaneData(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		db      *sql.DB
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := crud.InsertJaneData(tt.db)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("InsertJaneData() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("InsertJaneData() succeeded unexpectedly")
			}
		})
	}
}

func TestInsertJaneData(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		db      *sql.DB
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := crud.InsertJaneData(tt.db)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("InsertJaneData() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("InsertJaneData() succeeded unexpectedly")
			}
		})
	}
}

func TestInsertJaneData(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		db      *sql.DB
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := crud.InsertJaneData(tt.db)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("InsertJaneData() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("InsertJaneData() succeeded unexpectedly")
			}
		})
	}
}

func TestInsertJaneData(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		db      *sql.DB
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := crud.InsertJaneData(tt.db)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("InsertJaneData() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("InsertJaneData() succeeded unexpectedly")
			}
		})
	}
}

func TestInsertJaneData(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		db      *sql.DB
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := crud.InsertJaneData(tt.db)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("InsertJaneData() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("InsertJaneData() succeeded unexpectedly")
			}
		})
	}
}

func TestInsertJaneData(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		db      *sql.DB
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := crud.InsertJaneData(tt.db)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("InsertJaneData() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("InsertJaneData() succeeded unexpectedly")
			}
		})
	}
}

func TestInsertJaneData(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		db      *sql.DB
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := crud.InsertJaneData(tt.db)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("InsertJaneData() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("InsertJaneData() succeeded unexpectedly")
			}
		})
	}
}

func TestInsertJaneData(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		db      *sql.DB
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := crud.InsertJaneData(tt.db)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("InsertJaneData() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("InsertJaneData() succeeded unexpectedly")
			}
		})
	}
}

func TestInsertJaneData(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		db      *sql.DB
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := crud.InsertJaneData(tt.db)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("InsertJaneData() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("InsertJaneData() succeeded unexpectedly")
			}
		})
	}
}

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		db      *sql.DB
		user    *types.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := crud.CreateUser(tt.db, tt.user)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("CreateUser() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("CreateUser() succeeded unexpectedly")
			}
		})
	}
}

func TestCreateTables(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		db      *sql.DB
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := crud.CreateTables(tt.db)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("CreateTables() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("CreateTables() succeeded unexpectedly")
			}
		})
	}
}

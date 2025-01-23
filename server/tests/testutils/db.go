package testutils

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func init() {
	// Load .env file for tests
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal("Error loading .env file:", err)
	}
}

// GetTestDB returns a database connection for testing
func GetTestDB(t *testing.T) *sql.DB {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		t.Fatal("DB_URL environment variable is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		t.Fatal(err)
	}

	return db
} 
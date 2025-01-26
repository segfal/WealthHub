package income

import (
	"database/sql"
	"encoding/json"
	"math"
	"net/http"
	"server/analytics"
	"sort"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// SetupIncomeRoutes configures all the income-related routes for the API
func SetupIncomeRoutes(router *mux.Router, db *sql.DB) {
	repo := analytics.NewPostgresRepository(db)

	// Monthly Income route
	router.HandleFunc("/api/income/{accountId}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		accountID := vars["accountId"]
		
		// Get year and month from query params, default to current month
		now := time.Now()
		yearStr := r.URL.Query().Get("year")
		monthStr := r.URL.Query().Get("month")

		year := now.Year()
		month := int(now.Month())

		if yearStr != "" {
			var err error
			year, err = strconv.Atoi(yearStr)
			if err != nil {
				http.Error(w, "Invalid year parameter", http.StatusBadRequest)
				return
			}
		}

		if monthStr != "" {
			var err error
			month, err = strconv.Atoi(monthStr)
			if err != nil || month < 1 || month > 12 {
				http.Error(w, "Invalid month parameter (must be 1-12)", http.StatusBadRequest)
				return
			}
		}

		transactions, err := repo.GetMonthlyIncome(r.Context(), accountID, year, month)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Calculate total income
		var totalIncome float64
		for _, transaction := range transactions {
			totalIncome += math.Abs(transaction.Amount)
		}

		// Sort transactions by date (newest first)
		sort.Slice(transactions, func(i, j int) bool {
			return transactions[i].Date.After(transactions[j].Date)
		})

		response := map[string]interface{}{
			"transactions": transactions,
			"totalIncome": totalIncome,
			"year":        year,
			"month":       month,
			"monthName":   time.Month(month).String(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}).Methods("GET")
} 
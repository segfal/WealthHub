package bills

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"server/analytics"
	"sort"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// SetupBillRoutes configures all the bill-related routes for the API
func SetupBillRoutes(router *mux.Router, db *sql.DB) {
	repo := analytics.NewPostgresRepository(db)

	// Bills Route 
	router.HandleFunc("/api/bills/{accountId}", func(w http.ResponseWriter, r *http.Request) {
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

		billPayments, err := repo.GetBillPayments(r.Context(), accountID, year, month)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Transform the data into the expected format
		type BillPayment struct { 
			Category    string  `json:"category"`
			TotalSpent  float64 `json:"totalSpent"`
			Percentage  string  `json:"percentage"`  
		}

		// Calculate totals by merchant
		merchantTotals := make(map[string]float64)
		var totalSpent float64
		
		for _, payment := range billPayments {
			amount := math.Abs(payment.Amount)
			merchantTotals[payment.Merchant] += amount
			totalSpent += amount
		}

		var topBills []BillPayment
		for merchant, amount := range merchantTotals {
			percentage := (amount / totalSpent) * 100
			topBills = append(topBills, BillPayment{
				Category:   merchant,
				TotalSpent: amount,
				Percentage: fmt.Sprintf("%.2f", percentage),
			})
		}

		// Sort bills by amount spent
		sort.Slice(topBills, func(i, j int) bool {
			return topBills[i].TotalSpent > topBills[j].TotalSpent
		})

		response := map[string]interface{}{
			"topBills": topBills,
			"totalSpent": totalSpent,
			"year": year,
			"month": month,
			"monthName": time.Month(month).String(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}).Methods("GET")

	// Combined Bills and Income Analysis route
	router.HandleFunc("/api/analysis/bills-income/{accountId}", func(w http.ResponseWriter, r *http.Request) {
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

		// Get income transactions
		incomeTransactions, err := repo.GetMonthlyIncome(r.Context(), accountID, year, month)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Calculate total income
		var totalIncome float64
		for _, transaction := range incomeTransactions {
			totalIncome += math.Abs(transaction.Amount)
		}

		// Get bill payments
		billPayments, err := repo.GetBillPayments(r.Context(), accountID, year, month)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Calculate bill totals and percentages
		type BillAnalysis struct {
			Merchant    string  `json:"merchant"`
			Amount     float64 `json:"amount"`
			Percentage float64 `json:"percentage"`
			Date       string  `json:"date"`
		}

		merchantTotals := make(map[string]*BillAnalysis)
		var totalBills float64

		for _, payment := range billPayments {
			amount := math.Abs(payment.Amount)
			if analysis, exists := merchantTotals[payment.Merchant]; exists {
				analysis.Amount += amount
			} else {
				merchantTotals[payment.Merchant] = &BillAnalysis{
					Merchant: payment.Merchant,
					Amount:   amount,
					Date:     payment.Date.Format("2006-01-02"),
				}
			}
			totalBills += amount
		}

		// Calculate percentages of income for each bill
		var billsAnalysis []BillAnalysis
		for _, analysis := range merchantTotals {
			analysis.Percentage = (analysis.Amount / totalIncome) * 100
			billsAnalysis = append(billsAnalysis, *analysis)
		}

		// Sort bills by percentage (highest to lowest)
		sort.Slice(billsAnalysis, func(i, j int) bool {
			return billsAnalysis[i].Percentage > billsAnalysis[j].Percentage
		})

		response := map[string]interface{}{
			"monthlyIncome": totalIncome,
			"totalBills": totalBills,
			"billsToIncomeRatio": (totalBills / totalIncome) * 100,
			"remainingIncome": totalIncome - totalBills,
			"remainingIncomePercentage": 100 - ((totalBills / totalIncome) * 100),
			"bills": billsAnalysis,
			"year": year,
			"month": month,
			"monthName": time.Month(month).String(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}).Methods("GET")
} 
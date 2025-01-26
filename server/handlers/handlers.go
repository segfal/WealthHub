package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"server/analytics"
	"server/crud"
	"sort"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// SetupRoutes configures all the routes for the API
func SetupRoutes(router *mux.Router, db *sql.DB) {
	repo := analytics.NewPostgresRepository(db)
	analyticsService := analytics.NewService(repo)

	// User route
	router.HandleFunc("/api/user/{accountId}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		accountID := vars["accountId"]

		user, err := crud.GetUser(db, accountID)
		if err != nil {
			if err.Error() == "user not found" {
				http.Error(w, "User not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}).Methods("GET")

	// Analytics route
	router.HandleFunc("/api/analytics/{accountId}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		accountID := vars["accountId"]
		timeRange := r.URL.Query().Get("timeRange")
		if timeRange == "" {
			timeRange = "1 month"
		}

		analytics, err := analyticsService.AnalyzeSpending(r.Context(), accountID, timeRange)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(analytics)
	}).Methods("GET")

	// Categories route
	router.HandleFunc("/api/categories/{accountId}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		accountID := vars["accountId"]
		
		categoryTotals, err := repo.GetCategoryTotals(r.Context(), accountID, "1 month")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Transform the data into the expected format
		type CategoryTotal struct {
			Category   string  `json:"category"`
			TotalSpent float64 `json:"totalSpent"`
		}

		var topCategories []CategoryTotal
		for category, amount := range categoryTotals {
			topCategories = append(topCategories, CategoryTotal{
				Category:   category,
				TotalSpent: amount,
			})
		}

		response := map[string]interface{}{
			"topCategories": topCategories,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}).Methods("GET")

	//Bills Route 
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

	// Predictions route
	router.HandleFunc("/api/predictions/{accountId}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		accountID := vars["accountId"]
		
		predictions, err := analyticsService.PredictSpending(r.Context(), accountID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(predictions)
	}).Methods("GET")

	// Patterns route
	router.HandleFunc("/api/patterns/{accountId}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		accountID := vars["accountId"]
		
		patterns, err := analyticsService.GetTimePatterns(r.Context(), accountID, time.Now().AddDate(0, -1, 0), time.Now())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(patterns)
	}).Methods("GET")

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

	// Insights route
	router.HandleFunc("/api/insights/{accountId}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		accountID := vars["accountId"]
		
		// Get transactions for the past month
		transactions, err := repo.GetTransactions(r.Context(), accountID, "1 month")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Exclude utilities and rent from insights
		excludedCategories := map[string]bool{
			"Rent": true,
			"Utilities": true,
			"Insurance": true,
			"Phone Bill": true,
			"Internet": true,
			"Mortgage": true,
			"Water Bill": true,
			"Electric Bill": true,
			"Gas Bill": true,
		}

		// Group transactions by category
		categoryTotals := make(map[string]float64)
		totalSpent := 0.0
		for _, txn := range transactions {
			if !excludedCategories[txn.Category] {
				amount := math.Abs(txn.Amount)
				categoryTotals[txn.Category] += amount
				totalSpent += amount
			}
		}

		// Convert to sorted slice for top categories
		type CategoryTotal struct {
			Category   string  `json:"category"`
			TotalSpent float64 `json:"totalSpent"`
		}

		var topCategories []CategoryTotal
		for category, amount := range categoryTotals {
			topCategories = append(topCategories, CategoryTotal{
				Category:   category,
				TotalSpent: amount,
			})
		}

		// Sort by amount spent
		sort.Slice(topCategories, func(i, j int) bool {
			return topCategories[i].TotalSpent > topCategories[j].TotalSpent
		})

		// Get top 5 categories if available
		if len(topCategories) > 5 {
			topCategories = topCategories[:5]
		}

		// Calculate percentages
		var insightData []map[string]string
		for _, cat := range topCategories {
			percentage := (cat.TotalSpent / totalSpent) * 100
			insightData = append(insightData, map[string]string{
				"category":   cat.Category,
				"totalSpent": fmt.Sprintf("%.2f", cat.TotalSpent),
				"percentage": fmt.Sprintf("%.2f", percentage),
			})
		}

		insights := []map[string]interface{}{
			{
				"type":        "top_categories",
				"title":       "Top Spending Categories",
				"description": "Your highest spending areas (excluding bills)",
				"data":        insightData,
			},
		}

		response := map[string]interface{}{
			"insights":     insights,
			"totalSpent":   totalSpent,
			"monthlyAverage": totalSpent,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}).Methods("GET")
}

/*
	to call the api from the client, we need to use the following url:
	http://localhost:8080/api/user/1234567891
	if you want a prediction, you need to use the following url:
	http://localhost:8080/api/predictions/1234567891

*/
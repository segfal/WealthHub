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
		
		billPayments, err := repo.GetBillTotals(r.Context(), accountID, "1 month")
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

		response := map[string]interface{}{
			"topBills": topBills,
			"totalSpent": totalSpent,
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

	// Bills route
	router.HandleFunc("/api/bills/{accountId}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		accountID := vars["accountId"]
		
		// Get transactions for the past month
		transactions, err := repo.GetTransactions(r.Context(), accountID, "1 month")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Filter for bill-related categories
		billCategories := map[string]bool{
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

		var bills []map[string]interface{}
		totalAmount := 0.0

		for _, txn := range transactions {
			if billCategories[txn.Category] {
				amount := math.Abs(txn.Amount)
				totalAmount += amount
				
				// Determine bill status based on date
				var status string
				dueDate := time.Now().AddDate(0, 1, 0) // Due next month
				if txn.Date.After(time.Now()) {
					status = "upcoming"
				} else {
					status = "paid"
				}

				bills = append(bills, map[string]interface{}{
					"category": txn.Category,
					"amount":   fmt.Sprintf("%.2f", amount),
					"dueDate": dueDate.Format(time.RFC3339),
					"status":  status,
				})
			}
		}

		response := map[string]interface{}{
			"bills": bills,
			"monthlyTotal": totalAmount,
			"totalBills": len(bills),
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
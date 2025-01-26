package analytics

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers the analytics routes with the router
func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/analytics/spending", h.HandleSpendingAnalytics).Methods("GET")
	router.HandleFunc("/api/analytics/patterns", h.HandleTimePatterns).Methods("GET")
	router.HandleFunc("/api/analytics/predictions", h.HandlePredictions).Methods("GET")
	router.HandleFunc("/api/analytics/income/monthly", h.HandleMonthlyIncome).Methods("GET")
	router.HandleFunc("/api/analytics/bills/monthly", h.HandleMonthlyBills).Methods("GET")
}

// HandleSpendingAnalytics handles requests for spending analytics
func (h *Handler) HandleSpendingAnalytics(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling spending analytics request: %s", r.URL.String())
	
	accountID := r.URL.Query().Get("account_id")
	if accountID == "" {
		log.Printf("Missing account_id parameter")
		http.Error(w, "account_id is required", http.StatusBadRequest)
		return
	}

	timeRange := r.URL.Query().Get("time_range")
	if timeRange == "" {
		timeRange = "1 month" // default to 1 month if not specified
		log.Printf("Using default time range: %s", timeRange)
	}

	analytics, err := h.service.AnalyzeSpending(r.Context(), accountID, timeRange)
	if err != nil {
		log.Printf("Error analyzing spending: %v", err)
		http.Error(w, "Failed to analyze spending", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully analyzed spending for account %s", accountID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(analytics)
}

// HandleTimePatterns handles requests for time-based spending patterns
func (h *Handler) HandleTimePatterns(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling time patterns request: %s", r.URL.String())
	
	accountID := r.URL.Query().Get("account_id")
	if accountID == "" {
		log.Printf("Missing account_id parameter")
		http.Error(w, "account_id is required", http.StatusBadRequest)
		return
	}

	// Default to last month if no dates provided
	endDate := time.Now()
	startDate := endDate.AddDate(0, -1, 0)
	log.Printf("Using date range: %s to %s", startDate.Format(time.RFC3339), endDate.Format(time.RFC3339))

	patterns, err := h.service.GetTimePatterns(r.Context(), accountID, startDate, endDate)
	if err != nil {
		log.Printf("Error getting time patterns: %v", err)
		http.Error(w, "Failed to get time patterns", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully retrieved time patterns for account %s", accountID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(patterns)
}

// HandlePredictions handles requests for spending predictions
func (h *Handler) HandlePredictions(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling predictions request: %s", r.URL.String())
	
	accountID := r.URL.Query().Get("account_id")
	if accountID == "" {
		log.Printf("Missing account_id parameter")
		http.Error(w, "account_id is required", http.StatusBadRequest)
		return
	}

	predictions, err := h.service.PredictSpending(r.Context(), accountID)
	if err != nil {
		log.Printf("Error predicting spending: %v", err)
		http.Error(w, "Failed to predict spending", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully generated predictions for account %s", accountID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(predictions)
}

// HandleMonthlyIncome handles requests for monthly income data
func (h *Handler) HandleMonthlyIncome(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling monthly income request: %s", r.URL.String())
	
	accountID := r.URL.Query().Get("account_id")
	if accountID == "" {
		log.Printf("Missing account_id parameter")
		http.Error(w, "account_id is required", http.StatusBadRequest)
		return
	}

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
			log.Printf("Invalid year parameter: %s", yearStr)
			http.Error(w, "Invalid year parameter", http.StatusBadRequest)
			return
		}
	}

	if monthStr != "" {
		var err error
		month, err = strconv.Atoi(monthStr)
		if err != nil || month < 1 || month > 12 {
			log.Printf("Invalid month parameter: %s", monthStr)
			http.Error(w, "Invalid month parameter (must be 1-12)", http.StatusBadRequest)
			return
		}
	}

	transactions, err := h.service.GetMonthlyIncome(r.Context(), accountID, year, month)
	if err != nil {
		log.Printf("Error getting monthly income: %v", err)
		http.Error(w, "Failed to get monthly income", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"transactions": transactions,
		"year":        year,
		"month":       month,
	}

	log.Printf("Successfully retrieved monthly income for account %s", accountID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleMonthlyBills handles requests for monthly bill payment data
func (h *Handler) HandleMonthlyBills(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling monthly bills request: %s", r.URL.String())
	
	accountID := r.URL.Query().Get("account_id")
	if accountID == "" {
		log.Printf("Missing account_id parameter")
		http.Error(w, "account_id is required", http.StatusBadRequest)
		return
	}

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
			log.Printf("Invalid year parameter: %s", yearStr)
			http.Error(w, "Invalid year parameter", http.StatusBadRequest)
			return
		}
	}

	if monthStr != "" {
		var err error
		month, err = strconv.Atoi(monthStr)
		if err != nil || month < 1 || month > 12 {
			log.Printf("Invalid month parameter: %s", monthStr)
			http.Error(w, "Invalid month parameter (must be 1-12)", http.StatusBadRequest)
			return
		}
	}

	transactions, err := h.service.GetBillPayments(r.Context(), accountID, year, month)
	if err != nil {
		log.Printf("Error getting monthly bills: %v", err)
		http.Error(w, "Failed to get monthly bills", http.StatusInternalServerError)
		return
	}

	// Calculate totals by merchant
	merchantTotals := make(map[string]float64)
	var totalSpent float64
	
	for _, payment := range transactions {
		amount := math.Abs(payment.Amount)
		merchantTotals[payment.Merchant] += amount
		totalSpent += amount
	}

	// Transform to response format
	var billDetails []map[string]interface{}
	for merchant, amount := range merchantTotals {
		percentage := (amount / totalSpent) * 100
		billDetails = append(billDetails, map[string]interface{}{
			"merchant":   merchant,
			"amount":     amount,
			"percentage": fmt.Sprintf("%.2f", percentage),
		})
	}

	response := map[string]interface{}{
		"transactions": transactions,
		"billDetails": billDetails,
		"totalSpent":  totalSpent,
		"year":        year,
		"month":       month,
	}

	log.Printf("Successfully retrieved monthly bills for account %s", accountID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
} 
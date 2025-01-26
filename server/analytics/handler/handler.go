package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"server/analytics/repository"
	"server/analytics/service"
	"time"

	"github.com/gorilla/mux"
)

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{service: service}
}

// SetupRoutes configures all the analytics-related routes for the API
func SetupRoutes(router *mux.Router, db *sql.DB) {
	repo := repository.NewPostgresRepository(db)
	svc := service.NewService(repo)
	handler := NewHandler(svc)

	// Register all routes
	handler.RegisterRoutes(router)
}

// RegisterRoutes registers all analytics routes
func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/analytics/{accountId}", h.HandleSpendingAnalytics).Methods("GET")
	router.HandleFunc("/api/predictions/{accountId}", h.HandlePredictions).Methods("GET")
	router.HandleFunc("/api/patterns/{accountId}", h.HandleTimePatterns).Methods("GET")
	router.HandleFunc("/api/insights/{accountId}", h.HandleInsights).Methods("GET")
}

// HandleSpendingAnalytics handles requests for spending analytics
func (h *Handler) HandleSpendingAnalytics(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling spending analytics request: %s", r.URL.String())
	
	vars := mux.Vars(r)
	accountID := vars["accountId"]
	timeRange := r.URL.Query().Get("timeRange")
	if timeRange == "" {
		timeRange = "1 month"
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
	
	vars := mux.Vars(r)
	accountID := vars["accountId"]

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
	
	vars := mux.Vars(r)
	accountID := vars["accountId"]

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

// HandleInsights handles requests for spending insights
func (h *Handler) HandleInsights(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling insights request: %s", r.URL.String())
	
	vars := mux.Vars(r)
	accountID := vars["accountId"]

	// Get analytics for the past month
	analytics, err := h.service.AnalyzeSpending(r.Context(), accountID, "1 month")
	if err != nil {
		log.Printf("Error getting analytics: %v", err)
		http.Error(w, "Failed to get analytics", http.StatusInternalServerError)
		return
	}

	// Transform analytics data into insights format
	var insightData []map[string]string
	for _, cat := range analytics.TopCategories {
		insightData = append(insightData, map[string]string{
			"category":   cat.Category,
			"totalSpent": cat.TotalSpent,
			"percentage": cat.Percentage,
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
		"insights":       insights,
		"totalSpent":     analytics.TotalSpent,
		"monthlyAverage": analytics.MonthlyAverage,
	}

	log.Printf("Successfully generated insights for account %s", accountID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
} 


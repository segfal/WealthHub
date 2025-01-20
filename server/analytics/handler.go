package analytics

import (
	"encoding/json"
	"log"
	"net/http"
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
}

// HandleSpendingAnalytics handles requests for spending analytics
func (h *Handler) HandleSpendingAnalytics(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling spending analytics request: %s", r.URL.String())
	
	accountID := r.URL.Query().Get("accountId")
	if accountID == "" {
		log.Printf("Missing accountId parameter")
		http.Error(w, "accountId is required", http.StatusBadRequest)
		return
	}

	timeRange := r.URL.Query().Get("timeRange")
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
	
	accountID := r.URL.Query().Get("accountId")
	if accountID == "" {
		log.Printf("Missing accountId parameter")
		http.Error(w, "accountId is required", http.StatusBadRequest)
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
	
	accountID := r.URL.Query().Get("accountId")
	if accountID == "" {
		log.Printf("Missing accountId parameter")
		http.Error(w, "accountId is required", http.StatusBadRequest)
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
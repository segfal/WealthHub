package analytics

import (
	"encoding/json"
	"net/http"
	"time"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	if service == nil {
		panic("service is required")
	}
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/analytics/", h.handleAnalytics)
	mux.HandleFunc("/api/patterns/", h.handlePatterns)
	mux.HandleFunc("/api/predictions/", h.handlePredictions)
}

func (h *Handler) handleAnalytics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	accountID := r.URL.Path[len("/api/analytics/"):]
	if accountID == "" {
		http.Error(w, "Account ID is required", http.StatusBadRequest)
		return
	}

	timeRange := r.URL.Query().Get("timeRange")
	if timeRange == "" {
		timeRange = "1 month" // Default time range
	}

	analytics, err := h.service.GetSpendingAnalytics(r.Context(), accountID, timeRange)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(analytics)
}

func (h *Handler) handlePatterns(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	accountID := r.URL.Path[len("/api/patterns/"):]
	if accountID == "" {
		http.Error(w, "Account ID is required", http.StatusBadRequest)
		return
	}

	// Default to last month if no dates provided
	endDate := time.Now()
	startDate := endDate.AddDate(0, -1, 0)

	// Parse date parameters if provided
	if start := r.URL.Query().Get("start"); start != "" {
		parsedStart, err := time.Parse("2006-01-02", start)
		if err == nil {
			startDate = parsedStart
		}
	}
	if end := r.URL.Query().Get("end"); end != "" {
		parsedEnd, err := time.Parse("2006-01-02", end)
		if err == nil {
			endDate = parsedEnd
		}
	}

	patterns, err := h.service.AnalyzeTimePatterns(r.Context(), accountID, startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(patterns)
}

func (h *Handler) handlePredictions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	accountID := r.URL.Path[len("/api/predictions/"):]
	if accountID == "" {
		http.Error(w, "Account ID is required", http.StatusBadRequest)
		return
	}

	predictions, err := h.service.PredictFutureSpending(r.Context(), accountID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(predictions)
} 
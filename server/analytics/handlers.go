package analytics

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers all analytics routes with the provided mux
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/analytics/spending", h.HandleSpendingAnalytics)
	mux.HandleFunc("/api/analytics/patterns", h.HandleTimePatterns)
	mux.HandleFunc("/api/analytics/predictions", h.HandlePredictions)
}

// HandleSpendingAnalytics handles GET requests for spending analytics
func (h *Handler) HandleSpendingAnalytics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handleError(w, ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	// Get account ID from URL parameters
	accountID := r.URL.Query().Get("accountId")
	if accountID == "" {
		handleError(w, ErrMissingAccountID, http.StatusBadRequest)
		return
	}

	// Get time range from query parameters (default to 1 month)
	timeRange := r.URL.Query().Get("timeRange")
	if timeRange == "" {
		timeRange = "1 month"
	}

	// Get analytics
	analytics, err := h.service.AnalyzeSpending(r.Context(), accountID, timeRange)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(analytics)
}

// HandleTimePatterns handles GET requests for time-based spending patterns
func (h *Handler) HandleTimePatterns(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handleError(w, ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	// Get account ID from URL parameters
	accountID := r.URL.Query().Get("accountId")
	if accountID == "" {
		handleError(w, ErrMissingAccountID, http.StatusBadRequest)
		return
	}

	// Get date range from query parameters (default to last month)
	endDate := time.Now()
	startDate := endDate.AddDate(0, -1, 0)

	// Get patterns
	patterns, err := h.service.GetTimePatterns(r.Context(), accountID, startDate, endDate)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(patterns)
}

// HandlePredictions handles GET requests for spending predictions
func (h *Handler) HandlePredictions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handleError(w, ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	// Get account ID from URL parameters
	accountID := r.URL.Query().Get("accountId")
	if accountID == "" {
		handleError(w, ErrMissingAccountID, http.StatusBadRequest)
		return
	}

	// Get predictions
	predictions, err := h.service.PredictSpending(r.Context(), accountID)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(predictions)
}

func handleError(w http.ResponseWriter, err error, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error": err.Error(),
	})
}

var (
	ErrMethodNotAllowed = errors.New("method not allowed")
	ErrMissingAccountID = errors.New("account ID is required")
) 
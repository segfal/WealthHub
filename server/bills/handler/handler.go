package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"server/bills/repository"
	"server/bills/service"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{service: service}
}

// SetupBillRoutes configures all the bill-related routes
func SetupBillRoutes(router *mux.Router, db *sql.DB) {
	repo := repository.NewPostgresRepository(db)
	svc := service.NewService(repo)
	handler := NewHandler(svc)
	handler.RegisterRoutes(router)
}

// RegisterRoutes registers all bill routes
func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/bills/{accountId}", h.HandleGetBills).Methods("GET")
	router.HandleFunc("/api/bills/{accountId}/recurring", h.HandleGetRecurringBills).Methods("GET")
	router.HandleFunc("/api/bills/{accountId}/upcoming", h.HandleGetUpcomingBills).Methods("GET")
	router.HandleFunc("/api/bills/{accountId}/history/{merchant}", h.HandleGetBillHistory).Methods("GET")
}

// HandleGetBills handles requests for bill payments
func (h *Handler) HandleGetBills(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountID := vars["accountId"]

	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	month, _ := strconv.Atoi(r.URL.Query().Get("month"))

	if year == 0 || month == 0 {
		now := time.Now()
		year = now.Year()
		month = int(now.Month())
	}

	bills, err := h.service.GetBillsByMonth(r.Context(), accountID, year, month)
	if err != nil {
		log.Printf("Error getting bills: %v", err)
		http.Error(w, "Failed to get bills", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bills)
}

// HandleGetRecurringBills handles requests for recurring bills
func (h *Handler) HandleGetRecurringBills(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountID := vars["accountId"]

	bills, err := h.service.GetRecurringBills(r.Context(), accountID)
	if err != nil {
		log.Printf("Error getting recurring bills: %v", err)
		http.Error(w, "Failed to get recurring bills", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bills)
}

// HandleGetUpcomingBills handles requests for upcoming bills
func (h *Handler) HandleGetUpcomingBills(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountID := vars["accountId"]

	bills, err := h.service.GetUpcomingBills(r.Context(), accountID)
	if err != nil {
		log.Printf("Error getting upcoming bills: %v", err)
		http.Error(w, "Failed to get upcoming bills", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bills)
}

// HandleGetBillHistory handles requests for bill history
func (h *Handler) HandleGetBillHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountID := vars["accountId"]
	merchant := vars["merchant"]

	history, err := h.service.GetBillHistory(r.Context(), accountID, merchant)
	if err != nil {
		log.Printf("Error getting bill history: %v", err)
		http.Error(w, "Failed to get bill history", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
} 

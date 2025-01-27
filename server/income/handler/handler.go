package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"server/income/repository"
	"server/income/service"
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

// SetupIncomeRoutes configures all the income-related routes
func SetupIncomeRoutes(router *mux.Router, db *sql.DB) {
	repo := repository.NewPostgresRepository(db)
	svc := service.NewService(repo)
	handler := NewHandler(svc)
	handler.RegisterRoutes(router)
}

// RegisterRoutes registers all income routes
func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/income/{accountId}", h.HandleGetIncome).Methods("GET")
	router.HandleFunc("/api/income/{accountId}/monthly", h.HandleGetMonthlyIncome).Methods("GET")
}

// HandleGetIncome handles requests for income data
func (h *Handler) HandleGetIncome(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountID := vars["accountId"]

	income, err := h.service.GetIncome(r.Context(), accountID)
	if err != nil {
		log.Printf("Error getting income: %v", err)
		http.Error(w, "Failed to get income", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(income)
}

// HandleGetMonthlyIncome handles requests for monthly income data
func (h *Handler) HandleGetMonthlyIncome(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountID := vars["accountId"]

	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	month, _ := strconv.Atoi(r.URL.Query().Get("month"))

	if year == 0 || month == 0 {
		now := time.Now()
		year = now.Year()
		month = int(now.Month())
	}

	income, err := h.service.GetMonthlyIncome(r.Context(), accountID, year, month)
	if err != nil {
		log.Printf("Error getting monthly income: %v", err)
		http.Error(w, "Failed to get monthly income", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(income)
} 
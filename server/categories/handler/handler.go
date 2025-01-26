package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"server/categories/repository"
	"server/categories/service"

	"github.com/gorilla/mux"
)

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{service: service}
}

// SetupCategoryRoutes configures all the category-related routes
func SetupCategoryRoutes(router *mux.Router, db *sql.DB) {
	repo := repository.NewPostgresRepository(db)
	svc := service.NewService(repo)
	handler := NewHandler(svc)
	handler.RegisterRoutes(router)
}

// RegisterRoutes registers all category routes
func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/categories/{accountId}", h.HandleGetCategories).Methods("GET")
	router.HandleFunc("/api/categories/{accountId}/totals", h.HandleGetCategoryTotals).Methods("GET")
}

// HandleGetCategories handles requests for categories
func (h *Handler) HandleGetCategories(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountID := vars["accountId"]

	categories, err := h.service.GetCategories(r.Context(), accountID)
	if err != nil {
		log.Printf("Error getting categories: %v", err)
		http.Error(w, "Failed to get categories", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

// HandleGetCategoryTotals handles requests for category totals
func (h *Handler) HandleGetCategoryTotals(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountID := vars["accountId"]

	totals, err := h.service.GetCategoryTotals(r.Context(), accountID)
	if err != nil {
		log.Printf("Error getting category totals: %v", err)
		http.Error(w, "Failed to get category totals", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(totals)
} 
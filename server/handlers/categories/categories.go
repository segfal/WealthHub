package categories

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"server/analytics"

	"github.com/gorilla/mux"
)

// SetupCategoryRoutes configures all the category-related routes for the API
func SetupCategoryRoutes(router *mux.Router, db *sql.DB) {
	repo := analytics.NewPostgresRepository(db)

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
} 
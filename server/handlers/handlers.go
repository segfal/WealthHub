package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"server/analytics"
	"server/crud"
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

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(categoryTotals)
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
}

/*
	to call the api from the client, we need to use the following url:
	http://localhost:8080/api/user/1234567891
	if you want a prediction, you need to use the following url:
	http://localhost:8080/api/predictions/1234567891

*/
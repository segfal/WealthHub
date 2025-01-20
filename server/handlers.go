package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"server/analytics"
	"time"

	"github.com/gorilla/mux"
)

func setupRoutes(router *mux.Router, db *sql.DB) {
	repo := analytics.NewPostgresRepository(db)
	analyticsService := analytics.NewService(repo)

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
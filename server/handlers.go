package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func setupRoutes(router *mux.Router, db *sql.DB) {
	router.HandleFunc("/api/analytics/{accountId}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		accountID := vars["accountId"]
		timeRange := r.URL.Query().Get("timeRange")
		if timeRange == "" {
			timeRange = "1 month" // default time range
		}

		analytics, err := analyzeSpending(db, accountID, timeRange)
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
		
		transactions, err := getTransactions(db, accountID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		categories := make(map[string]float64)
		for _, t := range transactions {
			categories[t.Category] += t.Amount
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(categories)
	}).Methods("GET")

	router.HandleFunc("/api/predictions/{accountId}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		accountID := vars["accountId"]
		
		transactions, err := getTransactions(db, accountID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		categoryTotals := make(map[string]float64)
		for _, t := range transactions {
			categoryTotals[t.Category] += t.Amount
		}

		predictions := predictFutureSpending(transactions, categoryTotals)
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(predictions)
	}).Methods("GET")

	router.HandleFunc("/api/patterns/{accountId}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		accountID := vars["accountId"]
		// to call the url http://localhost:8080/api/patterns/1234567890?timeRange=1 month
		
		transactions, err := getTransactions(db, accountID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		patterns := analyzeTimePatterns(transactions)
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(patterns)
	}).Methods("GET")
} 
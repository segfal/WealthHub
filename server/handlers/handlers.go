package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	analyticsHandler "server/analytics/handler"
	billsHandler "server/bills/handler"
	categoriesHandler "server/categories/handler"
	"server/crud"
	incomeHandler "server/income/handler"

	"github.com/gorilla/mux"
)

// SetupRoutes configures all the routes for the API
func SetupRoutes(router *mux.Router, db *sql.DB) {
	// Setup routes from each package
	analyticsHandler.SetupRoutes(router, db)
	billsHandler.SetupBillRoutes(router, db)
	categoriesHandler.SetupCategoryRoutes(router, db)
	incomeHandler.SetupIncomeRoutes(router, db)

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
}

/*
	to call the api from the client, we need to use the following url:
	http://localhost:8080/api/user/1234567891
	if you want a prediction, you need to use the following url:
	http://localhost:8080/api/predictions/1234567891

*/
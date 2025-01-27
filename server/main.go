package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	analyticsHandler "server/analytics/handler"
	analyticsRepo "server/analytics/repository"
	analyticsService "server/analytics/service"
	"server/handlers"
	"time"

	gorilla_handlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// loggingMiddleware wraps an http.Handler and logs request details
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)
		
		// Create a custom response writer to capture the status code
		rw := &responseWriter{w: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)
		
		duration := time.Since(start)
		log.Printf("Completed %s %s with status %d in %v", r.Method, r.URL.Path, rw.status, duration)
	})
}

// responseWriter wraps http.ResponseWriter to capture the status code
type responseWriter struct {
	w      http.ResponseWriter
	status int
}

func (rw *responseWriter) Header() http.Header {
	return rw.w.Header()
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	return rw.w.Write(b)
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.w.WriteHeader(statusCode)
}

func main() {
	// Set up logging with timestamp
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("Starting server...")

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Get database URL from environment
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL environment variable is required")
	}

	// Connect to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}
	log.Printf("Successfully connected to database")

	// Initialize repositories and services
	repo := analyticsRepo.NewPostgresRepository(db)
	service := analyticsService.NewService(repo)
	analyticsH := analyticsHandler.NewHandler(service)

	// Set up router with middleware
	router := mux.NewRouter()
	router.Use(loggingMiddleware)

	// Add health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// Register all routes
	handlers.SetupRoutes(router, db)
	analyticsH.RegisterRoutes(router)

	// Set up CORS
	corsMiddleware := gorilla_handlers.CORS(
		gorilla_handlers.AllowedOrigins([]string{"*"}),
		gorilla_handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		gorilla_handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		gorilla_handlers.ExposedHeaders([]string{"Content-Length"}),
		gorilla_handlers.AllowCredentials(),
		gorilla_handlers.MaxAge(3600),
	)

	// Add recovery middleware to handle panics
	recoveryHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		router.ServeHTTP(w, r)
	})

	// Create final handler chain
	finalHandler := corsMiddleware(recoveryHandler)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Printf("Server listening on port %s", port)
	if err := http.ListenAndServe(":"+port, finalHandler); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
/*


*/

# FinanceBros Backend

This is the backend server for FinanceBros, a financial analytics and management platform. The server is built using Go and follows a clean architecture pattern.

## Project Structure

```
server/
├── analytics/           # Analytics feature package
│   ├── handler/        # HTTP handlers for analytics endpoints
│   ├── service/        # Business logic for analytics
│   └── repository/     # Data access layer for analytics
├── bills/              # Bills management feature package
│   ├── handler/        # HTTP handlers for bills endpoints
│   ├── service/        # Business logic for bills
│   └── repository/     # Data access layer for bills
├── categories/         # Categories feature package
│   ├── handler/        # HTTP handlers for categories endpoints
│   ├── service/        # Business logic for categories
│   └── repository/     # Data access layer for categories
├── income/             # Income tracking feature package
│   ├── handler/        # HTTP handlers for income endpoints
│   ├── service/        # Business logic for income
│   └── repository/     # Data access layer for income
├── types/              # Shared type definitions
├── handlers/           # Main route configuration
├── crud/              # Basic CRUD operations
└── main.go            # Application entry point
```

## API Endpoints

### User Endpoints
- `GET /api/user/{accountId}`
  - Example: `http://localhost:8080/api/user/1234567891`
  - Returns user information and account details

### Analytics Endpoints
- `GET /api/analytics/{accountId}`
  - Example: `http://localhost:8080/api/analytics/1234567891`
  - Returns spending analytics for the account
- `GET /api/predictions/{accountId}`
  - Example: `http://localhost:8080/api/predictions/1234567891`
  - Returns spending predictions
- `GET /api/patterns/{accountId}`
  - Example: `http://localhost:8080/api/patterns/1234567891`
  - Returns time-based spending patterns
- `GET /api/insights/{accountId}`
  - Example: `http://localhost:8080/api/insights/1234567891`
  - Returns spending insights and recommendations

### Bills Endpoints
- `GET /api/bills/{accountId}`
  - Example: `http://localhost:8080/api/bills/1234567891`
  - Returns all bill payments
- `GET /api/bills/{accountId}/recurring`
  - Example: `http://localhost:8080/api/bills/1234567891/recurring`
  - Returns recurring bill payments
- `GET /api/bills/{accountId}/upcoming`
  - Example: `http://localhost:8080/api/bills/1234567891/upcoming`
  - Returns upcoming bill payments
- `GET /api/bills/{accountId}/history/{merchant}`
  - Example: `http://localhost:8080/api/bills/1234567891/history/Netflix`
  - Returns bill payment history for a specific merchant

### Categories Endpoints
- `GET /api/categories/{accountId}`
  - Example: `http://localhost:8080/api/categories/1234567891`
  - Returns all spending categories
- `GET /api/categories/{accountId}/totals`
  - Example: `http://localhost:8080/api/categories/1234567891/totals`
  - Returns total spending by category

### Income Endpoints
- `GET /api/income/{accountId}`
  - Example: `http://localhost:8080/api/income/1234567891`
  - Returns all income transactions
- `GET /api/income/{accountId}/monthly`
  - Example: `http://localhost:8080/api/income/1234567891/monthly?year=2024&month=3`
  - Returns monthly income data

## Setup

1. Create a `.env` file in the server directory with:
```env
DB_URL=postgres://username:password@localhost:5432/dbname
PORT=8080
```

2. Install dependencies:
```bash
go mod download
```

3. Run the server:
```bash
go run main.go
```

## Architecture

The project follows a clean architecture pattern with the following layers:

1. **Handler Layer** (`handler/`)
   - Handles HTTP requests and responses
   - Input validation and response formatting
   - Routes definition

2. **Service Layer** (`service/`)
   - Business logic implementation
   - Data processing and analysis
   - No direct database access

3. **Repository Layer** (`repository/`)
   - Data access layer
   - Database queries and operations
   - Data persistence logic

Each feature (analytics, bills, categories, income) follows this layered architecture for better organization and maintainability.

## Database Schema

The application uses PostgreSQL with the following main tables:

1. **users**
   - account_id (primary key)
   - account_name
   - account_type
   - account_number
   - balance_current
   - balance_available
   - balance_currency
   - owner_name

2. **transactions**
   - transaction_id (primary key)
   - account_id (foreign key)
   - date
   - amount
   - category
   - merchant
   - location

## Error Handling

The API uses standard HTTP status codes:
- 200: Success
- 400: Bad Request
- 404: Not Found
- 500: Internal Server Error

All endpoints return JSON responses with appropriate error messages when applicable.

## Security

- CORS is enabled for development with appropriate middleware
- Panic recovery middleware is implemented
- Request logging for debugging and monitoring
- Environment variables for sensitive configuration

## Development Guide

### Adding a New Feature

When adding a new feature, follow these steps in order for optimal development flow:

1. **Define Types** (Start here)
   ```go
   // types/prediction.go
   package types

   type PredictionRequest struct {
       AccountID string `json:"account_id"`
       TimeRange string `json:"time_range"`
   }

   type PredictionResult struct {
       Category    string    `json:"category"`
       Likelihood  float64   `json:"likelihood"`
       NextDate    time.Time `json:"next_date"`
   }
   ```

2. **Repository Layer** (Build data access first)
   ```go
   // analytics/repository/repository.go
   type Repository interface {
       // Add new method to interface
       GetPredictionData(ctx context.Context, accountID string) ([]types.Transaction, error)
   }

   // analytics/repository/postgres.go
   func (r *postgresRepo) GetPredictionData(ctx context.Context, accountID string) ([]types.Transaction, error) {
       // Implement database query
       query := `SELECT ... FROM transactions WHERE ...`
       // Execute query and return results
   }
   ```

3. **Service Layer** (Implement business logic)
   ```go
   // analytics/service/service.go
   type Service interface {
       // Add new method to interface
       GeneratePredictions(ctx context.Context, req types.PredictionRequest) ([]types.PredictionResult, error)
   }

   func (s *service) GeneratePredictions(ctx context.Context, req types.PredictionRequest) ([]types.PredictionResult, error) {
       // Get data from repository
       data, err := s.repo.GetPredictionData(ctx, req.AccountID)
       if err != nil {
           return nil, err
       }

       // Implement prediction logic
       var predictions []types.PredictionResult
       // Process data and generate predictions
       return predictions, nil
   }
   ```

4. **Handler Layer** (Add API endpoint last)
   ```go
   // analytics/handler/handler.go
   func (h *Handler) RegisterRoutes(router *mux.Router) {
       // Add new route
       router.HandleFunc("/api/predictions/{accountId}", h.HandlePredictions).Methods("GET")
   }

   func (h *Handler) HandlePredictions(w http.ResponseWriter, r *http.Request) {
       // Extract parameters
       vars := mux.Vars(r)
       accountID := vars["accountId"]

       // Create request
       req := types.PredictionRequest{
           AccountID: accountID,
           TimeRange: r.URL.Query().Get("timeRange"),
       }

       // Call service
       predictions, err := h.service.GeneratePredictions(r.Context(), req)
       if err != nil {
           http.Error(w, err.Error(), http.StatusInternalServerError)
           return
       }

       // Return response
       w.Header().Set("Content-Type", "application/json")
       json.NewEncoder(w).Encode(predictions)
   }
   ```

### Development Best Practices

1. **Order of Development**
   - Types → Repository → Service → Handler
   - This order ensures dependencies are available when needed
   - Allows for easier testing and mocking

2. **Testing Strategy**
   ```go
   // Start with repository tests
   func TestGetPredictionData(t *testing.T) {
       // Test database queries
   }

   // Then service tests
   func TestGeneratePredictions(t *testing.T) {
       // Test business logic with mocked repository
   }

   // Finally handler tests
   func TestHandlePredictions(t *testing.T) {
       // Test HTTP handling with mocked service
   }
   ```

3. **Feature Checklist**
   - [ ] Define types and interfaces
   - [ ] Implement repository methods
   - [ ] Add repository tests
   - [ ] Implement service logic
   - [ ] Add service tests
   - [ ] Create HTTP handler
   - [ ] Add handler tests
   - [ ] Update API documentation
   - [ ] Add logging and error handling

4. **Code Organization**
   ```
   feature/
   ├── types/
   │   └── feature.go         # Type definitions
   ├── repository/
   │   ├── repository.go      # Interface definition
   │   ├── postgres.go        # Implementation
   │   └── postgres_test.go   # Tests
   ├── service/
   │   ├── service.go         # Business logic
   │   └── service_test.go    # Tests
   └── handler/
       ├── handler.go         # HTTP handling
       └── handler_test.go    # Tests
   ```

5. **Error Handling Pattern**
   ```go
   // Repository layer - return specific errors
   if err != nil {
       return nil, fmt.Errorf("failed to query predictions: %w", err)
   }

   // Service layer - add context
   if err != nil {
       return nil, fmt.Errorf("failed to generate predictions: %w", err)
   }

   // Handler layer - convert to HTTP errors
   if err != nil {
       log.Printf("Error generating predictions: %v", err)
       http.Error(w, "Failed to generate predictions", http.StatusInternalServerError)
       return
   }
   ```

6. **Logging Pattern**
   ```go
   // Add context-rich logs at each layer
   log.Printf("Generating predictions for account %s with timeRange %s", 
       accountID, timeRange)

   // Log important operations and errors
   if err != nil {
       log.Printf("Failed to generate predictions: %v", err)
   }
   ```

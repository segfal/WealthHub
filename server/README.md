# FinanceBros - Financial Analytics Backend

A sophisticated financial analytics backend that analyzes spending patterns, predicts future expenses, and provides insights into financial behavior.

## Project Structure

```
server/
├── analytics/              # Analytics domain package
│   ├── handlers.go        # HTTP request handlers
│   ├── postgres.go        # PostgreSQL repository implementation
│   ├── repository.go      # Repository interface
│   └── service.go         # Business logic and analytics calculations
├── cmd/                   # Command-line applications
│   └── setup_db/         # Database setup utility
│       └── main.go       # Database initialization
├── crud/                  # Basic CRUD operations
│   └── crud.go           # Database operations
├── types/                 # Shared type definitions
│   ├── analytics.go      # Analytics-related types
│   └── transaction.go    # Transaction-related types
├── main.go               # Server entry point
└── README.md             # This documentation
```

## Analytics System

### System Design
```mermaid
graph TD
    Client[Client] --> |HTTP Request| Router[HTTP Router]
    Router --> |GET /api/analytics| AH[Analytics Handler]
    Router --> |GET /api/patterns| PH[Patterns Handler]
    Router --> |GET /api/predictions| FH[Future Predictions Handler]
    
    AH --> |Service Call| AS[Analytics Service]
    PH --> |Service Call| AS
    FH --> |Service Call| AS
    
    AS --> |Repository Call| AR[Analytics Repository]
    AR --> |SQL Query| DB[(PostgreSQL)]
    
    subgraph "Analytics Domain"
        AS --> |Time Pattern Analysis| TPA[Time Pattern Analyzer]
        AS --> |Spending Prediction| SP[Spending Predictor]
        AS --> |Category Analysis| CA[Category Analyzer]
    end
```

### Analytics Components

1. **Time Pattern Analysis** ([analytics/service.go](analytics/service.go))
   - Groups transactions by day of week and hour
   - Calculates frequency and average spend for each time slot
   - Implementation details:
     ```go
     // Key function: AnalyzeTimePatterns
     func (s *service) AnalyzeTimePatterns(ctx context.Context, accountID string, startDate, endDate time.Time) ([]types.TimePattern, error) {
         // Retrieves transactions and groups them by day and hour
         // Calculates frequency and average spend
     }
     ```
   - Mathematical model for average spend calculation:
     $$ A_{spend}(d,h) = \frac{\sum_{t \in T_{d,h}} |amount(t)|}{|T_{d,h}|} $$
     where:
     - $A_{spend}(d,h)$ is the average spend for day $d$ and hour $h$
     - $T_{d,h}$ is the set of transactions at day $d$ and hour $h$
     - $|amount(t)|$ is the absolute value of transaction amount
   - SQL Query Used:
     ```sql
     SELECT transaction_id, account_id, date, amount, category, merchant, location
     FROM transactions 
     WHERE account_id = $1 
       AND date >= NOW() - $2::INTERVAL
     ORDER BY date DESC
     ```

2. **Spending Prediction** ([analytics/service.go](analytics/service.go))
   - Uses 6-month historical data for predictions
   - Calculates likelihood scores based on frequency and amount
   - Implementation details:
     ```go
     // Key function: PredictFutureSpending
     func (s *service) PredictFutureSpending(ctx context.Context, accountID string) ([]types.PredictedSpend, error) {
         // Analyzes 6-month transaction history
         // Calculates spending patterns and likelihood scores
     }
     ```
   - Mathematical models:
     $$ frequency_{norm} = min(\frac{n_{transactions}}{180} \times 30, 1.0) $$
     $$ amount_{norm} = min(\frac{avg\_amount}{1000}, 1.0) $$
     $$ likelihood = \frac{frequency_{norm} + amount_{norm}}{2} $$
     where:
     - $n_{transactions}$ is the number of transactions in 6 months
     - $avg\_amount$ is the average transaction amount
     - Normalization factors: 180 days (6 months), 30 days (monthly), $1000 (amount threshold)
   - Prediction Algorithm Steps:
     1. Group transactions by category
     2. Calculate average time between transactions
     3. Determine spending frequency
     4. Normalize amounts and frequencies
     5. Calculate likelihood scores
     6. Generate warnings for high-likelihood events

3. **Category Analysis** ([analytics/service.go](analytics/service.go))
   - Aggregates spending by category
   - Calculates percentage distribution
   - Implementation details:
     ```go
     // Key function: GetSpendingAnalytics
     func (s *service) GetSpendingAnalytics(ctx context.Context, accountID string, timeRange string) (*types.SpendingAnalytics, error) {
         // Retrieves category totals and calculates percentages
     }
     ```
   - Mathematical model:
     $$ category\_percentage = \frac{total\_amount\_in\_category}{total\_amount\_all\_categories} \times 100 $$
   - SQL Query Used:
     ```sql
     SELECT category, COALESCE(SUM(ABS(amount)), 0) as total
     FROM transactions 
     WHERE account_id = $1 
       AND date >= NOW() - $2::INTERVAL
     GROUP BY category
     ORDER BY total DESC
     ```

### Implementation Details

1. **Repository Layer** ([analytics/repository.go](analytics/repository.go), [analytics/postgres.go](analytics/postgres.go))
   ```mermaid
   classDiagram
       class Repository {
           <<interface>>
           +GetTransactions(ctx, accountID, timeRange) []Transaction
           +GetCategoryTotals(ctx, accountID, timeRange) map[string]float64
       }
       class PostgresRepo {
           -db *sql.DB
           +GetTransactions(ctx, accountID, timeRange) []Transaction
           +GetCategoryTotals(ctx, accountID, timeRange) map[string]float64
       }
       Repository <|.. PostgresRepo
   ```
   - Interface definition:
     ```go
     type Repository interface {
         GetTransactions(ctx context.Context, accountID string, timeRange string) ([]types.Transaction, error)
         GetCategoryTotals(ctx context.Context, accountID string, timeRange string) (map[string]float64, error)
     }
     ```
   - PostgreSQL implementation uses parameterized queries to prevent SQL injection
   - Handles connection pooling and transaction management

2. **Service Layer** ([analytics/service.go](analytics/service.go))
   ```mermaid
   classDiagram
       class Service {
           <<interface>>
           +GetSpendingAnalytics(ctx, accountID, timeRange) *SpendingAnalytics
           +AnalyzeTimePatterns(ctx, accountID, startDate, endDate) []TimePattern
           +PredictFutureSpending(ctx, accountID) []PredictedSpend
       }
       class AnalyticsService {
           -repo Repository
           +GetSpendingAnalytics(ctx, accountID, timeRange) *SpendingAnalytics
           +AnalyzeTimePatterns(ctx, accountID, startDate, endDate) []TimePattern
           +PredictFutureSpending(ctx, accountID) []PredictedSpend
       }
       Service <|.. AnalyticsService
   ```
   - Business logic implementation
   - Error handling and validation
   - Data transformation and analytics calculations

3. **Handler Layer** ([analytics/handlers.go](analytics/handlers.go))
   ```mermaid
   classDiagram
       class Handler {
           -service Service
           +RegisterRoutes(mux *http.ServeMux)
           -handleAnalytics(w http.ResponseWriter, r *http.Request)
           -handlePatterns(w http.ResponseWriter, r *http.Request)
           -handlePredictions(w http.ResponseWriter, r *http.Request)
       }
   ```
   - HTTP request handling
   - Input validation
   - Response formatting
   - Error response handling

### Data Types ([types/analytics.go](types/analytics.go))

1. **SpendingAnalytics**
   ```go
   type SpendingAnalytics struct {
       TopCategories     []CategorySpend   `json:"topCategories"`
       SpendingPatterns  []TimePattern     `json:"spendingPatterns"`
       PredictedSpending []PredictedSpend  `json:"predictedSpending"`
       TotalSpent        float64           `json:"totalSpent"`
       MonthlyAverage    float64           `json:"monthlyAverage"`
   }
   ```

2. **TimePattern**
   ```go
   type TimePattern struct {
       TimeOfDay    string  `json:"timeOfDay"`
       DayOfWeek    string  `json:"dayOfWeek"`
       Frequency    int     `json:"frequency"`
       AverageSpend float64 `json:"averageSpend"`
   }
   ```

3. **PredictedSpend**
   ```go
   type PredictedSpend struct {
       Category      string    `json:"category"`
       Likelihood    float64   `json:"likelihood"`
       PredictedDate time.Time `json:"predictedDate"`
       Warning       string    `json:"warning,omitempty"`
   }
   ```

### API Endpoints and Examples

1. **GET /api/analytics/{accountId}**
   - File: [analytics/handlers.go](analytics/handlers.go)
   - Handler: `handleAnalytics`
   - Example Request:
     ```bash
     curl "http://localhost:8080/api/analytics/1234567891?timeRange=1%20month"
     ```
   - Example Response:
     ```json
     {
       "topCategories": [
         {
           "category": "Groceries",
           "totalSpent": "543.21",
           "percentage": "32.5"
         }
       ],
       "spendingPatterns": [
         {
           "timeOfDay": "18:00",
           "dayOfWeek": "Friday",
           "frequency": 12,
           "averageSpend": 45.67
         }
       ],
       "predictedSpending": [
         {
           "category": "Groceries",
           "likelihood": 0.85,
           "predictedDate": "2024-02-01T18:00:00Z",
           "warning": "High likelihood (85%) of spending in Groceries category around Feb 01"
         }
       ],
       "totalSpent": 1672.43,
       "monthlyAverage": 1672.43
     }
     ```

2. **GET /api/patterns/{accountId}**
   - File: [analytics/handlers.go](analytics/handlers.go)
   - Handler: `handlePatterns`
   - Example Request:
     ```bash
     curl "http://localhost:8080/api/patterns/1234567891?start=2024-01-01&end=2024-02-01"
     ```
   - Example Response:
     ```json
     [
       {
         "timeOfDay": "12:00",
         "dayOfWeek": "Monday",
         "frequency": 8,
         "averageSpend": 25.50
       }
     ]
     ```

3. **GET /api/predictions/{accountId}**
   - File: [analytics/handlers.go](analytics/handlers.go)
   - Handler: `handlePredictions`
   - Example Request:
     ```bash
     curl "http://localhost:8080/api/predictions/1234567891"
     ```
   - Example Response:
     ```json
     [
       {
         "category": "Dining",
         "likelihood": 0.75,
         "predictedDate": "2024-02-03T19:00:00Z",
         "warning": "High likelihood (75%) of spending in Dining category around Feb 03"
       }
     ]
     ```

### Error Handling

The API uses standard HTTP status codes with detailed error messages:

1. **400 Bad Request**
   ```json
   {
     "error": "Invalid account ID format"
   }
   ```

2. **404 Not Found**
   ```json
   {
     "error": "Account not found: 1234567891"
   }
   ```

3. **500 Internal Server Error**
   ```json
   {
     "error": "Failed to query database: connection refused"
   }
   ```

Implementation ([analytics/handlers.go](analytics/handlers.go)):
```go
func handleError(w http.ResponseWriter, err error, status int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(map[string]string{
        "error": err.Error(),
    })
}
```

### Security Considerations

1. **Environment Variables** ([.env.example](.env.example))
   ```env
   DB_URL=postgresql://postgres:password@host:port/dbname
   DB_NAME=dbname
   DB_USER=username
   DB_PASSWORD=password
   DB_HOST=host
   DB_PORT=port
   PORT=8080
   ```

2. **SQL Injection Prevention** ([analytics/postgres.go](analytics/postgres.go))
   ```go
   // Using parameterized queries
   rows, err := r.db.QueryContext(ctx, query, accountID, timeRange)
   ```

3. **Input Validation** ([analytics/handlers.go](analytics/handlers.go))
   ```go
   // Validate account ID format
   if !validateAccountID(accountID) {
       handleError(w, errors.New("invalid account ID format"), http.StatusBadRequest)
       return
   }
   ```

4. **CORS Configuration** ([main.go](main.go))
   ```go
   // CORS middleware setup
   w.Header().Set("Access-Control-Allow-Origin", "*")
   w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
   w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
   ```

### Performance Optimization

1. **Database Indexing**
   ```sql
   CREATE INDEX idx_transactions_account_date ON transactions(account_id, date);
   CREATE INDEX idx_transactions_category ON transactions(category);
   ```

2. **Connection Pooling**
   ```go
   // Database configuration
   db.SetMaxOpenConns(25)
   db.SetMaxIdleConns(25)
   db.SetConnMaxLifetime(5 * time.Minute)
   ```

3. **Query Optimization**
   - Using appropriate indexes
   - Limiting result sets
   - Efficient joins and aggregations

4. **Caching Strategy**
   - In-memory caching for frequently accessed data
   - Cache invalidation on data updates
   - Configurable cache duration

## Testing

### Unit Tests
Location: [analytics/service_test.go](analytics/service_test.go)
```go
func TestAnalyzeTimePatterns(t *testing.T) {
    // Test cases for time pattern analysis
}

func TestPredictFutureSpending(t *testing.T) {
    // Test cases for spending predictions
}
```

### Integration Tests
Location: [analytics/integration_test.go](analytics/integration_test.go)
```go
func TestAnalyticsEndToEnd(t *testing.T) {
    // End-to-end test cases
}
```

### Load Tests
Location: [tests/load/main.go](tests/load/main.go)
```go
func BenchmarkAnalyticsEndpoint(b *testing.B) {
    // Load test implementation
}
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines.

## License

MIT License - see [LICENSE](LICENSE) file for details 
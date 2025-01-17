# FinanceBros - Financial Analytics Backend

A sophisticated financial analytics backend that analyzes spending patterns, predicts future expenses, and provides insights into financial behavior.

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

1. **Time Pattern Analysis**
   - Groups transactions by day of week and hour
   - Calculates frequency and average spend for each time slot
   - Mathematical model for average spend calculation:
     ```latex
     A_{spend}(d,h) = \frac{\sum_{t \in T_{d,h}} |amount(t)|}{|T_{d,h}|}
     ```
     where:
     - $A_{spend}(d,h)$ is the average spend for day $d$ and hour $h$
     - $T_{d,h}$ is the set of transactions at day $d$ and hour $h$
     - $|amount(t)|$ is the absolute value of transaction amount

2. **Spending Prediction**
   - Uses 6-month historical data for predictions
   - Calculates likelihood scores based on frequency and amount
   - Mathematical models:
     ```latex
     frequency_{norm} = min(\frac{n_{transactions}}{180} \times 30, 1.0)
     amount_{norm} = min(\frac{avg\_amount}{1000}, 1.0)
     likelihood = \frac{frequency_{norm} + amount_{norm}}{2}
     ```
     where:
     - $n_{transactions}$ is the number of transactions in 6 months
     - $avg\_amount$ is the average transaction amount
     - Normalization factors: 180 days (6 months), 30 days (monthly), $1000 (amount threshold)

3. **Category Analysis**
   - Aggregates spending by category
   - Calculates percentage distribution
   - Mathematical model:
     ```latex
     category\_percentage = \frac{total\_amount\_in\_category}{total\_amount\_all\_categories} \times 100
     ```

### API Endpoints

1. **GET /api/analytics/{accountId}**
   ```typescript
   type SpendingAnalytics = {
     topCategories: Array<{
       category: string;
       totalSpent: string;
       percentage: string;
     }>;
     spendingPatterns: Array<{
       timeOfDay: string;
       dayOfWeek: string;
       frequency: number;
       averageSpend: number;
     }>;
     predictedSpending: Array<{
       category: string;
       likelihood: number;
       predictedDate: string;
       warning?: string;
     }>;
     totalSpent: number;
     monthlyAverage: number;
   }
   ```
   Query Parameters:
   - `timeRange`: string ("1 month" | "3 months" | "6 months" | "1 year")

2. **GET /api/patterns/{accountId}**
   ```typescript
   type TimePattern = {
     timeOfDay: string;
     dayOfWeek: string;
     frequency: number;
     averageSpend: number;
   }
   ```
   Query Parameters:
   - `start`: string (YYYY-MM-DD)
   - `end`: string (YYYY-MM-DD)

3. **GET /api/predictions/{accountId}**
   ```typescript
   type PredictedSpend = {
     category: string;
     likelihood: number;
     predictedDate: string;
     warning?: string;
   }
   ```

### Implementation Details

1. **Repository Layer**
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

2. **Service Layer**
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

3. **Handler Layer**
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

### Example Usage

1. Get comprehensive analytics:
   ```bash
   curl "http://localhost:8080/api/analytics/1234567891?timeRange=1%20month"
   ```

2. Get time patterns for a specific period:
   ```bash
   curl "http://localhost:8080/api/patterns/1234567891?start=2024-01-01&end=2024-02-01"
   ```

3. Get spending predictions:
   ```bash
   curl "http://localhost:8080/api/predictions/1234567891"
   ```

### Error Handling

The API uses standard HTTP status codes:
- 200: Success
- 400: Bad Request (invalid account ID, invalid time range)
- 500: Internal Server Error (database errors, calculation errors)

All error responses include a message:
```json
{
  "error": "Detailed error message"
}
```

## Setup and Installation

### Prerequisites

1. Go 1.19 or higher
2. PostgreSQL 12 or higher
3. Git

### Required Go Packages

```bash
go get -u github.com/gorilla/mux
go get -u github.com/lib/pq
go get -u github.com/joho/godotenv
```

### Environment Setup

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd financebros
   ```

2. Create `.env` file with database credentials:
   ```env
   DB_URL=postgresql://postgres:password@host:port/dbname
   DB_NAME=dbname
   DB_USER=username
   DB_PASSWORD=password
   DB_HOST=host
   DB_PORT=port
   PORT=8080
   ```

3. Initialize the database:
   ```bash
   psql -U postgres -d your_database -f init.sql
   ```

### Running the Server

1. Start the server:
   ```bash
   go run .
   ```

2. The server will start on port 8080 (or the port specified in .env)

### API Usage Examples

1. Get Analytics:
   ```bash
   curl "http://localhost:8080/api/analytics/123?timeRange=1%20month"
   ```

2. Get Categories:
   ```bash
   curl "http://localhost:8080/api/categories/123"
   ```

3. Get Predictions:
   ```bash
   curl "http://localhost:8080/api/predictions/123"
   ```

4. Get Patterns:
   ```bash
   curl "http://localhost:8080/api/patterns/123"
   ```

## Response Formats

### Analytics Response
```json
{
  "topCategories": [
    {
      "category": "string",
      "totalSpent": 0.0,
      "percentage": 0.0
    }
  ],
  "spendingPatterns": [
    {
      "timeOfDay": "15:00",
      "dayOfWeek": "Monday",
      "frequency": 0,
      "averageSpend": 0.0
    }
  ],
  "predictedSpending": [
    {
      "category": "string",
      "likelihood": 0.0,
      "predictedDate": "2024-01-05T15:04:05Z",
      "warning": "string"
    }
  ],
  "totalSpent": 0.0,
  "monthlyAverage": 0.0
}
```

## Error Handling

The API uses standard HTTP status codes:
- 200: Success
- 400: Bad Request
- 404: Not Found
- 500: Internal Server Error

All error responses include a message in the response body:
```json
{
  "error": "Error message description"
}
```

## Security Considerations

1. All database credentials are stored in environment variables
2. CORS is enabled for cross-origin requests
3. Input validation is performed on all API endpoints
4. SQL injection prevention through parameterized queries

## Performance Optimization

1. Database queries are optimized with indexes
2. Connection pooling is implemented
3. Results are cached where appropriate
4. Batch processing for large datasets

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

MIT License - see LICENSE file for details 
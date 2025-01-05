# FinanceBros - Financial Analytics Backend

A sophisticated financial analytics backend that analyzes spending patterns, predicts future expenses, and provides insights into financial behavior.

## System Design

### Architecture Overview
```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│   Client    │ ──► │  REST API    │ ──► │  Database   │
│  (Frontend) │ ◄── │  (Go Server) │ ◄── │ (PostgreSQL)│
└─────────────┘     └──────────────┘     └─────────────┘
```

### Components

1. **Database Layer (PostgreSQL)**
   - Stores transaction data
   - Maintains account information
   - Schema:
     ```sql
     accounts (
         account_id VARCHAR(20) PRIMARY KEY,
         account_name VARCHAR(50),
         account_type VARCHAR(20),
         account_number VARCHAR(10),
         balance_current DECIMAL(10, 2),
         balance_available DECIMAL(10, 2),
         currency VARCHAR(10),
         owner_name VARCHAR(50),
         bank_name VARCHAR(50),
         routing_number VARCHAR(10),
         branch VARCHAR(100)
     )

     transactions (
         transaction_id VARCHAR(20) PRIMARY KEY,
         account_id VARCHAR(20),
         date DATE,
         amount DECIMAL(10, 2),
         category VARCHAR(50),
         merchant VARCHAR(50),
         location VARCHAR(100)
     )
     ```

2. **API Layer**
   - RESTful endpoints for data access
   - Analytics processing
   - Prediction generation

3. **Analytics Engine**
   - Spending pattern analysis
   - Predictive modeling
   - Category-based analysis
   - Time-based pattern recognition

### API Endpoints

1. **GET /api/analytics/{accountId}**
   - Comprehensive spending analysis
   - Query params: `timeRange` (e.g., "1 month", "3 months", "1 year")
   - Returns: Full analytics including top categories, patterns, and predictions

2. **GET /api/categories/{accountId}**
   - Category-wise spending breakdown
   - Returns: Map of categories to total spent

3. **GET /api/predictions/{accountId}**
   - Future spending predictions
   - Returns: Array of predicted spending events with likelihoods

4. **GET /api/patterns/{accountId}**
   - Temporal spending patterns
   - Returns: Analysis of spending patterns by time and day

## Mathematical Models

### Spending Prediction Model

The likelihood of future spending is calculated using a normalized frequency and amount model:

$$ L = \frac{min(f_{n}, 1) + min(a_{n}, 1)}{2} $$

Where:
- $L$ is the likelihood score (0 to 1)
- $f_{n}$ is the normalized frequency: $f_{n} = \frac{frequency}{30}$
- $a_{n}$ is the normalized amount: $a_{n} = \frac{average\_amount}{1000}$

### Pattern Recognition

Time patterns are analyzed using frequency and average spend calculations:

$$ A_{spend} = \frac{\sum_{i=1}^{n} amount_i}{n} $$

Where:
- $A_{spend}$ is the average spend for a time period
- $n$ is the number of transactions in that period

### Prediction Intervals

Future spending dates are predicted using average time between transactions:

$$ T_{avg} = \frac{\sum_{i=1}^{n-1} (t_{i+1} - t_i)}{n-1} $$

Where:
- $T_{avg}$ is the average time between transactions
- $t_i$ is the timestamp of transaction i

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
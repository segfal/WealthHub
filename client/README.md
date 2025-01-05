Clone the repository

npm install

npm run dev


# Install Dependencies

npm install

# Run the application

npm run dev



## 🧩 Components

### Homepage [src/Homepage.tsx]
Main dashboard layout featuring:
- Dark theme with Robinhood-inspired design
- User welcome section with balance display
- Quick stats cards for financial overview
- Tabbed navigation system
- Framer Motion animations

### SpendingOverview [src/components/SpendingOverview.tsx]
Displays key financial metrics:
- Total spent amount
- Average transaction value
- Transaction count
- Interactive spending trend charts
- Real-time data fetching from `/api/analytics/{accountId}`

### SpendingCategories [src/components/SpendingCategories.tsx]
Visualizes spending patterns:
- Interactive pie chart for category distribution
- Progress bars for spending breakdown
- Percentage calculations
- Color-coded categories
- Data from `/api/categories/{accountId}`

### SpendingPredictions [src/components/SpendingPredictions.tsx]
Provides spending forecasts:
- Line charts for trend visualization
- Confidence indicators
- Warning system for unusual patterns
- Category-wise predictions
- Data from `/api/predictions/{accountId}`

### SpendingPatterns [src/components/SpendingPatterns.tsx]
Analyzes spending behavior:
- Time-of-day analysis
- Day-of-week breakdown
- Recurring transaction detection
- Pattern visualization
- Data from `/api/patterns/{accountId}`
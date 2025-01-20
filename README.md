# FinanceBros - Financial Analytics Platform

A modern financial analytics platform built with Go and React, providing spending analysis, pattern recognition, and predictive analytics.

## Prerequisites

- Node.js (v16 or higher)
- Go (v1.21 or higher)
- PostgreSQL (v14 or higher)
- Git

## Project Structure

```
FinanceBros/
├── client/           # React frontend
├── server/           # Go backend
│   ├── analytics/    # Analytics services
│   ├── crud/         # Database operations
│   └── types/        # Shared types
```

## Quick Start

1. Clone the repository:
```bash
git clone https://github.com/yourusername/FinanceBros.git
cd FinanceBros
```

2. Install Dependencies:

Feeling lazy? We've got you covered! Just run our setup script:
```bash
chmod +x setup.sh  # Make it executable first
./setup.sh         # Let the magic happen
```
Or if you're feeling extra lazy, do it all in one go:
```bash
bash setup.sh
```

The script will check if you have all the required tools (npm, node, and Go) and guide you through the installation process.

3. Set up the database:
```bash
# Create a PostgreSQL database named 'financebros'
createdb financebros
```

4. Configure environment variables:
```bash
# In server/.env
DB_URL=postgresql://postgres:password@localhost:5432/financebros
PORT=8080
```

5. Install dependencies:
```bash
# Install frontend dependencies
cd client
npm install

# Install backend dependencies
cd ../server
go mod download
```

6. Start the development servers:
```bash
# From the root directory
npm run dev
```

This will start both the frontend and backend servers concurrently:
- Frontend: http://localhost:5173
- Backend: http://localhost:8080

## Development

### Frontend (React + TypeScript)

The frontend is built with:
- React
- TypeScript
- Tailwind CSS
- Recharts for data visualization
- Framer Motion for animations

To run only the frontend:
```bash
cd client
npm run dev
```

### Backend (Go)

The backend features:
- RESTful API
- PostgreSQL database
- Hot reloading with Air
- CORS support
- Comprehensive logging

To run only the backend:
```bash
cd server
./run.sh
```

## API Endpoints

- `GET /api/analytics/spending?accountId={id}&timeRange={range}`
  - Get spending analytics and category breakdown

- `GET /api/analytics/patterns?accountId={id}`
  - Get spending patterns and time-based analysis

- `GET /api/analytics/predictions?accountId={id}`
  - Get spending predictions and trends

## Features

- 📊 Spending Analytics
  - Category breakdown
  - Monthly averages
  - Spending trends

- 🕒 Time Pattern Analysis
  - Day of week patterns
  - Time of day analysis
  - Recurring transaction detection

- 📈 Predictive Analytics
  - Future spending predictions
  - Category-wise forecasting
  - Confidence scores

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 
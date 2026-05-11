# Finance Tracker API

A REST API for tracking finances and managing budgets in Go.

## Overview

Finance Tracker is a backend service that allows users to:
- Set and manage monthly budgets
- Track expenses
- View expense totals and summaries
- Create accounts and authenticate securely - currently in development in feat/auth branch
## Project Structure

```
finance-tracker/
├── cmd/
│   └── finance-api/         # Application entry point
│       └── main.go
├── internal/
│   ├── auth/                # Authentication logic - feat/auth branch
│   │   └── auth.go
│   ├── config/              # Configuration management
│   │   └── env.go
│   ├── data/                # Data models
│   │   ├── budget.go
│   │   ├── expenses.go
│   │   └── expenses_test.go
│   ├── database/            # Database operations
│   │   ├── auth.go
│   │   └── database.go
│   └── handlers/            # HTTP request handlers
│       ├── auth.go
│       ├── budget.go
│       ├── budget_test.go
│       └── expenses.go
├── migrations/              # Database schema migrations
│   └── pg_querys.sql
├── scripts/                 # Utility scripts
│   └── db-migrate.go
├── go.mod
└── README.md
```

## Installation

### Prerequisites

- Go 1.26.2 or higher
- PostgreSQL database
- Git

### Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/jjma22/finance-tracker.git
   cd finance-tracker
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   
   Create a `.env` file in the project root with the following variables:
   ```env
   DB_host=localhost
   DB_port=5432
   DB_user=postgres
   DB_password=your_password
   DB_name=finance_tracker_database
   ```

4. **Run database migrations**
   Run migrations/pg_querys.sql on database

5. **Build and run**
   ```bash
   go run cmd/finance-api/main.go
   ```

   The API will be available at `http://127.0.0.1:9090`

## API Endpoints

### Authentication - in development, see feat/auth branch

- `POST /login` - Authenticate user and receive session token
- `POST /create/user` - Create a new user account

### Monthly Budget

- `POST /monthlybudget` - Set monthly budget (budget validated with middleware)
- `GET /monthlybudget/{id}` - Get budget details
- `PUT /monthlybudget/{id}` - Update monthly budget

### Expenses

- `POST /expense` - Add new expense (expense validated with middleware)
- `GET /expense` - Get all expenses
- `GET /expense/{id}` - Get specific expense
- `GET /expense/total` - Get total expenses
- `PUT /expense/update/{id}` - Update expense
- `DELETE /expense/delete/{id}` - Delete expense

## CI
CI is configured via .github/workflows/go.yaml
For each pull request into main, a containerised db will be spun up and initilaised using scripts/db-migrate.go. Tests will be ran using this database.
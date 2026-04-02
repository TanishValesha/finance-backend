# Finance Data Processing and Access Control Backend

A RESTful backend API for a finance dashboard system built with **Go**, **Gin**, **PostgreSQL**, and **GORM**. The system supports user role management, financial record keeping, dashboard analytics, and role-based access control.

---

## Tech Stack

| Layer            | Technology            |
| ---------------- | --------------------- |
| Language         | Go                    |
| Framework        | Gin                   |
| Database         | PostgreSQL (Supabase) |
| ORM              | GORM                  |
| Authentication   | JWT                   |
| Password Hashing | bcrypt                |

---

## Project Structure

```
finance-backend/
├── cmd/
│   └── seed/
│       └── main.go          # Database seeder
├── config/
│   └── db.go                # Database connection
├── handlers/
│   ├── auth.go              # Register, Login
│   ├── user.go              # User management
│   ├── transaction.go       # Transaction CRUD
│   └── dashboard.go         # Dashboard summary APIs
├── middleware/
│   ├── auth.go              # JWT validation + active status check
│   └── rbac.go              # Role-based access control
├── models/
│   ├── user.go              # User model
│   └── transaction.go       # Transaction model
├── routes/
│   └── routes.go            # All route definitions
├── services/
│   ├── auth_service.go      # JWT generation, password hashing
│   └── dashboard_service.go # Aggregation logic
├── utils/
│   └── response.go          # Consistent JSON response helpers
├── .env.example
├── go.mod
└── README.md
```

---

## Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL database (local or Supabase)

### Setup

**1. Clone the repository**

```bash
git clone https://github.com/TanishValesha/finance-backend
cd finance-backend
```

**2. Install dependencies**

```bash
go mod tidy
```

**3. Configure environment**

```bash
cp .env.example .env
```

Edit `.env` with your database credentials:

```env
DB_HOST=your_db_host
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=your_db_name
JWT_SECRET=your_super_secret_key
PORT=8080
```

**4. Run the server**

```bash
go run main.go
```

**5. Seed the database** (optional but recommended for quick testing)

```bash
cd cmd/seed
go run main.go
```

This creates the following test accounts:

| Role    | Email               | Password   |
| ------- | ------------------- | ---------- |
| Admin   | admin@finance.com   | admin123   |
| Analyst | analyst@finance.com | analyst123 |
| Viewer  | viewer@finance.com  | viewer123  |

---

## API Reference

### Authentication

> All protected routes require `Authorization: Bearer <token>` header.

| Method | Endpoint             | Access | Description                                |
| ------ | -------------------- | ------ | ------------------------------------------ |
| POST   | `/api/auth/register` | Public | Register a new user (default role: viewer) |
| POST   | `/api/auth/login`    | Public | Login and receive JWT token                |

---

### User Management

| Method | Endpoint                | Access    | Description                   |
| ------ | ----------------------- | --------- | ----------------------------- |
| GET    | `/api/me`               | All roles | Get current logged in user    |
| GET    | `/api/users`            | Admin     | Get all users                 |
| PATCH  | `/api/users/:id/role`   | Admin     | Update a user's role          |
| PATCH  | `/api/users/:id/status` | Admin     | Activate or deactivate a user |

---

### Transactions

| Method | Endpoint                | Access    | Description                                |
| ------ | ----------------------- | --------- | ------------------------------------------ |
| GET    | `/api/transactions`     | All roles | Get transactions (own only for non-admins) |
| GET    | `/api/transactions/:id` | All roles | Get single transaction                     |
| POST   | `/api/transactions/`    | Admin     | Create a transaction                       |
| PUT    | `/api/transactions/:id` | Admin     | Update a transaction                       |
| DELETE | `/api/transactions/:id` | Admin     | Soft delete a transaction                  |

#### Filtering & Pagination

```
GET /api/transactions?type=income
GET /api/transactions?category=food
GET /api/transactions?from=2026-01-01&to=2026-03-31
GET /api/transactions?type=expense&category=food&page=2&limit=5
```

| Query Param | Type   | Description                                                                                     |
| ----------- | ------ | ----------------------------------------------------------------------------------------------- |
| `type`      | string | `income` or `expense`                                                                           |
| `category`  | string | `salary`, `freelance`, `food`, `transport`, `utilities`, `entertainment`, `healthcare`, `other` |
| `from`      | string | Start date `YYYY-MM-DD`                                                                         |
| `to`        | string | End date `YYYY-MM-DD`                                                                           |
| `page`      | int    | Page number (default: 1)                                                                        |
| `limit`     | int    | Records per page (default: 10, max: 100)                                                        |

---

### Dashboard

| Method | Endpoint                         | Access         | Description                         |
| ------ | -------------------------------- | -------------- | ----------------------------------- |
| GET    | `/api/dashboard/summary`         | Analyst, Admin | Total income, expenses, net balance |
| GET    | `/api/dashboard/category-totals` | Analyst, Admin | Spending grouped by category        |
| GET    | `/api/dashboard/trends`          | Analyst, Admin | Monthly income vs expense           |
| GET    | `/api/dashboard/recent`          | Analyst, Admin | Last 5 transactions                 |

---

## Role Permissions

| Action                 | Viewer | Analyst | Admin |
| ---------------------- | ------ | ------- | ----- |
| Login / Register       | ✓      | ✓       | ✓     |
| View own transactions  | ✓      | ✓       | ✓     |
| View all transactions  | ✗      | ✗       | ✓     |
| View dashboard summary | ✗      | ✓       | ✓     |
| Create transactions    | ✗      | ✗       | ✓     |
| Update transactions    | ✗      | ✗       | ✓     |
| Delete transactions    | ✗      | ✗       | ✓     |
| Manage users           | ✗      | ✗       | ✓     |

---

## Design Decisions & Assumptions

**Authentication**
JWT is used for stateless authentication. Tokens expire after 24 hours. Every protected request validates the token AND checks the user's active status in the database — so deactivated users are blocked immediately even if their token is still valid.

**Soft Delete**
Transactions are never hard deleted. GORM's `DeletedAt` field is used to mark records as deleted while keeping them in the database for audit purposes. All queries automatically exclude soft-deleted records.

**Ownership**
Non-admin users (viewer, analyst) can only view transactions they created. Admins can view and manage all transactions.

**Role Enforcement**
Access control is applied at the middleware level using route groups — not inside individual handlers. This keeps handlers clean and makes permissions easy to audit in one place (`routes/routes.go`).

**Admin Self-Protection**
An admin cannot deactivate or change the status of their own account to prevent accidental lockout.

**Dashboard Aggregation**
All dashboard calculations are done directly in PostgreSQL using raw SQL aggregation queries (`SUM`, `GROUP BY`, `DATE_TRUNC`) rather than in application memory. This is more efficient and scales better with large datasets.

**Seed Script**
A seed script is provided at `cmd/seed/main.go` that creates test users and 18 sample transactions spread across 4 months. It is safe to run multiple times — existing records are skipped.

---

## Sample Request & Response

**Login**

```bash
POST /api/auth/login
Content-Type: application/json

{
  "email": "admin@finance.com",
  "password": "admin123"
}
```

```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": 1,
      "name": "Admin User",
      "email": "admin@finance.com",
      "role": "admin"
    }
  }
}
```

**Dashboard Summary**

```bash
GET /api/dashboard/summary
Authorization: Bearer
```

```json
{
  "success": true,
  "message": "Summary fetched",
  "data": {
    "total_income": 370000,
    "total_expenses": 25200,
    "net_balance": 344800
  }
}
```

---

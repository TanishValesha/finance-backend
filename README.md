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

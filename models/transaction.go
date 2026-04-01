package models

import (
	"time"

	"gorm.io/gorm"
)

type TransactionType string
type TransactionCategory string

const (
	TypeIncome  TransactionType = "income"
	TypeExpense TransactionType = "expense"
)

const (
	CategorySalary        TransactionCategory = "salary"
	CategoryFreelance     TransactionCategory = "freelance"
	CategoryFood          TransactionCategory = "food"
	CategoryTransport     TransactionCategory = "transport"
	CategoryUtilities     TransactionCategory = "utilities"
	CategoryEntertainment TransactionCategory = "entertainment"
	CategoryHealthcare    TransactionCategory = "healthcare"
	CategoryOther         TransactionCategory = "other"
)

type Transaction struct {
	ID          uint                `gorm:"primaryKey;autoIncrement" json:"id"`
	Amount      float64             `gorm:"not null" json:"amount"`
	Type        TransactionType     `gorm:"not null" json:"type"`
	Category    TransactionCategory `gorm:"not null" json:"category"`
	Date        time.Time           `gorm:"not null" json:"date"`
	Notes       string              `json:"notes"`
	CreatedByID uint                `gorm:"not null" json:"created_by"`
	CreatedBy   User                `gorm:"foreignKey:CreatedByID" json:"-"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	DeletedAt   gorm.DeletedAt      `gorm:"index" json:"-"`
}

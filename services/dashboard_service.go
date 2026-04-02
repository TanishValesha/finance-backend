package services

import (
	"finance-backend/config"
	"time"
)

type SummaryResult struct {
	TotalIncome   float64 `json:"total_income"`
	TotalExpenses float64 `json:"total_expenses"`
	NetBalance    float64 `json:"net_balance"`
}

type CategoryTotal struct {
	Category string  `json:"category"`
	Total    float64 `json:"total"`
}

type RecentTransaction struct {
	ID       uint      `json:"id"`
	Amount   float64   `json:"amount"`
	Type     string    `json:"type"`
	Category string    `json:"category"`
	Date     time.Time `json:"date"`
	Notes    string    `json:"notes"`
}

type MonthlyTrend struct {
	Month   string  `json:"month"`
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
}

func GetSummary() (SummaryResult, error) {
	var result SummaryResult

	// Calculate total income
	if err := config.DB.Raw(`SELECT COALESCE(SUM(amount), 0) FROM transactions WHERE type = 'income' AND deleted_at IS NULL`).Scan(&result.TotalIncome).Error; err != nil {
		return result, err
	}

	// Calculate total expenses
	if err := config.DB.Raw(`SELECT COALESCE(SUM(amount), 0) FROM transactions WHERE type = 'expense' AND deleted_at IS NULL`).Scan(&result.TotalExpenses).Error; err != nil {
		return result, err
	}

	result.NetBalance = result.TotalIncome - result.TotalExpenses

	return result, nil
}

func GetCategoryTotals() ([]CategoryTotal, error) {
	var results []CategoryTotal

	if err := config.DB.Raw(`SELECT category, COALESCE(SUM(amount), 0) AS total FROM transactions WHERE deleted_at IS NULL GROUP BY category ORDER BY total DESC`).Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func GetMonthlyTrends() ([]MonthlyTrend, error) {
	var results []MonthlyTrend

	if err := config.DB.Raw(`
		SELECT
			TO_CHAR(DATE_TRUNC('month', date), 'YYYY-MM') as month,
			COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) as income,
			COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) as expense
		FROM transactions
		WHERE deleted_at IS NULL
		GROUP BY DATE_TRUNC('month', date)
		ORDER BY DATE_TRUNC('month', date) ASC
	`).Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func GetRecentTransactions(limit int) ([]RecentTransaction, error) {
	var results []RecentTransaction

	if err := config.DB.Raw(`SELECT id, amount, type, category, date, notes FROM transactions WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT ?`, limit).Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil

}

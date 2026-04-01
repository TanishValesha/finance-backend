package services

import "finance-backend/config"

type SummaryResult struct {
	TotalIncome   float64 `json:"total_income"`
	TotalExpenses float64 `json:"total_expenses"`
	NetBalance    float64 `json:"net_balance"`
}

type CategoryTotal struct {
	Category string  `json:"category"`
	Total    float64 `json:"total"`
}

func GetSummary() (SummaryResult, error) {
	var result SummaryResult

	// Calculate total income
	config.DB.Raw(`SELECT COALESCE(SUM(amount), 0) FROM transactions WHERE type = 'income' AND deleted_at IS NULL`).Scan(&result.TotalIncome)

	// Calculate total expenses
	config.DB.Raw(`SELECT COALESCE(SUM(amount), 0) FROM transactions WHERE type = 'expense' AND deleted_at IS NULL`).Scan(&result.TotalExpenses)

	result.NetBalance = result.TotalIncome - result.TotalExpenses

	return result, nil
}

func GetCategoryTotals() ([]CategoryTotal, error) {
	var results []CategoryTotal

	config.DB.Raw(`SELECT category, COALESCE(SUM(amount), 0) AS total FROM transactions WHERE deleted_at IS NULL GROUP BY category ORDER BY total DESC`).Scan(&results)

	return results, nil
}

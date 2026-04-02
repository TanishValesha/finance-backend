package handlers

import (
	"finance-backend/services"
	"finance-backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetSummary(c *gin.Context) {
	summary, err := services.GetSummary()
	if err != nil {
		utils.Error(c, 500, "Failed to fetch summary")
		return
	}

	utils.Success(c, http.StatusOK, "Summary fetched", summary)
}

func GetCategoryBreakdown(c *gin.Context) {
	breakdown, err := services.GetCategoryTotals()
	if err != nil {
		utils.Error(c, 500, "Failed to fetch category breakdown")
		return
	}

	utils.Success(c, http.StatusOK, "Category breakdown fetched", breakdown)
}

func GetMonthlyTrends(c *gin.Context) {
	trends, err := services.GetMonthlyTrends()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch monthly trends")
		return
	}

	utils.Success(c, http.StatusOK, "Monthly trends fetched", trends)
}

func GetRecentTransactions(c *gin.Context) {
	transactions, err := services.GetRecentTransactions(4)
	if err != nil {
		utils.Error(c, 500, "Failed to fetch recent transactions")
		return
	}

	utils.Success(c, http.StatusOK, "Recent transactions fetched", transactions)
}

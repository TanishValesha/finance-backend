package handlers

import (
	"finance-backend/config"
	"finance-backend/models"
	"finance-backend/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateTransactionInput struct {
	Amount   float64 `json:"amount" binding:"required,gt=0"`
	Type     string  `json:"type" binding:"required"`
	Category string  `json:"category" binding:"required"`
	Date     string  `json:"date" binding:"required"` // format: "2006-01-02"
	Notes    string  `json:"notes"`
}

type UpdateTransactionInput struct {
	Amount   float64 `json:"amount" binding:"omitempty,gt=0"`
	Type     string  `json:"type"`
	Category string  `json:"category"`
	Date     string  `json:"date"`
	Notes    string  `json:"notes"`
}

func validateTransactionCategory(category string) bool {
	validCategories := []string{
		"salary", "freelance", "food", "transport", "utilities",
		"entertainment", "healthcare", "other",
	}

	for _, c := range validCategories {
		if category == c {
			return true
		}
	}
	return false
}

func validateTransactionType(t string) bool {
	return t == "income" || t == "expense"
}

func CreateTransaction(c *gin.Context) {
	var input CreateTransactionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if !validateTransactionType(input.Type) {
		utils.Error(c, http.StatusBadRequest, "Type must be income or expense")
		return
	}

	if !validateTransactionCategory(input.Category) {
		utils.Error(c, http.StatusBadRequest, "Invalid category")
		return
	}

	parsedDate, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid date format. Use YYYY-MM-DD")
		return
	}

	userID, _ := c.Get("userID")

	transaction := models.Transaction{
		Amount:      input.Amount,
		Type:        models.TransactionType(input.Type),
		Category:    models.TransactionCategory(input.Category),
		Date:        parsedDate,
		Notes:       input.Notes,
		CreatedByID: userID.(uint),
	}

	if err := config.DB.Create(&transaction).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to create transaction")
		return
	}

	utils.Success(c, http.StatusCreated, "Transaction created", transaction)

}

func GetTransactions(c *gin.Context) {
	var transactions []models.Transaction

	typeOfTransaction := c.Query("type")
	category := c.Query("category")

	query := config.DB.Model(&models.Transaction{})

	// Optional Filters

	if typeOfTransaction != "" {
		if !validateTransactionType(typeOfTransaction) {
			utils.Error(c, http.StatusBadRequest, "Type must be income or expense")
			return
		}
		query = query.Where("type = ?", typeOfTransaction)
	}

	if category != "" {
		if !validateTransactionCategory(category) {
			utils.Error(c, http.StatusBadRequest, "Invalid category")
			return
		}
		query = query.Where("category = ?", category)
	}

	if err := query.Order("date desc").Find(&transactions).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch transactions")
		return
	}

	utils.Success(c, http.StatusOK, "Transactions fetched", transactions)
}

func GetTransactionByID(c *gin.Context) {
	id := c.Param("id")
	var transaction models.Transaction

	if err := config.DB.First(&transaction, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Transaction not found")
		return
	}

	utils.Success(c, http.StatusOK, "Transaction fetched", transaction)
}

func UpdateTransaction(c *gin.Context) {
	id := c.Param("id")

	var transaction models.Transaction
	if err := config.DB.First(&transaction, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Transaction not found")
		return
	}

	var input UpdateTransactionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	updates := map[string]interface{}{}

	if input.Amount > 0 {
		updates["amount"] = input.Amount
	}

	if input.Type != "" {
		if !validateTransactionType(input.Type) {
			utils.Error(c, http.StatusBadRequest, "Type must be income or expense")
			return
		}
		updates["type"] = input.Type
	}

	if input.Category != "" {
		if !validateTransactionCategory(input.Category) {
			utils.Error(c, http.StatusBadRequest, "Invalid category")
			return
		}
		updates["category"] = input.Category
	}

	if input.Notes != "" {
		updates["notes"] = input.Notes
	}

	if input.Date != "" {
		parsedDate, err := time.Parse("2006-01-02", input.Date)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, "Invalid date format. Use YYYY-MM-DD")
			return
		}
		updates["date"] = parsedDate
	}

	if err := config.DB.Model(&transaction).Updates(updates); err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to update transaction")
		return
	}
	utils.Success(c, http.StatusOK, "Transaction updated", transaction)

}

func DeleteTransaction(c *gin.Context) {
	id := c.Param("id")

	var transaction models.Transaction
	if err := config.DB.First(&transaction, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Transaction not found")
		return
	}

	if err := config.DB.Delete(&transaction).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to delete transaction")
		return
	}

	utils.Success(c, http.StatusOK, "Transaction deleted", nil)
}

package handlers

import (
	"finance-backend/config"
	"finance-backend/models"
	"finance-backend/utils"
	"net/http"
	"strconv"
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
	userID, _ := c.Get("userID")
	userRole, _ := c.Get("userRole")
	var transactions []models.Transaction

	typeOfTransaction := c.Query("type")
	category := c.Query("category")
	from := c.Query("from")
	to := c.Query("to")
	p := c.Query("page")
	l := c.Query("limit")

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

	if from != "" {
		parsedFrom, err := time.Parse("2006-01-02", from)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, "Invalid from date format. Use YYYY-MM-DD")
			return
		}
		query = query.Where("date >= ?", parsedFrom)
	}

	if to != "" {
		parsedTo, err := time.Parse("2006-01-02", to)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, "Invalid to date format. Use YYYY-MM-DD")
			return
		}
		query = query.Where("date <= ?", parsedTo)
	}

	// Pagination

	page := 1
	limit := 10

	if p != "" {
		parsedPage, err := strconv.Atoi(p)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, "Invalid page format")
			return
		}
		if parsedPage < 1 {
			utils.Error(c, http.StatusBadRequest, "Page must be greater than 0")
			return
		}
		page = parsedPage
	}

	if l != "" {
		parsedLimit, err := strconv.Atoi(l)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, "Invalid limit format")
			return
		}
		if parsedLimit < 1 {
			utils.Error(c, http.StatusBadRequest, "Limit must be greater than 0")
			return
		}
		limit = parsedLimit
	}

	offset := (page - 1) * limit

	var total int64
	if err := query.Count(&total).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to count transactions")
		return
	}

	if userRole.(string) != "admin" {
		query = query.Where("created_by_id = ?", userID)
	}

	if err := query.Order("date desc").Offset(offset).Limit(limit).Find(&transactions).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch transactions")
		return
	}

	utils.Success(c, http.StatusOK, "Transactions fetched", gin.H{
		"transactions": transactions,
		"total":        total,
		"page":         page,
		"limit":        limit,
		"total_pages":  (int(total) + limit - 1) / limit,
	})
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

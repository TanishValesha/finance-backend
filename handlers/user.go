package handlers

import (
	"finance-backend/config"
	"finance-backend/models"
	"finance-backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetMe(c *gin.Context) {
	userID, _ := c.Get("userID")

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "User not found")
		return
	}

	utils.Success(c, http.StatusOK, "User fetched", gin.H{
		"id":        user.ID,
		"name":      user.Name,
		"email":     user.Email,
		"role":      user.Role,
		"is_active": user.IsActive,
	})

}

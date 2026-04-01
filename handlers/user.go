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

func GetAllUsers(c *gin.Context) {
	var users []models.User
	if err := config.DB.Find(&users).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	utils.Success(c, http.StatusOK, "Users fetched", users)
}

type UpdateRoleInput struct {
	Role string `json:"role" binding:"required"`
}

func UpdateUserRole(c *gin.Context) {
	id := c.Param("id")
	var input UpdateRoleInput

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if input.Role != "viewer" && input.Role != "analyst" && input.Role != "admin" {
		utils.Error(c, http.StatusBadRequest, "Invalid role. Must be viewer, analyst, or admin")
		return
	}

	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "User not found")
		return
	}

	config.DB.Model(&user).Update("role", input.Role)
	utils.Success(c, http.StatusOK, "Role updated successfully", nil)
}

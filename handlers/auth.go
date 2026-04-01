package handlers

import (
	"finance-backend/config"
	"finance-backend/models"
	"finance-backend/services"
	"finance-backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, http.StatusConflict, err.Error())
		return
	}

	var existingUser models.User

	if err := config.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		utils.Error(c, http.StatusConflict, "Email already registered")
		return
	}

	hashedPassword, err := services.HashPassword(input.Password)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	user := models.User{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: hashedPassword,
		Role:         models.RoleViewer, //deafult role
	}

	if err := config.DB.Create(&user).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to create user")
		return
	}

	utils.Success(c, http.StatusCreated, "User registered successfully", gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
	})
}

func Login(c *gin.Context) {
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, http.StatusConflict, err.Error())
		return
	}

	var user models.User

	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		utils.Error(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	if !user.IsActive {
		utils.Error(c, http.StatusForbidden, "Account is inactive")
		return
	}

	if !services.CheckPassword(input.Password, user.PasswordHash) {
		utils.Error(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	token, err := services.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	utils.Success(c, http.StatusOK, "Login successful", gin.H{
		"token": token,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})

}

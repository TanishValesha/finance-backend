package routes

import (
	"finance-backend/handlers"
	"finance-backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	auth := r.Group("/api/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
	}

	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		api.GET("/me", handlers.GetMe)

		api.GET("/transactions", handlers.GetTransactions)
		api.GET("/transactions/:id", handlers.GetTransactionByID)

		adminTransactions := api.Group("/transactions")
		adminTransactions.Use(middleware.RequiredRoles("admin"))
		{
			adminTransactions.POST("/", handlers.CreateTransaction)
			adminTransactions.PUT("/transactions/:id", handlers.UpdateTransaction)
			adminTransactions.DELETE("/transactions/:id", handlers.DeleteTransaction)
		}

		dashboard := api.Group("/dashboard")
		dashboard.Use(middleware.RequiredRoles("admin", "analyst"))
		{
			dashboard.GET("/summary", handlers.GetSummary)
			dashboard.GET("/category-breakdown", handlers.GetCategoryBreakdown)
		}

		users := api.Group("/users")
		users.Use(middleware.RequiredRoles("admin"))
		{
			users.GET("/", handlers.GetAllUsers)
			users.PATCH("/:id/role", handlers.UpdateUserRole)
		}
	}

	return r
}

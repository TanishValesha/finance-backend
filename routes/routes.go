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

		users := api.Group("/users")
		users.Use(middleware.RequiredRoles("admin"))
		{
			users.GET("/", handlers.GetAllUsers)
			users.PATCH("/:id/role", handlers.UpdateUserRole)
		}
	}

	return r
}

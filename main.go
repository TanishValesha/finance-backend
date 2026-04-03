package main

import (
	"finance-backend/config"
	"finance-backend/routes"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	config.ConnectDB()

	// config.DB.AutoMigrate(&models.User{}, &models.Transaction{})

	r := routes.SetupRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)

	r.Run(":" + port)

}

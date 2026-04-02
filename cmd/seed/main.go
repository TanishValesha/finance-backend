package main

import (
	"finance-backend/config"
	"finance-backend/models"
	"finance-backend/services"
	"log"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect DB
	config.ConnectDB()

	// Auto migrate
	config.DB.AutoMigrate(&models.User{}, &models.Transaction{})

	log.Println("Seeding database...")

	seedUsers()
	seedTransactions()

	log.Println("Seeding complete!")
}

func seedUsers() {
	users := []struct {
		Name     string
		Email    string
		Password string
		Role     models.Role
	}{
		{"Admin User", "admin@finance.com", "admin123", models.RoleAdmin},
		{"Analyst User", "analyst@finance.com", "analyst123", models.RoleAnalyst},
		{"Viewer User", "viewer@finance.com", "viewer123", models.RoleViewer},
	}

	for _, u := range users {
		// Skip if already exists
		var existing models.User
		if err := config.DB.Where("email = ?", u.Email).First(&existing).Error; err == nil {
			log.Printf("User %s already exists, skipping\n", u.Email)
			continue
		}

		hashedPassword, err := services.HashPassword(u.Password)
		if err != nil {
			log.Printf("Failed to hash password for %s: %v\n", u.Email, err)
			continue
		}

		user := models.User{
			Name:         u.Name,
			Email:        u.Email,
			PasswordHash: hashedPassword,
			Role:         u.Role,
			IsActive:     true,
		}

		if err := config.DB.Create(&user).Error; err != nil {
			log.Printf("Failed to create user %s: %v\n", u.Email, err)
			continue
		}

		log.Printf("Created user: %s (%s)\n", u.Email, u.Role)
	}

}

func seedTransactions() {
	var admin models.User
	if err := config.DB.Where("email = ?", "admin@finance.com").First(&admin).Error; err != nil {
		log.Println("Admin user not found, skipping transaction seeding")
		return
	}

	var count int64
	config.DB.Model(&models.Transaction{}).Count(&count)
	if count > 0 {
		log.Println("Transactions already exist, skipping")
		return
	}

	now := time.Now()

	transactions := []models.Transaction{
		// January
		{Amount: 85000, Type: models.TypeIncome, Category: models.CategorySalary, Date: date(now, -3, 1), Notes: "January salary", CreatedByID: admin.ID},
		{Amount: 1200, Type: models.TypeExpense, Category: models.CategoryFood, Date: date(now, -3, 5), Notes: "Groceries", CreatedByID: admin.ID},
		{Amount: 3500, Type: models.TypeExpense, Category: models.CategoryUtilities, Date: date(now, -3, 10), Notes: "Electricity and water", CreatedByID: admin.ID},
		{Amount: 800, Type: models.TypeExpense, Category: models.CategoryTransport, Date: date(now, -3, 15), Notes: "Fuel", CreatedByID: admin.ID},
		{Amount: 12000, Type: models.TypeIncome, Category: models.CategoryFreelance, Date: date(now, -3, 20), Notes: "Freelance project", CreatedByID: admin.ID},

		// February
		{Amount: 85000, Type: models.TypeIncome, Category: models.CategorySalary, Date: date(now, -2, 1), Notes: "February salary", CreatedByID: admin.ID},
		{Amount: 1500, Type: models.TypeExpense, Category: models.CategoryFood, Date: date(now, -2, 7), Notes: "Groceries + dining out", CreatedByID: admin.ID},
		{Amount: 2500, Type: models.TypeExpense, Category: models.CategoryEntertainment, Date: date(now, -2, 14), Notes: "Weekend trip", CreatedByID: admin.ID},
		{Amount: 4000, Type: models.TypeExpense, Category: models.CategoryHealthcare, Date: date(now, -2, 18), Notes: "Medical checkup", CreatedByID: admin.ID},
		{Amount: 900, Type: models.TypeExpense, Category: models.CategoryTransport, Date: date(now, -2, 22), Notes: "Cab rides", CreatedByID: admin.ID},

		// March
		{Amount: 85000, Type: models.TypeIncome, Category: models.CategorySalary, Date: date(now, -1, 1), Notes: "March salary", CreatedByID: admin.ID},
		{Amount: 18000, Type: models.TypeIncome, Category: models.CategoryFreelance, Date: date(now, -1, 10), Notes: "Freelance project", CreatedByID: admin.ID},
		{Amount: 1100, Type: models.TypeExpense, Category: models.CategoryFood, Date: date(now, -1, 12), Notes: "Groceries", CreatedByID: admin.ID},
		{Amount: 3200, Type: models.TypeExpense, Category: models.CategoryUtilities, Date: date(now, -1, 16), Notes: "Bills", CreatedByID: admin.ID},
		{Amount: 5500, Type: models.TypeExpense, Category: models.CategoryOther, Date: date(now, -1, 25), Notes: "Miscellaneous", CreatedByID: admin.ID},

		// Current month
		{Amount: 85000, Type: models.TypeIncome, Category: models.CategorySalary, Date: date(now, 0, 1), Notes: "This month salary", CreatedByID: admin.ID},
		{Amount: 1300, Type: models.TypeExpense, Category: models.CategoryFood, Date: date(now, 0, 3), Notes: "Groceries", CreatedByID: admin.ID},
		{Amount: 700, Type: models.TypeExpense, Category: models.CategoryTransport, Date: date(now, 0, 6), Notes: "Fuel", CreatedByID: admin.ID},
	}

	for _, t := range transactions {
		if err := config.DB.Create(&t).Error; err != nil {
			log.Printf("Failed to create transaction: %v\n", err)
			continue
		}
	}

	log.Printf("Created %d transactions\n", len(transactions))

}

func date(base time.Time, monthOffset int, day int) time.Time {
	return time.Date(base.Year(), base.Month()+time.Month(monthOffset), day, 0, 0, 0, 0, time.UTC)
}

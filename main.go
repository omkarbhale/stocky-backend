package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"

	"stockybackend/src/controllers"
	"stockybackend/src/database"
	"stockybackend/src/middlewares"
	"stockybackend/src/models"
	"stockybackend/src/routes"
)

func main() {
	fmt.Println("Started...")

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}

	database.Connect()
	models.SeedDatabase(database.DB, false)

	r := gin.Default()

	// Add the logger middleware
	r.Use(middlewares.LoggerMiddleware())

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello world",
		})
	})

	routes.RegisterUserRoutes(r)
	routes.RegisterRewardRoutes(r)
	routes.RegisterSymbolRoutes(r)
	routes.RegisterPortfolioRoutes(r)
	routes.RegisterBonusRoutes(r)

	// // Start the price update thread
	// controllers.TestPriceSimulation()
	controllers.GeneratePast12HoursPrices(database.DB)
	go startPriceUpdater(database.DB)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}

func startPriceUpdater(db *gorm.DB) {
	for {
		now := time.Now()
		nextHour := now.Truncate(time.Hour).Add(time.Hour)
		time.Sleep(time.Until(nextHour))
		controllers.UpdateSymbolPrices(db)
	}
}

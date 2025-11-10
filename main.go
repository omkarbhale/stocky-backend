package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"stockybackend/src/controllers"
	"stockybackend/src/database"
	"stockybackend/src/models"
	"stockybackend/src/routes"
)

func main() {
	fmt.Println("Started...")

	database.Connect()
	models.SeedDatabase(database.DB, false)

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello world",
		})
	})

	routes.RegisterUserRoutes(r)
	routes.RegisterRewardRoutes(r)
	routes.RegisterSymbolRoutes(r)

	// // Start the price update thread
	// controllers.TestPriceSimulation()
	controllers.GeneratePast12HoursPrices(database.DB)
	go startPriceUpdater(database.DB)

	r.Run(":8080") // TODO Use dotenv for PORT
}

func startPriceUpdater(db *gorm.DB) {
	for {
		now := time.Now()
		nextHour := now.Truncate(time.Hour).Add(time.Hour)
		time.Sleep(time.Until(nextHour))
		controllers.UpdateSymbolPrices(db)
	}
}

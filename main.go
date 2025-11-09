package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

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

	r.Run(":8080") // TODO Use dotenv for PORT
}

package main

import (
	"fmt"

	"stockybackend/src/database"
	"stockybackend/src/models"
)

func main() {
	fmt.Println("Started...")

	database.Connect()
	models.SeedDatabase(database.DB, false)
}

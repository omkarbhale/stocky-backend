package controllers

import (
	"math/rand"
	"net/http"
	"time"

	"gorm.io/gorm"

	"stockybackend/src/database"
	"stockybackend/src/models"

	"github.com/gin-gonic/gin"
)

const (
	priceChangeMin = -0.008 // -0.8%
	priceChangeMax = 0.008  // +0.8%
	priceDrift     = 0.0001 // Small positive drift
)

func GetSymbols(c *gin.Context) {
	symbols := []struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}{}
	database.DB.Model(&models.Symbol{}).Select("id, name").Scan(&symbols)
	c.JSON(http.StatusOK, symbols)
}

func UpdateSymbolPrices(db *gorm.DB) {
	var symbols []models.Symbol
	db.Find(&symbols)

	for _, symbol := range symbols {
		var latestPrice models.SymbolPriceHistory
		db.Where("symbol_id = ?", symbol.ID).Order("date DESC, time_hour DESC").First(&latestPrice)

		newPrice := 100.0 // Default starting price
		if latestPrice.ID != 0 {
			change := priceChangeMin + rand.Float64()*(priceChangeMax-priceChangeMin) + priceDrift
			newPrice = latestPrice.Price * (1 + change)
		}

		newPriceEntry := models.SymbolPriceHistory{
			SymbolID: symbol.ID,
			Price:    newPrice,
			TimeHour: uint(time.Now().Hour()),
			Date:     time.Now(),
		}
		db.Create(&newPriceEntry)
	}
}

// func SimulatePriceSeries(initialPrice float64, steps int) []float64 {
// 	prices := make([]float64, steps)
// 	prices[0] = initialPrice

// 	for i := 1; i < steps; i++ {
// 		change := priceChangeMin + rand.Float64()*(priceChangeMax-priceChangeMin) + priceDrift
// 		prices[i] = prices[i-1] * (1 + change)
// 	}

// 	return prices
// }

// func TestPriceSimulation() {
// 	initialPrice := 10000.0
// 	steps := 200 // Simulate for 24 hours
// 	prices := SimulatePriceSeries(initialPrice, steps)

// 	fmt.Println("Simulated Prices:")
// 	for i, price := range prices {
// 		fmt.Printf("Hour %d: %.2f\n", i, price)
// 	}
// }

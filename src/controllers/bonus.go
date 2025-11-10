package controllers

import (
	"fmt"
	"net/http"
	"stockybackend/src/database"
	"stockybackend/src/models"

	"github.com/gin-gonic/gin"
)

func GetHoldingsPerStockSymbol(c *gin.Context) {
	// Step 1: parse userId param
	userIdParam := c.Param("userId")
	var userId uint
	if _, err := fmt.Sscanf(userIdParam, "%d", &userId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user ID"})
		return
	}

	// Step 2: ensure user exists
	var user models.User
	if err := database.DB.First(&user, userId).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("User with ID %d not found", userId)})
		return
	}

	// Step 3: Struct for response
	type PortfolioItem struct {
		SymbolID   uint    `json:"symbolId"`
		SymbolName string  `json:"symbolName"`
		Quantity   float64 `json:"quantity"`
		Price      float64 `json:"price"`
		ValueINR   float64 `json:"valueInr"`
	}

	var portfolio []PortfolioItem
	totalValue := 0.0

	// Step 4: Main query
	// Explanation:
	// - Sum all rewards (holdings) per symbol for this user
	// - Join with latest available price per symbol using DISTINCT ON (Postgres)
	// - Multiply quantity * price for current value
	query := `
	SELECT
	    r.symbol_id AS symbol_id,
	    s.name AS symbol_name,
	    SUM(r.quantity) AS quantity,
	    p.price AS price,
	    SUM(r.quantity) * p.price AS value_inr
	FROM rewards r
	JOIN symbols s ON s.id = r.symbol_id
	JOIN (
	    SELECT DISTINCT ON (symbol_id)
	        symbol_id,
	        price
	    FROM symbol_price_histories
	    ORDER BY symbol_id, date DESC, time_hour DESC
	) p ON p.symbol_id = r.symbol_id
	WHERE r.user_id = ?
	GROUP BY r.symbol_id, s.name, p.price
	ORDER BY s.name;
	`

	if err := database.DB.Raw(query, userId).Scan(&portfolio).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Database error (portfolio query)",
			"error":   err.Error(),
		})
		return
	}

	// Step 5: Compute total portfolio value
	for _, item := range portfolio {
		totalValue += item.ValueINR
	}

	// Step 6: Return response
	c.JSON(http.StatusOK, gin.H{
		"userId":        userId,
		"portfolio":     portfolio,
		"totalValueINR": totalValue,
		"symbolCount":   len(portfolio),
	})
}

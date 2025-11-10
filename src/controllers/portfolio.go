package controllers

import (
	"fmt"
	"net/http"
	"stockybackend/src/database"
	"stockybackend/src/models"
	"time"

	"github.com/gin-gonic/gin"
)

func GetInrValueTillYesterday(c *gin.Context) {
	userIdParam := c.Param("userId")
	var userId uint
	if _, err := fmt.Sscanf(userIdParam, "%d", &userId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user ID"})
		return
	}

	// ensure user exists
	var user models.User
	if err := database.DB.First(&user, userId).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("User with ID %d not found", userId)})
		return
	}

	// Compute yesterday’s date
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1).Format("2006-01-02")

	type DailyValue struct {
		Date     string  `json:"date"`
		TotalINR float64 `json:"totalInr"`
	}
	var results []DailyValue

	var allDates []time.Time
	database.DB.Model(&models.SymbolPriceHistory{}).
		Select("DISTINCT date").
		Where("date <= ?", yesterday).
		Order("date asc").
		Pluck("date", &allDates)

	for _, d := range allDates {
		dateStr := d.Format("2006-01-02")

		type row struct {
			TotalINR float64
		}
		var result row

		err := database.DB.Raw(`
			SELECT SUM(h.total_quantity * p.price) AS total_inr
			FROM (
			    SELECT symbol_id, SUM(quantity) AS total_quantity
			    FROM rewards
			    WHERE user_id = ?
			      AND timestamp <= ?
			    GROUP BY symbol_id
			) AS h
			JOIN symbol_price_histories p
			  ON p.symbol_id = h.symbol_id
			WHERE p.date = ?
			  AND p.time_hour = 23
		`, userId, dateStr+" 23:59:59", dateStr).Scan(&result).Error

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "DB error", "error": err.Error()})
			return
		}

		results = append(results, DailyValue{
			Date:     dateStr,
			TotalINR: result.TotalINR,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"userId":  userId,
		"history": results,
	})
}

func GetTodaysUserStats(c *gin.Context) {
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

	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())

	type TodayReward struct {
		SymbolID   uint    `json:"symbolId"`
		SymbolName string  `json:"symbolName"`
		Quantity   float64 `json:"quantity"`
	}

	// result1: get today's total shares rewarded (grouped by symbol)
	var todayRewards []TodayReward
	if err := database.DB.Raw(`
		SELECT r.symbol_id AS symbol_id,
		       s.name AS symbol_name,
		       SUM(r.quantity) AS quantity
		FROM rewards r
		JOIN symbols s ON s.id = r.symbol_id
		WHERE r.user_id = ?
		  AND r.timestamp BETWEEN ? AND ?
		GROUP BY r.symbol_id, s.name
		ORDER BY s.name;
	`, userId, startOfDay, endOfDay).Scan(&todayRewards).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "DB error (todayRewards)", "error": err.Error()})
		return
	}

	type HoldingValue struct {
		SymbolID uint
		Price    float64
		Quantity float64
	}

	// result2: calculate current INR value of entire portfolio
	var holdings []HoldingValue
	if err := database.DB.Raw(`
		SELECT r.symbol_id AS symbol_id,
		       SUM(r.quantity) AS quantity,
		       p.price AS price
		FROM rewards r
		JOIN (
		    SELECT DISTINCT ON (symbol_id) symbol_id, price
		    FROM symbol_price_histories
		    ORDER BY symbol_id, date DESC, time_hour DESC
		) p ON p.symbol_id = r.symbol_id
		WHERE r.user_id = ?
		GROUP BY r.symbol_id, p.price;
	`, userId).Scan(&holdings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "DB error (holdings)", "error": err.Error()})
		return
	}

	totalValue := 0.0
	for _, h := range holdings {
		totalValue += h.Quantity * h.Price
	}

	c.JSON(http.StatusOK, gin.H{
		"userId":          userId,
		"todayRewards":    todayRewards,
		"currentValueINR": totalValue,
	})
}

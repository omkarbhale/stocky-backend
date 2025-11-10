package controllers

import (
	"fmt"
	"net/http"
	"stockybackend/src/database"
	"stockybackend/src/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateReward(c *gin.Context) {
	var body struct {
		UserId    *uint      `json:"userId"`
		SymbolID  *uint      `json:"symbolId"`
		Quantity  *float64   `json:"quantity"`
		Timestamp *time.Time `json:"time"`
	}

	// --- Basic JSON Validations ---
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON: " + err.Error()})
		return
	}
	if body.UserId == nil || body.SymbolID == nil || body.Quantity == nil || body.Timestamp == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing required fields"})
		return
	}

	// Check userId and symbolId in db
	var user models.User
	if err := database.DB.First(&user, *body.UserId).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("User with ID %d not found", *body.UserId)})
		return
	}

	var symbol models.Symbol
	if err := database.DB.First(&symbol, *body.SymbolID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Symbol with ID %d not found", *body.SymbolID)})
		return
	}

	// Assume 4% transaction fees
	// Get latest price for that symbol to calculate 4% transaction fees
	var latestPrice struct {
		Price float64
	}
	err := database.DB.Raw(`
		SELECT price FROM symbol_price_histories
		WHERE symbol_id = ?
		ORDER BY date DESC, time_hour DESC
		LIMIT 1
	`, *body.SymbolID).Scan(&latestPrice).Error

	if err != nil || latestPrice.Price == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Could not fetch latest price for symbol"})
		return
	}

	// Save reward
	reward := models.Reward{
		UserID:    *body.UserId,
		SymbolID:  *body.SymbolID,
		Quantity:  *body.Quantity,
		Timestamp: *body.Timestamp,
	}
	if err := database.DB.Create(&reward).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Could not save reward: " + err.Error()})
		return
	}

	// Compute ledger values
	baseValue := latestPrice.Price * (*body.Quantity)
	fees := baseValue * 0.04
	totalCost := baseValue + fees

	// Save transaction
	transaction := models.Transaction{
		Description: fmt.Sprintf("Purchase of %.2f shares of %s for user %s", *body.Quantity, symbol.Name, user.Name),
	}
	if err := database.DB.Create(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create ledger transaction", "error": err.Error()})
		return
	}

	var cashAcc, stockAcc, feeAcc models.Account
	database.DB.First(&cashAcc, "name = ?", "Cash")
	database.DB.First(&stockAcc, "name = ?", "StockInvestments")
	database.DB.First(&feeAcc, "name = ?", "TransactionFees")

	// Save entries
	entries := []models.Entry{
		{TransactionID: transaction.ID, AccountID: stockAcc.ID, Type: "debit", Amount: baseValue},
		{TransactionID: transaction.ID, AccountID: feeAcc.ID, Type: "debit", Amount: fees},
		{TransactionID: transaction.ID, AccountID: cashAcc.ID, Type: "credit", Amount: totalCost},
	}

	for _, e := range entries {
		if err := database.DB.Create(&e).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create ledger entry", "error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User rewarded & ledger updated",
		"user":    user,
		"symbol":  symbol,
		"reward":  reward,
		"ledger": gin.H{
			"transaction": transaction.Description,
			"entries": []gin.H{
				{"account": stockAcc.Name, "type": "debit", "amount": baseValue},
				{"account": feeAcc.Name, "type": "debit", "amount": fees},
				{"account": cashAcc.Name, "type": "credit", "amount": totalCost},
			},
		},
	})
}

func GetRewardsForUser(c *gin.Context) {
	userIdParam := c.Param("userId")
	userIdUint64, err := strconv.ParseUint(userIdParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user ID"})
		return
	}

	var body struct {
		Date *time.Time `json:"date"`
	}

	_ = c.ShouldBindJSON(&body)

	var user models.User
	if err := database.DB.First(&user, userIdUint64).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("User with ID %d not found", userIdUint64)})
		return
	}

	var dateToMatch time.Time
	if body.Date == nil {
		dateToMatch = time.Now()
	} else {
		dateToMatch = *body.Date
	}

	type RewardResponse struct {
		UserID    uint       `json:"user_id"`
		SymbolID  uint       `json:"symbol_id"`
		Quantity  float64    `json:"quantity"`
		Timestamp *time.Time `json:"timestamp"`
	}

	var rewards []RewardResponse
	err2 := database.DB.Model(&models.Reward{}).Select("user_id", "symbol_id", "quantity", "timestamp").Where("USER_ID=? AND DATE(timestamp) = DATE(?)", userIdUint64, dateToMatch).Find(&rewards).Error
	if err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not get rewards from database"})
		return
	}

	c.JSON(http.StatusOK, rewards)
}

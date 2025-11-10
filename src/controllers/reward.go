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

	// Basic json validations
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON: " + err.Error()})
		return
	}
	if body.UserId == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Field userId missing"})
		return
	}
	if body.SymbolID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Field symbolId missing"})
		return
	}
	if body.Quantity == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Field quantity missing"})
		return
	}
	if body.Timestamp == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Field time missing"})
		return
	}

	// Validations that given IDs are present in db
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

	// Save reward to DB
	reward := models.Reward{
		UserID:    *body.UserId,
		SymbolID:  *body.SymbolID,
		Quantity:  *body.Quantity,
		Timestamp: *body.Timestamp,
	}
	result := database.DB.Create(&reward)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Coult not save reward to database: " + result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User rewarded", "user": user, "symbol": symbol, "reward": reward})
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

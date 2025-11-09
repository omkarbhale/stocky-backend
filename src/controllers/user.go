package controllers

import (
	"net/http"
	"stockybackend/src/database"
	"stockybackend/src/models"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	var body struct {
		Name string `json:"name"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON"})
		return
	}
	if body.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Name is missing"})
		return
	}

	user := models.User{Name: body.Name}
	result := database.DB.Create(&user) // pass pointer!

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created", "user": user})
}

func GetUsers(c *gin.Context) {
	users := []models.User{}
	database.DB.Model(models.User{}).Find(&users)
	c.JSON(http.StatusOK, users)
}

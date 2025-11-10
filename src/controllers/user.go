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

	database.DB.Create(models.User{
		Name: body.Name,
	})

	c.JSON(http.StatusCreated, gin.H{"message": "User created"})
}

func GetUsers(c *gin.Context) {
	users := []models.User{}
	database.DB.Model(models.User{}).Find(&users)
	c.JSON(http.StatusOK, users)
}

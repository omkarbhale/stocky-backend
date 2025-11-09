package controllers

import (
	"net/http"
	"stockybackend/src/database"
	"stockybackend/src/models"

	"github.com/gin-gonic/gin"
)

func GetSymbols(c *gin.Context) {
	symbols := []struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}{}
	database.DB.Model(&models.Symbol{}).Select("id, name").Scan(&symbols)
	c.JSON(http.StatusOK, symbols)
}

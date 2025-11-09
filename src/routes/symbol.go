package routes

import (
	"stockybackend/src/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterSymbolRoutes(r *gin.Engine) {
	symbolRoutes := r.Group("/symbol")
	{
		symbolRoutes.GET("/", controllers.GetSymbols)
	}
}

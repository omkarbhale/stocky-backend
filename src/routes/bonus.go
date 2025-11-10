package routes

import (
	"stockybackend/src/controllers"

	"github.com/gin-gonic/gin"
)

// *(Bonus: Add `/portfolio/{userId}` to show holdings per stock symbol with current INR value.)*

func RegisterBonusRoutes(r *gin.Engine) {
	r.GET("/portfolio/:userId", controllers.GetHoldingsPerStockSymbol)
}

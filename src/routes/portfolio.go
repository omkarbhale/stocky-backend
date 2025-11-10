package routes

import (
	"stockybackend/src/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterPortfolioRoutes(r *gin.Engine) {
	r.GET("/historical-inr/:userId", controllers.GetInrValueTillYesterday)
	r.GET("/stats/:userId", controllers.GetTodaysUserStats)
}

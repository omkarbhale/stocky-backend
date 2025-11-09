package routes

import (
	"stockybackend/src/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRewardRoutes(c *gin.Engine) {
	rewardRoutes := c.Group("/reward")
	{
		rewardRoutes.POST("/", controllers.CreateReward)
	}
}

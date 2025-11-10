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

	// Outside group because this route "/today-stocks" was part of the assignment
	c.GET("/today-stocks/:userId", controllers.GetRewardsForUser)
}

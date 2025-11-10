package routes

import (
	"stockybackend/src/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRewardRoutes(r *gin.Engine) {
	rewardRoutes := r.Group("/reward")
	{
		rewardRoutes.POST("/", controllers.CreateReward)
	}

	// Outside group because this route "/today-stocks" was part of the assignment
	r.GET("/today-stocks/:userId", controllers.GetRewardsForUser)
}

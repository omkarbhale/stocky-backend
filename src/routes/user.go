package routes

import (
	"stockybackend/src/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine) {
	r.POST("/user", controllers.CreateUser)
	r.GET("/user", controllers.GetUsers)
}

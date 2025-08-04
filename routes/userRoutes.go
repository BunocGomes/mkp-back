package routes

import (
	"github.com/BunocGomes/mkp-back/controller"
	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(router *gin.Engine) {
	userGroup := router.Group("/api/v1/users")
	{
		userGroup.POST("/", controller.CreateUser)
		userGroup.GET("/", controller.GetAllUsers)
		userGroup.GET("/:id", controller.GetUser)
		userGroup.PUT("/:id", controller.UpdateUser)
		userGroup.DELETE("/:id", controller.DeleteUser)

		userGroup.POST("/login", controller.Login)
	}

}

package routes

import (
	controller "github.com/BunocGomes/mkp-back/controller"
	"github.com/BunocGomes/mkp-back/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	userGroup := router.Group("/api/v1/users")
	userGroup.Use(middleware.AuthMiddleware())
	{
		userGroup.POST("/", controller.CreateUser)
		userGroup.GET("/", controller.GetAllUsers)
		userGroup.GET("/:id", controller.GetUser)
		userGroup.PUT("/:id", controller.UpdateUser)
		userGroup.DELETE("/:id", controller.DeleteUser)
	}

	loginUserGroup := router.Group("/api/v1/auth")
	{
		loginUserGroup.POST("/login", controller.Login)
	}

	createUserGroup := router.Group("/api/v1/users")
	{
		createUserGroup.POST("/", controller.CreateUser)
	}

}

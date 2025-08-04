package routes

import (
	"github.com/BunocGomes/mkp-back/controller"
	"github.com/BunocGomes/mkp-back/middleware"
	"github.com/gin-gonic/gin"
)

func SetupProfileRoutes(router *gin.Engine) {
	profileGroup := router.Group("/api/v1/profiles")
	profileGroup.Use(middleware.AuthMiddleware())
	{

		profileGroup.GET("/:userId", controller.GetProfileByUserID)
		profileGroup.PUT("/me", controller.UpdateMyProfile)
	}
}

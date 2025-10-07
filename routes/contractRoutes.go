package routes

import (
	"github.com/BunocGomes/mkp-back/controller"
	"github.com/BunocGomes/mkp-back/middleware"
	"github.com/gin-gonic/gin"
)

func SetupContractRoutes(router *gin.Engine) {
	contract := router.Group("/contracts")
	contract.Use(middleware.AuthMiddleware())
	{
		contract.GET("/meus", controller.GetMeusContratos)
	}
}
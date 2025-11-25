package routes

import (
	"github.com/BunocGomes/mkp-back/controller"
	"github.com/BunocGomes/mkp-back/middleware"
	"github.com/gin-gonic/gin"
)

func SetupContractRoutes(router *gin.Engine) {
	contract := router.Group("/api/v1/contracts")
	contract.Use(middleware.AuthMiddleware())
	{
		contract.GET("/meus", controller.GetMeusContratos)
		contract.PATCH("/:id/status", controller.UpdateContractStatus)
	}
}

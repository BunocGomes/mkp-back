package routes

import (
	"github.com/BunocGomes/mkp-back/controller"
	"github.com/BunocGomes/mkp-back/middleware"
	"github.com/gin-gonic/gin"
)

func SetupProposaltRoutes(router *gin.Engine) {
	proposal := router.Group("/proposals")
	proposal.Use(middleware.AuthMiddleware())
	{
		proposal.POST("/", controller.CreateProposta)
		proposal.PUT("/:id", controller.UpdateProposta)
		proposal.DELETE("/:id", controller.DeleteProposta)
		proposal.GET("/:id/propostas", controller.GetPropostasByProjeto)
		proposal.POST("/:id/accept", controller.AcceptProposta)
	}
}

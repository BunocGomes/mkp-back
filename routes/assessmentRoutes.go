package routes

import (
	"github.com/BunocGomes/mkp-back/controller"	
	"github.com/BunocGomes/mkp-back/middleware"
	"github.com/gin-gonic/gin"
)

func SetupAvaliacaoRoutes(router *gin.Engine) {
	avaliacoes := router.Group("/api/v1/avaliacoes")
	avaliacoes.Use(middleware.AuthMiddleware()) // Requer autenticação
	{
		avaliacoes.POST("/", controller.CreateAvaliacao)

		avaliacoes.GET("/usuario/:userID", controller.GetAvaliacoesDeUsuario)

	}
}
package routes

import (
	"github.com/BunocGomes/mkp-back/controller"
	"github.com/BunocGomes/mkp-back/middleware"
	"github.com/gin-gonic/gin"
)

func SetupProjectRoutes(router *gin.Engine) {
	// Agrupa todas as rotas de projetos sob /projetos
	projetos := router.Group("/projetos")

	// Aplica os middlewares de autenticação e autorização para o papel "empresa"
	projetos.Use(middleware.AuthMiddleware())

	{
		projetos.POST("/", controller.CreateProject)
		projetos.GET("/", controller.GetProjectsForEmpresa)
		projetos.GET("/:id", controller.GetProjectByID)
		projetos.PUT("/:id", controller.UpdateProject)
		projetos.DELETE("/:id", controller.DeleteProject)
	}
}

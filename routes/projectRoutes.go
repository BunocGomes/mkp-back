package routes

import (
	"github.com/BunocGomes/mkp-back/controller"
	"github.com/BunocGomes/mkp-back/middleware"
	"github.com/gin-gonic/gin"
)

func SetupProjectRoutes(router *gin.Engine) {
	projetos := router.Group("/api/v1/projetos")

	projetos.Use(middleware.AuthMiddleware())

	{
		projetos.POST("/", controller.CreateProject)
		projetos.GET("/", controller.ListProjects)
		projetos.GET("/:id", controller.GetProjectByID)
		projetos.PUT("/:id", controller.UpdateProject)
		projetos.DELETE("/:id", controller.DeleteProject)
	}
}

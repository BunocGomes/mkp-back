package routes

import (
	"github.com/BunocGomes/mkp-back/chat"
	"github.com/BunocGomes/mkp-back/controller"
	"github.com/BunocGomes/mkp-back/middleware"
	"github.com/gin-gonic/gin"
)

func SetupChatRoutes(router *gin.Engine, hub *chat.Hub) {

	chatGroup := router.Group("/chat")
	chatGroup.Use(middleware.AuthMiddleware())
	{
		chatGroup.GET("/history/:userID", controller.GetMessageHistory)
		chatGroup.GET("/ws", func(c *gin.Context) {
			controller.ServeWs(hub, c)
		})
	}
}

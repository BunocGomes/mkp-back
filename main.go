package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"

	"github.com/BunocGomes/mkp-back/chat"
	"github.com/BunocGomes/mkp-back/helper"
	"github.com/BunocGomes/mkp-back/routes"
	"github.com/gin-gonic/gin"
)

var hub *chat.Hub

func main() {
	// A chamada volta a ser simples, sem argumentos.
	helper.InitializeDatabase()
	hub = chat.NewHub()
	go hub.Run()

	log.Println("--- INICIANDO SERVIDOR COM CORS LIBERADO ---") // <--- Procure por isso nos logs!
	router := gin.Default()

	// CONFIGURAÇÃO CORS "NUCLEAR" (PERMISSIVA)
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true, // <--- Isso resolve o problema de origens específicas
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// Configuração das rotas (pode precisar de ajuste se usar injeção de dependência)
	routes.UserRoutes(router)
	routes.SetupProfileRoutes(router)
	routes.SetupProjectRoutes(router)
	routes.SetupProposaltRoutes(router)
	routes.SetupContractRoutes(router)
	routes.SetupChatRoutes(router, hub)
	routes.SetupAvaliacaoRoutes(router)

	log.Println("Servidor iniciado em http://localhost:8080")
	router.Run(":8080")
}

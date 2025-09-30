package main

import (
	"log"

	"github.com/BunocGomes/mkp-back/helper"
	"github.com/BunocGomes/mkp-back/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	// A chamada volta a ser simples, sem argumentos.
	helper.InitializeDatabase()

	log.Println("Iniciando o servidor web...")
	router := gin.Default()

	// Configuração das rotas (pode precisar de ajuste se usar injeção de dependência)
	routes.UserRoutes(router)
	routes.SetupProfileRoutes(router)
	routes.SetupProjectRoutes(router)

	log.Println("Servidor iniciado em http://localhost:8080")
	router.Run(":8080")
}

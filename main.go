package main

import (
	"log"

	"github.com/BunocGomes/mkp-back/helper"
	"github.com/BunocGomes/mkp-back/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Atenção: Não foi possível carregar o arquivo .env")
	}
	helper.InitializeDatabase()
	log.Println("Iniciando o servidor web...")
	router := gin.Default()
	routes.SetupUserRoutes(router)
	log.Println("Servidor iniciado em http://localhost:8080")
	router.Run(":8080")
}

package main

import (
	"github.com/BunocGomes/MKP-back/database" 
	
	"github.com/gin-gonic/gin"
)

func main() {
	// Conecta aos bancos de dados
	database.Connect()      // Conecta ao PostgreSQL
	database.ConnectMongo() // Conecta ao MongoDB

	// Inicia o servidor Gin
	router := gin.Default()

	// Rota de teste para verificar se o servidor está no ar
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.Run(":8080") // Escuta e serve na porta 8080
}
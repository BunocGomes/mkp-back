package controller

import (
	"net/http"
	"strconv"

	"github.com/BunocGomes/mkp-back/dto"
	"github.com/BunocGomes/mkp-back/service"
	"github.com/gin-gonic/gin"
)

func CreateAvaliacao(c *gin.Context) {
	var input dto.CreateAvaliacaoDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	avaliadorUserID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Contexto de usuário 'userId' não encontrado"})
		return
	}

	avaliadorRole, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Contexto de usuário 'role' não encontrado"})
		return
	}

	_, err := service.CreateAvaliacao(input, avaliadorUserID.(uint), avaliadorRole.(string))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Avaliação enviada com sucesso"})
}

func GetAvaliacoesDeUsuario(c *gin.Context) {
	userIDStr := c.Param("userID")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuário inválido"})
		return
	}

	avaliacoes, err := service.GetAvaliacoesDeUsuario(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, avaliacoes)
}
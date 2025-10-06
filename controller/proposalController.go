package controller

import (
	"net/http"
	"strconv"

	"github.com/BunocGomes/mkp-back/dto"
	"github.com/BunocGomes/mkp-back/service"
	"github.com/gin-gonic/gin"
)

// CreateProposta manipula a criação de uma nova proposta por um freelancer.
func CreateProposta(c *gin.Context) {
	freelancerID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ID do freelancer não encontrado no token"})
		return
	}

	var input dto.CreatePropostaDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	proposta, err := service.CreateProposta(input, freelancerID.(uint))
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, proposta)
}

// GetPropostasByProjeto manipula a listagem de propostas de um projeto para a empresa dona.
func GetPropostasByProjeto(c *gin.Context) {
	empresaID, exists := c.Get("empresaId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ID da empresa não encontrado no token"})
		return
	}

	projetoID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do projeto inválido"})
		return
	}

	propostas, err := service.GetPropostasByProjetoID(uint(projetoID), empresaID.(uint))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, propostas)
}

func UpdateProposta(c *gin.Context) {
	freelancerID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ID do freelancer não encontrado no token"})
		return
	}

	propostaID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID da proposta inválido"})
		return
	}

	var input dto.UpdatePropostaDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	proposta, err := service.UpdateProposta(uint(propostaID), freelancerID.(uint), input)
	if err != nil {
		// Usa StatusForbidden para erros de permissão ou regra de negócio
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, proposta)
}

// DeleteProposta manipula a exclusão de uma proposta.
func DeleteProposta(c *gin.Context) {
	freelancerID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ID do freelancer não encontrado no token"})
		return
	}

	propostaID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID da proposta inválido"})
		return
	}

	err = service.DeleteProposta(uint(propostaID), freelancerID.(uint))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Proposta deletada com sucesso"})
}

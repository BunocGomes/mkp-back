package controller

import (
	"net/http"
	"strconv"

	"github.com/BunocGomes/mkp-back/service"
	"github.com/gin-gonic/gin"
)

// GetMeusContratos manipula a listagem de contratos do usuário logado.
func GetMeusContratos(c *gin.Context) {
	userID, _ := c.Get("userId")
	role, _ := c.Get("role")

	contratos, err := service.GetMeusContratos(userID.(uint), role.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, contratos)
}

func UpdateContractStatus(c *gin.Context) {
	userID, _ := c.Get("userId") // ID do usuário logado

	contratoIDStr := c.Param("id")
	contratoID, err := strconv.ParseUint(contratoIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do contrato inválido"})
		return
	}

	var input struct {
		Status string `json:"status" binding:"required"` // "concluido" ou "cancelado"
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Chama o service (você precisará ajustar a validação de permissão no service com base no userID)
	err = service.UpdateContractStatus(uint(contratoID), userID.(uint), input.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status do contrato atualizado com sucesso"})
}

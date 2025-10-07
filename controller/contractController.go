package controller

import (
	"net/http"

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
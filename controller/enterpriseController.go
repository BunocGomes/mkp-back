package controller

import (
	"net/http"

	"github.com/BunocGomes/mkp-back/dto"
	"github.com/BunocGomes/mkp-back/service"
	"github.com/gin-gonic/gin"
)

// CreateEmpresaAccount é o handler para registrar uma nova empresa.
func CreateEmpresaAccount(c *gin.Context) {
	var empresaDTO dto.RegisterEmpresaDTO
	if err := c.ShouldBindJSON(&empresaDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newEmpresa, err := service.CreateEmpresaAccount(empresaDTO)
	if err != nil {
		if err.Error() == "e-mail já cadastrado" || err.Error() == "CNPJ já cadastrado" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao criar conta de empresa"})
		return
	}

	c.JSON(http.StatusCreated, newEmpresa)
}

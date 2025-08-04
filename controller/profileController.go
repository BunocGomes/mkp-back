package controller

import (
	"net/http"
	"strconv"

	"github.com/BunocGomes/mkp-back/dto"
	"github.com/BunocGomes/mkp-back/service"
	"github.com/gin-gonic/gin"
)

// GetProfileByUserID é o handler para buscar um perfil público.
func GetProfileByUserID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuário inválido"})
		return
	}

	profile, err := service.GetProfileByUserID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, profile)
}

// UpdateMyProfile é o handler para o usuário logado atualizar seu próprio perfil.
func UpdateMyProfile(c *gin.Context) {
	// 1. O userID é injetado pelo middleware de autenticação.
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Contexto de usuário inválido"})
		return
	}
	userIDUint := userID.(uint)

	// 2. Faz o bind do JSON do corpo da requisição para o DTO.
	var profileDTO dto.UpdateProfileDTO
	if err := c.ShouldBindJSON(&profileDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. Chama o serviço para executar a atualização.
	updatedProfile, err := service.UpdateProfile(userIDUint, profileDTO)
	if err != nil {
		if err.Error() == "perfil não encontrado" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao atualizar o perfil", "details": err.Error()})
		return
	}

	// 4. Retorna o perfil atualizado com sucesso.
	c.JSON(http.StatusOK, updatedProfile)
}

package controller

import (
	"net/http"
	"strconv"

	"github.com/BunocGomes/mkp-back/dto"
	"github.com/BunocGomes/mkp-back/helper"
	"github.com/BunocGomes/mkp-back/service"
	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	var userDTO dto.RegisterUserDTO
	if err := c.ShouldBindJSON(&userDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	createdUser, err := service.CreateUser(userDTO)
	if err != nil {
		if err.Error() == "e-mail já cadastrado" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao criar usuário"})
		return
	}
	c.JSON(http.StatusCreated, createdUser)
}

func GetAllUsers(c *gin.Context) {
	users, err := service.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao buscar usuários"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	user, err := service.GetUserByID(uint(id))
	if err != nil {
		if err.Error() == "usuário não encontrado" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao buscar usuário"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	var userDTO dto.UpdateUserDTO
	if err := c.ShouldBindJSON(&userDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedUser, err := service.UpdateUser(uint(id), userDTO)
	if err != nil {
		if err.Error() == "usuário não encontrado" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao atualizar usuário"})
		return
	}
	c.JSON(http.StatusOK, updatedUser)
}

func DeleteUser(c *gin.Context) {

	token := helper.ExtractToken(c)

	// 2. Verifica se o usuário tem a permissão de SuperAdmin
	isAllowed, err := helper.IsSuperAdmin(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido ou expirado"})
		return
	}
	if !isAllowed {
		c.JSON(http.StatusForbidden, gin.H{"error": "Acesso negado. Apenas SuperAdmins podem deletar usuários."})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	err = service.DeleteUser(uint(id))
	if err != nil {
		if err.Error() == "usuário não encontrado" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao deletar usuário"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func Login(c *gin.Context) {
	var loginDTO dto.LoginDTO
	if err := c.ShouldBindJSON(&loginDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := service.Login(loginDTO)
	if err != nil {
		// Retorna 401 Unauthorized para qualquer erro de login
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	response := dto.LoginResponseDTO{Token: token}
	c.JSON(http.StatusOK, response)
}

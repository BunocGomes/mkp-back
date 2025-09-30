package controller

import (
	"net/http"
	"strconv"

	"github.com/BunocGomes/mkp-back/dto"
	"github.com/BunocGomes/mkp-back/service"
	"github.com/gin-gonic/gin"
)

// CreateProject manipula a criação de um novo projeto.
func CreateProject(c *gin.Context) {
	empresaID, exists := c.Get("empresaId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ID da empresa não encontrado no token"})
		return
	}

	var input dto.CreateProjectDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	projeto, err := service.CreateProject(input, empresaID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, projeto)
}

// GetProjectsForEmpresa manipula a listagem de projetos da empresa logada.
func GetProjectsForEmpresa(c *gin.Context) {
	empresaID, exists := c.Get("empresaId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ID da empresa não encontrado no token"})
		return
	}

	projetos, err := service.GetProjectsByEmpresaID(empresaID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, projetos)
}

// GetProjectByID manipula a busca por um projeto específico.
func GetProjectByID(c *gin.Context) {
	empresaID, exists := c.Get("empresaId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ID da empresa não encontrado no token"})
		return
	}

	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do projeto inválido"})
		return
	}

	projeto, err := service.GetProjectByID(uint(projectID), empresaID.(uint))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, projeto)
}

// UpdateProject manipula a atualização de um projeto.
func UpdateProject(c *gin.Context) {
	empresaID, exists := c.Get("empresaId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ID da empresa não encontrado no token"})
		return
	}

	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do projeto inválido"})
		return
	}

	var input dto.UpdateProjectDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	projeto, err := service.UpdateProject(uint(projectID), empresaID.(uint), input)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, projeto)
}

// DeleteProject manipula a exclusão de um projeto.
func DeleteProject(c *gin.Context) {
	empresaID, exists := c.Get("empresaId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ID da empresa não encontrado no token"})
		return
	}

	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do projeto inválido"})
		return
	}

	err = service.DeleteProject(uint(projectID), empresaID.(uint))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Projeto deletado com sucesso"})
}

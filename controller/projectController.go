package controller

import (
	"net/http"
	"strconv"

	"github.com/BunocGomes/mkp-back/dto"
	"github.com/BunocGomes/mkp-back/service"
	"github.com/gin-gonic/gin"
)

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

func ListProjects(c *gin.Context) {
	role, _ := c.Get("role")

	// Se for uma empresa, retorna apenas os projetos dela
	if role.(string) == "empresa" {
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
		return
	}

	// Para freelancers e outros, retorna a busca pública
	searchTerm := c.DefaultQuery("search", "")

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	projetos, err := service.SearchOpenProjects(searchTerm, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, projetos)
}

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

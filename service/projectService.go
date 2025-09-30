package service

import (
	"errors"
	"fmt"

	"github.com/BunocGomes/mkp-back/database"
	"github.com/BunocGomes/mkp-back/dto"
	"github.com/BunocGomes/mkp-back/models"
)

// CreateProject cria um novo projeto para uma empresa.
func CreateProject(input dto.CreateProjectDTO, empresaID uint) (dto.ProjectResponseDTO, error) {
	projeto := models.Projeto{
		Titulo:    input.Titulo,
		Descricao: input.Descricao,
		Orcamento: input.Orcamento,
		EmpresaID: empresaID, // Associa o projeto à empresa do usuário logado
	}

	if err := database.DB.Create(&projeto).Error; err != nil {
		return dto.ProjectResponseDTO{}, err
	}

	response := dto.ProjectResponseDTO{
		ID:        projeto.ID,
		Titulo:    projeto.Titulo,
		Descricao: projeto.Descricao,
		Status:    projeto.Status,
		Orcamento: projeto.Orcamento,
		EmpresaID: projeto.EmpresaID,
		CreatedAt: projeto.CreatedAt,
		UpdatedAt: projeto.UpdatedAt,
	}

	return response, nil
}

// GetProjectsByEmpresaID retorna todos os projetos de uma empresa específica.
func GetProjectsByEmpresaID(empresaID uint) ([]dto.ProjectResponseDTO, error) {
	var projetos []models.Projeto
	if err := database.DB.Where("empresa_id = ?", empresaID).Order("created_at desc").Find(&projetos).Error; err != nil {
		return nil, err
	}

	var response []dto.ProjectResponseDTO
	for _, p := range projetos {
		response = append(response, dto.ProjectResponseDTO{
			ID:        p.ID,
			Titulo:    p.Titulo,
			Descricao: p.Descricao,
			Status:    p.Status,
			Orcamento: p.Orcamento,
			EmpresaID: p.EmpresaID,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		})
	}
	return response, nil
}

// GetProjectByID retorna um projeto específico, verificando a posse.
func GetProjectByID(projectID, empresaID uint) (dto.ProjectResponseDTO, error) {
	var projeto models.Projeto
	if err := database.DB.First(&projeto, projectID).Error; err != nil {
		return dto.ProjectResponseDTO{}, errors.New("projeto não encontrado")
	}

	if projeto.EmpresaID != empresaID {
		return dto.ProjectResponseDTO{}, errors.New("acesso negado: este projeto não pertence à sua empresa")
	}

	response := dto.ProjectResponseDTO{
		ID:        projeto.ID,
		Titulo:    projeto.Titulo,
		Descricao: projeto.Descricao,
		Status:    projeto.Status,
		Orcamento: projeto.Orcamento,
		EmpresaID: projeto.EmpresaID,
		CreatedAt: projeto.CreatedAt,
		UpdatedAt: projeto.UpdatedAt,
	}

	return response, nil
}

// UpdateProject atualiza um projeto, verificando a posse.
func UpdateProject(projectID, empresaID uint, input dto.UpdateProjectDTO) (dto.ProjectResponseDTO, error) {
	var projeto models.Projeto
	if err := database.DB.First(&projeto, projectID).Error; err != nil {
		return dto.ProjectResponseDTO{}, errors.New("projeto não encontrado")
	}

	if projeto.EmpresaID != empresaID {
		return dto.ProjectResponseDTO{}, errors.New("acesso negado: você não pode atualizar este projeto")
	}

	if err := database.DB.Model(&projeto).Updates(input).Error; err != nil {
		return dto.ProjectResponseDTO{}, err
	}

	return GetProjectByID(projectID, empresaID)
}

// DeleteProject deleta um projeto, verificando a posse.
func DeleteProject(projectID, empresaID uint) error {
	var projeto models.Projeto
	if err := database.DB.First(&projeto, projectID).Error; err != nil {
		return errors.New("projeto não encontrado")
	}

	if projeto.EmpresaID != empresaID {
		return errors.New("acesso negado: você não pode deletar este projeto")
	}

	if err := database.DB.Delete(&projeto).Error; err != nil {
		return fmt.Errorf("falha ao deletar o projeto: %w", err)
	}

	return nil
}

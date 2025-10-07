package service

import (
	"github.com/BunocGomes/mkp-back/database"
	"github.com/BunocGomes/mkp-back/dto"
	"github.com/BunocGomes/mkp-back/models"
)

func GetMeusContratos(userID uint, role string) ([]dto.ContratoResponseDTO, error) {
	var contratos []models.Contrato
	var err error

	query := database.DB.Preload("Projeto").Preload("Empresa").Preload("Freelancer")

	if role == "empresa" {
		var user models.User
		if err := database.DB.Select("empresa_id").First(&user, userID).Error; err != nil {
			return nil, err
		}
		err = query.Where("empresa_id = ?", user.EmpresaID).Order("created_at desc").Find(&contratos).Error
	} else {
		err = query.Where("freelancer_id = ?", userID).Order("created_at desc").Find(&contratos).Error
	}

	if err != nil {
		return nil, err
	}

	var response []dto.ContratoResponseDTO
	for _, c := range contratos {
		response = append(response, dto.ContratoResponseDTO{
			ID:              c.ID,
			ValorFinal:      c.ValorFinal,
			DataInicio:      c.DataInicio,
			DataFimPrevista: c.DataFimPrevista,
			Status:          c.Status,
			ProjetoID:       c.ProjetoID,
			ProjetoTitulo:   c.Projeto.Titulo,
			EmpresaNome:     c.Empresa.NomeFantasia,
			FreelancerNome:  c.Freelancer.Nome,
		})
	}

	return response, nil
}

package service

import (
	"errors"
	"time"

	"github.com/BunocGomes/mkp-back/database"
	"github.com/BunocGomes/mkp-back/dto"
	"github.com/BunocGomes/mkp-back/models"
	"gorm.io/gorm"
)

func CreateProposta(input dto.CreatePropostaDTO, freelancerID uint) (dto.PropostaResponseDTO, error) {
	var projeto models.Projeto
	if err := database.DB.First(&projeto, input.ProjetoID).Error; err != nil {
		return dto.PropostaResponseDTO{}, errors.New("projeto não encontrado")
	}
	if projeto.Status != "aberto" {
		return dto.PropostaResponseDTO{}, errors.New("este projeto não está mais aceitando propostas")
	}

	var existente models.Proposta
	err := database.DB.Where("projeto_id = ? AND freelancer_id = ?", input.ProjetoID, freelancerID).First(&existente).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.PropostaResponseDTO{}, errors.New("você já enviou uma proposta para este projeto")
	}

	proposta := models.Proposta{
		ProjetoID:     input.ProjetoID,
		FreelancerID:  freelancerID,
		ValorProposto: input.ValorProposto,
		Descricao:     input.Descricao,
		PrazoEstimado: input.PrazoEstimado,
	}

	if err := database.DB.Create(&proposta).Error; err != nil {
		return dto.PropostaResponseDTO{}, err
	}

	var freelancer models.User
	database.DB.First(&freelancer, freelancerID)

	response := dto.PropostaResponseDTO{
		ID:            proposta.ID,
		ValorProposto: proposta.ValorProposto,
		Descricao:     proposta.Descricao,
		PrazoEstimado: proposta.PrazoEstimado,
		Status:        proposta.Status,
		ProjetoID:     proposta.ProjetoID,
		CreatedAt:     proposta.CreatedAt,
		Freelancer: struct {
			ID   uint   `json:"id"`
			Nome string `json:"nome"`
		}{
			ID:   freelancer.ID,
			Nome: freelancer.Nome,
		},
	}

	return response, nil
}

func GetPropostasByProjetoID(projetoID, empresaID uint) ([]dto.PropostaResponseDTO, error) {
	var projeto models.Projeto
	if err := database.DB.Select("empresa_id").First(&projeto, projetoID).Error; err != nil {
		return nil, errors.New("projeto não encontrado")
	}
	if projeto.EmpresaID != empresaID {
		return nil, errors.New("acesso negado: este projeto não pertence à sua empresa")
	}

	var propostas []models.Proposta
	if err := database.DB.Preload("Freelancer").Where("projeto_id = ?", projetoID).Find(&propostas).Error; err != nil {
		return nil, err
	}

	var response []dto.PropostaResponseDTO
	for _, p := range propostas {
		response = append(response, dto.PropostaResponseDTO{
			ID:            p.ID,
			ValorProposto: p.ValorProposto,
			Descricao:     p.Descricao,
			PrazoEstimado: p.PrazoEstimado,
			Status:        p.Status,
			ProjetoID:     p.ProjetoID,
			CreatedAt:     p.CreatedAt,
			Freelancer: struct {
				ID   uint   `json:"id"`
				Nome string `json:"nome"`
			}{
				ID:   p.FreelancerID,
				Nome: p.Freelancer.Nome,
			},
		})
	}

	return response, nil
}

func UpdateProposta(propostaID uint, freelancerID uint, input dto.UpdatePropostaDTO) (dto.PropostaResponseDTO, error) {
	var proposta models.Proposta
	if err := database.DB.First(&proposta, propostaID).Error; err != nil {
		return dto.PropostaResponseDTO{}, errors.New("proposta não encontrada")
	}

	if proposta.FreelancerID != freelancerID {
		return dto.PropostaResponseDTO{}, errors.New("acesso negado: você não é o autor desta proposta")
	}

	if proposta.Status != "enviada" {
		return dto.PropostaResponseDTO{}, errors.New("esta proposta não pode mais ser editada, pois já foi visualizada ou respondida")
	}

	if input.ValorProposto != nil {
		proposta.ValorProposto = *input.ValorProposto
	}
	if input.Descricao != nil {
		proposta.Descricao = *input.Descricao
	}
	if input.PrazoEstimado != nil {
		proposta.PrazoEstimado = *input.PrazoEstimado
	}

	if err := database.DB.Save(&proposta).Error; err != nil {
		return dto.PropostaResponseDTO{}, err
	}

	var freelancer models.User
	database.DB.First(&freelancer, freelancerID)

	response := dto.PropostaResponseDTO{
		ID:            proposta.ID,
		ValorProposto: proposta.ValorProposto,
		Descricao:     proposta.Descricao,
		PrazoEstimado: proposta.PrazoEstimado,
		Status:        proposta.Status,
		ProjetoID:     proposta.ProjetoID,
		CreatedAt:     proposta.CreatedAt,
		Freelancer: struct {
			ID   uint   `json:"id"`
			Nome string `json:"nome"`
		}{
			ID:   freelancer.ID,
			Nome: freelancer.Nome,
		},
	}

	return response, nil
}

func DeleteProposta(propostaID uint, freelancerID uint) error {
	var proposta models.Proposta
	if err := database.DB.First(&proposta, propostaID).Error; err != nil {
		return errors.New("proposta não encontrada")
	}

	if proposta.FreelancerID != freelancerID {
		return errors.New("acesso negado: você não é o autor desta proposta")
	}

	if proposta.Status != "enviada" {
		return errors.New("esta proposta não pode mais ser deletada")
	}

	if err := database.DB.Delete(&proposta).Error; err != nil {
		return err
	}

	return nil
}

func AceitarProposta(propostaID uint, empresaID uint) (dto.ContratoResponseDTO, error) {
	tx := database.DB.Begin()
	if tx.Error != nil {
		return dto.ContratoResponseDTO{}, tx.Error
	}

	var proposta models.Proposta
	if err := tx.Preload("Projeto").First(&proposta, propostaID).Error; err != nil {
		tx.Rollback()
		return dto.ContratoResponseDTO{}, errors.New("proposta não encontrada")
	}

	if proposta.Projeto.EmpresaID != empresaID {
		tx.Rollback()
		return dto.ContratoResponseDTO{}, errors.New("acesso negado: este projeto não pertence à sua empresa")
	}

	if proposta.Status != "enviada" || proposta.Projeto.Status != "aberto" {
		tx.Rollback()
		return dto.ContratoResponseDTO{}, errors.New("esta proposta ou projeto não está mais disponível para contratação")
	}

	contrato := models.Contrato{
		ValorFinal:      proposta.ValorProposto,
		DataInicio:      time.Now(),
		DataFimPrevista: time.Now().AddDate(0, 0, proposta.PrazoEstimado),
		Status:          "ativo",
		ProjetoID:       proposta.ProjetoID,
		EmpresaID:       empresaID,
		FreelancerID:    proposta.FreelancerID,
		PropostaID:      proposta.ID,
	}
	if err := tx.Create(&contrato).Error; err != nil {
		tx.Rollback()
		return dto.ContratoResponseDTO{}, err
	}

	if err := tx.Model(&proposta).Update("status", "aceita").Error; err != nil {
		tx.Rollback()
		return dto.ContratoResponseDTO{}, err
	}

	if err := tx.Model(&models.Projeto{}).Where("id = ?", proposta.ProjetoID).Update("status", "em_andamento").Error; err != nil {
		tx.Rollback()
		return dto.ContratoResponseDTO{}, err
	}

	if err := tx.Model(&models.Proposta{}).Where("projeto_id = ? AND id != ?", proposta.ProjetoID, proposta.ID).Update("status", "recusada").Error; err != nil {
		tx.Rollback()
		return dto.ContratoResponseDTO{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return dto.ContratoResponseDTO{}, err
	}

	var freelancer models.User
	var empresa models.Empresa
	database.DB.First(&freelancer, contrato.FreelancerID)
	database.DB.First(&empresa, contrato.EmpresaID)

	response := dto.ContratoResponseDTO{
		ID:              contrato.ID,
		ValorFinal:      contrato.ValorFinal,
		DataInicio:      contrato.DataInicio,
		DataFimPrevista: contrato.DataFimPrevista,
		Status:          contrato.Status,
		ProjetoID:       contrato.ProjetoID,
		ProjetoTitulo:   proposta.Projeto.Titulo,
		EmpresaNome:     empresa.NomeFantasia,
		FreelancerNome:  freelancer.Nome,
	}

	return response, nil
}

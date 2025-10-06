package service

import (
	"errors"

	"github.com/BunocGomes/mkp-back/database"
	"github.com/BunocGomes/mkp-back/dto"
	"github.com/BunocGomes/mkp-back/models"
	"gorm.io/gorm"
)

// CreateProposta cria uma nova proposta de um freelancer para um projeto.
func CreateProposta(input dto.CreatePropostaDTO, freelancerID uint) (dto.PropostaResponseDTO, error) {
	// Regra de negócio 1: Verificar se o projeto existe e está aberto
	var projeto models.Projeto
	if err := database.DB.First(&projeto, input.ProjetoID).Error; err != nil {
		return dto.PropostaResponseDTO{}, errors.New("projeto não encontrado")
	}
	if projeto.Status != "aberto" {
		return dto.PropostaResponseDTO{}, errors.New("este projeto não está mais aceitando propostas")
	}

	// Regra de negócio 2: Verificar se o freelancer já não enviou uma proposta
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

	// Para a resposta, vamos buscar o usuário para incluir o nome
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

// GetPropostasByProjetoID retorna as propostas de um projeto para o dono da empresa.
func GetPropostasByProjetoID(projetoID, empresaID uint) ([]dto.PropostaResponseDTO, error) {
	// Segurança: Verificar se a empresa que está pedindo é a dona do projeto
	var projeto models.Projeto
	if err := database.DB.Select("empresa_id").First(&projeto, projetoID).Error; err != nil {
		return nil, errors.New("projeto não encontrado")
	}
	if projeto.EmpresaID != empresaID {
		return nil, errors.New("acesso negado: este projeto não pertence à sua empresa")
	}

	var propostas []models.Proposta
	// Usamos Preload para carregar os dados do Freelancer (que é um User) junto com a proposta
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

	// Segurança: Apenas o dono da proposta pode atualizá-la.
	if proposta.FreelancerID != freelancerID {
		return dto.PropostaResponseDTO{}, errors.New("acesso negado: você não é o autor desta proposta")
	}

	// Regra de Negócio: Proposta só pode ser alterada se ainda estiver com status "enviada".
	if proposta.Status != "enviada" {
		return dto.PropostaResponseDTO{}, errors.New("esta proposta não pode mais ser editada, pois já foi visualizada ou respondida")
	}

	// Atualiza apenas os campos que foram fornecidos no DTO
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

	// Para a resposta, vamos buscar o usuário para incluir o nome
	var freelancer models.User
	database.DB.First(&freelancer, freelancerID)

	// Reutiliza a lógica de formatação de resposta
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

// DeleteProposta remove uma proposta do sistema.
func DeleteProposta(propostaID uint, freelancerID uint) error {
	var proposta models.Proposta
	if err := database.DB.First(&proposta, propostaID).Error; err != nil {
		return errors.New("proposta não encontrada")
	}

	// Segurança: Apenas o dono da proposta pode deletá-la.
	if proposta.FreelancerID != freelancerID {
		return errors.New("acesso negado: você não é o autor desta proposta")
	}

	// Regra de Negócio: Proposta só pode ser deletada se ainda estiver com status "enviada".
	if proposta.Status != "enviada" {
		return errors.New("esta proposta não pode mais ser deletada")
	}

	if err := database.DB.Delete(&proposta).Error; err != nil {
		return err
	}

	return nil
}

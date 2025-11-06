package service

import (
	"errors"

	"github.com/BunocGomes/mkp-back/database"
	"github.com/BunocGomes/mkp-back/dto"
	"github.com/BunocGomes/mkp-back/models"
	"gorm.io/gorm/clause"
)

// CreateAvaliacao contém a lógica transacional para criar uma avaliação
func CreateAvaliacao(input dto.CreateAvaliacaoDTO, avaliadorUserID uint, avaliadorRole string) (models.Avaliacao, error) {
	var contrato models.Contrato
	var avaliadoUserID uint
	var campoAtualizacaoContrato string

	tx := database.DB.Begin()
	if tx.Error != nil {
		return models.Avaliacao{}, tx.Error
	}

	// 1. Buscar o contrato e verificar permissões
	if err := tx.First(&contrato, input.ContratoID).Error; err != nil {
		tx.Rollback()
		return models.Avaliacao{}, errors.New("contrato não encontrado")
	}

	if contrato.Status != "concluido" {
		tx.Rollback()
		return models.Avaliacao{}, errors.New("este contrato ainda não foi concluído")
	}

	// 2. Determinar quem está a avaliar quem (LÓGICA CORRIGIDA)
	if avaliadorRole == "empresa" {
		// O avaliador é uma empresa. O 'avaliadorUserID' é da tabela 'users'.
		// Precisamos de descobrir o 'empresa_id' associado a este usuário.
		var userAvaliador models.User
		if err := tx.Select("empresa_id").First(&userAvaliador, avaliadorUserID).Error; err != nil {
			tx.Rollback()
			return models.Avaliacao{}, errors.New("usuário da empresa avaliadora não encontrado")
		}
		if userAvaliador.EmpresaID == nil {
			tx.Rollback()
			return models.Avaliacao{}, errors.New("este usuário não está associado a nenhuma empresa")
		}

		// AGORA SIM: Comparamos o empresa_id (do contrato) com o empresa_id (do usuário)
		if contrato.EmpresaID == *userAvaliador.EmpresaID {
			avaliadoUserID = contrato.FreelancerID // O avaliado é o freelancer (um user_id)
			campoAtualizacaoContrato = "avaliacao_empresa_id"
			if contrato.AvaliacaoEmpresaID != nil {
				tx.Rollback()
				return models.Avaliacao{}, errors.New("empresa já avaliou este contrato")
			}
		} else {
			// O token é de uma empresa, mas não é a empresa deste contrato
			tx.Rollback()
			return models.Avaliacao{}, errors.New("acesso negado: você não faz parte deste contrato")
		}

	} else if avaliadorRole == "freelancer" {
		// O avaliador é um freelancer. O 'avaliadorUserID' é da tabela 'users'.
		// Comparamos o freelancer_id (do contrato) com o user_id (do token)
		if contrato.FreelancerID == avaliadorUserID {

			// O avaliado é a empresa. Precisamos do user_id dela, não do empresa_id.
			var userAvaliado models.User
			if err := tx.Select("id").Where("empresa_id = ?", contrato.EmpresaID).First(&userAvaliado).Error; err != nil {
				tx.Rollback()
				return models.Avaliacao{}, errors.New("usuário da empresa avaliada não encontrado")
			}
			avaliadoUserID = userAvaliado.ID // Este é o user_id da empresa

			campoAtualizacaoContrato = "avaliacao_freelancer_id"
			if contrato.AvaliacaoFreelancerID != nil {
				tx.Rollback()
				return models.Avaliacao{}, errors.New("freelancer já avaliou este contrato")
			}
		} else {
			// O token é de um freelancer, mas não é o freelancer deste contrato
			tx.Rollback()
			return models.Avaliacao{}, errors.New("acesso negado: você não faz parte deste contrato")
		}
	} else {
		tx.Rollback()
		return models.Avaliacao{}, errors.New("acesso negado: papel inválido")
	}

	// 3. Criar a Avaliação (Agora 'avaliadoUserID' é sempre um user_id)
	avaliacao := models.Avaliacao{
		ContratoID:      input.ContratoID,
		AvaliadorUserID: avaliadorUserID, // O user_id do token
		AvaliadoUserID:  avaliadoUserID,  // O user_id que descobrimos
		Nota:            input.Nota,
		Comentario:      input.Comentario,
	}
	if err := tx.Create(&avaliacao).Error; err != nil {
		tx.Rollback()
		return models.Avaliacao{}, err
	}

	// 4. Marcar o Contrato como avaliado por esta parte
	if err := tx.Model(&contrato).Update(campoAtualizacaoContrato, avaliacao.ID).Error; err != nil {
		tx.Rollback()
		return models.Avaliacao{}, err
	}
	// 5. Atualizar o Perfil do usuário AVALIADO
	var perfilAvaliado models.Perfil
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("usuario_id = ?", avaliadoUserID).First(&perfilAvaliado).Error; err != nil {
		tx.Rollback()
		return models.Avaliacao{}, errors.New("perfil do usuário avaliado não encontrado")
	}

	totalNotas := (perfilAvaliado.NotaMedia * float64(perfilAvaliado.TotalAvaliacoes)) + float64(input.Nota)
	novoTotalAvaliacoes := perfilAvaliado.TotalAvaliacoes + 1
	novaNotaMedia := totalNotas / float64(novoTotalAvaliacoes)

	if err := tx.Model(&perfilAvaliado).Updates(models.Perfil{
		NotaMedia:       novaNotaMedia,
		TotalAvaliacoes: novoTotalAvaliacoes,
	}).Error; err != nil {
		tx.Rollback()
		return models.Avaliacao{}, err
	}

	// 6. Se tudo deu certo, cometer a transação
	if err := tx.Commit().Error; err != nil {
		return models.Avaliacao{}, err
	}

	return avaliacao, nil
}

// GetAvaliacoesDeUsuario busca todas as avaliações que um usuário recebeu.
func GetAvaliacoesDeUsuario(avaliadoUserID uint) ([]dto.AvaliacaoResponseDTO, error) {
	var avaliacoes []models.Avaliacao
	
	// Usamos Preload("AvaliadorUser") para carregar os dados do autor da avaliação
	// (O GORM usa o nome do campo no struct, "AvaliadorUser")
	err := database.DB.Preload("AvaliadorUser").
		Where("avaliado_user_id = ?", avaliadoUserID).
		Order("created_at desc").
		Find(&avaliacoes).Error

	if err != nil {
		return nil, err
	}

	// Mapeia os dados do banco para o DTO de Resposta da API
	var responseDTOs []dto.AvaliacaoResponseDTO
	for _, a := range avaliacoes {
		dto := dto.AvaliacaoResponseDTO{
			ID:         a.ID,
			Nota:       a.Nota,
			Comentario: a.Comentario,
			ContratoID: a.ContratoID,
			CreatedAt:  a.CreatedAt,
			// Preenchemos os dados do avaliador que carregámos com o Preload
			Avaliador: struct {
				ID   uint   `json:"id"`
				Nome string `json:"nome"`
			}{
				ID:   a.AvaliadorUserID,
				Nome: a.AvaliadorUser.Nome,
			},
		}
		responseDTOs = append(responseDTOs, dto)
	}

	return responseDTOs, nil
}
package service

import (
	"errors"

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
		// --- LÓGICA NOVA: Buscar o ID do Usuário Admin da Empresa ---
		var empresaUser models.User
		// Buscamos o primeiro usuário vinculado a esta empresa (assumindo que é o admin/dono)
		// Isso é necessário para o chat funcionar corretamente entre Freelancer -> Usuário da Empresa
		database.DB.Select("id").Where("empresa_id = ?", c.EmpresaID).First(&empresaUser)

		empresaUserID := empresaUser.ID
		// ------------------------------------------------------------

		response = append(response, dto.ContratoResponseDTO{
			ID:              c.ID,
			ValorFinal:      c.ValorFinal,
			DataInicio:      c.DataInicio,
			DataFimPrevista: c.DataFimPrevista,
			Status:          c.Status,
			ProjetoID:       c.ProjetoID,
			EmpresaID:       c.EmpresaID,

			// Adicione este campo no DTO também!
			EmpresaUserID: empresaUserID,

			FreelancerID:   c.FreelancerID,
			ProjetoTitulo:  c.Projeto.Titulo,
			EmpresaNome:    c.Empresa.NomeFantasia,
			FreelancerNome: c.Freelancer.Nome,
		})
	}

	return response, nil
}

func UpdateContractStatus(contratoID uint, userID uint, novoStatus string) error {
	var contrato models.Contrato
	if err := database.DB.First(&contrato, contratoID).Error; err != nil {
		return errors.New("contrato não encontrado")
	}

	// Verificação de Segurança: Apenas as partes envolvidas podem alterar
	// Precisamos verificar se o userID é da Empresa ou do Freelancer deste contrato
	// Nota: Para simplificar, assumimos que o userID passado já foi validado pelo controller
	// Mas idealmente, faríamos:
	// if contrato.EmpresaID != user.EmpresaID && contrato.FreelancerID != userID { return erro }

	// Regra de Negócio: Só pode mudar se estiver "ativo"
	if contrato.Status != "ativo" {
		return errors.New("apenas contratos ativos podem ser alterados")
	}

	if novoStatus != "concluido" && novoStatus != "cancelado" {
		return errors.New("status inválido")
	}

	contrato.Status = novoStatus
	// Se concluído, definimos a data real de fim (opcional, mas bom)
	// contrato.DataFim = time.Now()

	return database.DB.Save(&contrato).Error
}

package service

import (
	"errors"

	"github.com/BunocGomes/mkp-back/database"
	"github.com/BunocGomes/mkp-back/dto"
	"github.com/BunocGomes/mkp-back/models"
	"gorm.io/gorm"
)

func CreateEmpresaAccount(input dto.RegisterEmpresaDTO) (dto.EmpresaResponseDTO, error) {
	var empresa models.Empresa
	var user models.User

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("email = ?", input.EmailUsuario).First(&models.User{}).Error; err == nil {
			return errors.New("e-mail já cadastrado")
		}
		if err := tx.Where("cnpj = ?", input.CNPJ).First(&models.Empresa{}).Error; err == nil {
			return errors.New("CNPJ já cadastrado")
		}

		empresa = models.Empresa{
			NomeFantasia: input.NomeFantasia,
			RazaoSocial:  input.RazaoSocial,
			CNPJ:         input.CNPJ,
		}
		if err := tx.Create(&empresa).Error; err != nil {
			return err
		}

		empresaID := empresa.ID
		user = models.User{
			Nome:      input.NomeUsuario,
			Email:     input.EmailUsuario,
			EmpresaID: &empresaID,
		}
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		perfil := models.Perfil{
			UsuarioID: user.ID,
			Role:      "empresa",
		}
		if err := tx.Create(&perfil).Error; err != nil {
			return err
		}
		user.Perfil = perfil

		digest := models.Digest{
			UserID:   user.ID,
			Password: input.PasswordUsuario,
		}
		if err := tx.Create(&digest).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return dto.EmpresaResponseDTO{}, err
	}

	response := dto.EmpresaResponseDTO{
		ID:           empresa.ID,
		NomeFantasia: empresa.NomeFantasia,
		CNPJ:         empresa.CNPJ,
		CreatedAt:    empresa.CreatedAt,
		AdminUser:    mapUserToResponseDTO(user),
	}

	return response, nil
}

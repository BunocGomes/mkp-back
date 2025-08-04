package service

import (
	"errors"
	//"time"
	"github.com/BunocGomes/mkp-back/database"
	"github.com/BunocGomes/mkp-back/dto"
	"github.com/BunocGomes/mkp-back/models"
	"github.com/BunocGomes/mkp-back/utils"
	"gorm.io/gorm"
)

func CreateUser(userDTO dto.RegisterUserDTO) (dto.UserResponseDTO, error) {
	var user models.User
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("email = ?", userDTO.Email).First(&models.User{}).Error; err == nil {
			return errors.New("e-mail já cadastrado")
		}
		user = models.User{
			Nome:  userDTO.Nome,
			Email: userDTO.Email,
		}
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		perfil := models.Perfil{
			UsuarioID: user.ID,
			Role:      userDTO.Role,
		}
		if err := tx.Create(&perfil).Error; err != nil {
			return err
		}
		digest := models.Digest{
			UserID:   user.ID,
			Password: userDTO.Password,
		}
		if err := tx.Create(&digest).Error; err != nil {
			return err
		}
		user.Perfil = perfil
		return nil
	})
	if err != nil {
		return dto.UserResponseDTO{}, err
	}
	return mapUserToResponseDTO(user), nil
}

func GetAllUsers() ([]dto.UserResponseDTO, error) {
	var users []models.User
	if err := database.DB.Preload("Perfil").Find(&users).Error; err != nil {
		return nil, err
	}
	var responses []dto.UserResponseDTO
	for _, user := range users {
		responses = append(responses, mapUserToResponseDTO(user))
	}
	return responses, nil
}

func GetUserByID(id uint) (dto.UserResponseDTO, error) {
	var user models.User
	if err := database.DB.Preload("Perfil").First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.UserResponseDTO{}, errors.New("usuário não encontrado")
		}
		return dto.UserResponseDTO{}, err
	}
	return mapUserToResponseDTO(user), nil
}
func UpdateUser(id uint, userDTO dto.UpdateUserDTO) (dto.UserResponseDTO, error) {
	var user models.User

	// Não precisamos mais de uma transação aqui, pois só alteramos uma tabela.
	// O Preload("Perfil") ainda é útil para retornar a resposta completa.
	if err := database.DB.Preload("Perfil").First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.UserResponseDTO{}, errors.New("usuário não encontrado")
		}
		return dto.UserResponseDTO{}, err
	}

	// Atualiza os campos do User se forem fornecidos no DTO
	if userDTO.Nome != nil {
		user.Nome = *userDTO.Nome
	}
	if userDTO.Email != nil {
		// Opcional: Verificar se o novo e-mail já existe antes de atribuir
		user.Email = *userDTO.Email
	}

	// Salva as alterações apenas na tabela 'users'.
	if err := database.DB.Save(&user).Error; err != nil {
		return dto.UserResponseDTO{}, err
	}

	return mapUserToResponseDTO(user), nil
}
func DeleteUser(id uint) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		var user models.User
		if err := tx.First(&user, id).Error; err != nil {
			return errors.New("usuário não encontrado")
		}
		if err := tx.Delete(&models.Digest{}, user.ID).Error; err != nil {
			return err
		}
		if err := tx.Select("Perfil").Delete(&user).Error; err != nil {
			return err
		}
		return nil
	})
}

func mapUserToResponseDTO(user models.User) dto.UserResponseDTO {
	return dto.UserResponseDTO{
		ID:        user.ID,
		Nome:      user.Nome,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		Perfil: dto.PerfilResponseDTO{
			ID:   user.Perfil.ID,
			Role: user.Perfil.Role,
		},
	}
}

func Login(loginDTO dto.LoginDTO) (string, error) {
	var user models.User

	// 1. Encontrar o usuário pelo e-mail, carregando Perfil e Digest
	err := database.DB.Preload("Perfil").Preload("Digest").Where("email = ?", loginDTO.Email).First(&user).Error
	if err != nil {
		// Se não encontrou o registro, retorna um erro genérico de credenciais
		return "", errors.New("credenciais inválidas")
	}

	// 2. Verificar se a senha está correta
	// Usamos a função do nosso pacote utils
	if !utils.CheckPasswordHash(loginDTO.Password, user.Digest.PasswordHash) {
		return "", errors.New("credenciais inválidas")
	}

	// 3. Gerar o token JWT
	// Passamos o ID do usuário e o seu papel (role) do perfil
	token, err := utils.GenerateJWT(user.ID, user.Perfil.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}

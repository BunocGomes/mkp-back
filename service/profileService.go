package service

import (
	"errors"
	"log"

	"github.com/BunocGomes/mkp-back/database"
	"github.com/BunocGomes/mkp-back/dto"
	"github.com/BunocGomes/mkp-back/models"
	"gorm.io/gorm"
)

// GetProfileByUserID busca um perfil completo pelo ID do usuário.
func GetProfileByUserID(userID uint) (dto.ProfileResponseDTO, error) {
	var user models.User
	err := database.DB.
		Preload("Perfil.Skills").
		Preload("Perfil.PortfolioItems").
		Preload("Perfil.SocialLinks").
		First(&user, userID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ProfileResponseDTO{}, errors.New("perfil de usuário não encontrado")
		}
		return dto.ProfileResponseDTO{}, err
	}

	// Log de depuração para ver o que o GORM carregou
	log.Printf("Dados carregados do banco: User.ID=%d, Perfil.ID=%d, Skills=%d, PortfolioItems=%d, SocialLinks=%d",
		user.ID, user.Perfil.ID, len(user.Perfil.Skills), len(user.Perfil.PortfolioItems), len(user.Perfil.SocialLinks))

	return mapProfileToResponseDTO(user), nil
}

// UpdateProfile atualiza o perfil de um usuário dentro de uma transação.
func UpdateProfile(userID uint, profileDTO dto.UpdateProfileDTO) (dto.ProfileResponseDTO, error) {
	log.Printf("Iniciando atualização para o usuário ID: %d", userID)
	log.Printf("Dados recebidos (DTO): %+v", profileDTO)

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		var perfil models.Perfil
		if err := tx.Where("usuario_id = ?", userID).First(&perfil).Error; err != nil {
			return errors.New("perfil não encontrado")
		}

		if profileDTO.Titulo != nil {
			perfil.Titulo = *profileDTO.Titulo
		}
		if profileDTO.Bio != nil {
			perfil.Bio = *profileDTO.Bio
		}
		if profileDTO.AvatarURL != nil {
			perfil.AvatarURL = *profileDTO.AvatarURL
		}
		if err := tx.Save(&perfil).Error; err != nil {
			return err
		}

		if profileDTO.Skills != nil {
			var skills []*models.Skill
			for _, skillDTO := range profileDTO.Skills {
				skill := models.Skill{}
				if err := tx.Where("nome = ?", skillDTO.Nome).FirstOrCreate(&skill, models.Skill{Nome: skillDTO.Nome}).Error; err != nil {
					return err
				}
				skills = append(skills, &skill)
			}
			if err := tx.Model(&perfil).Association("Skills").Replace(skills); err != nil {
				return err
			}
		}
		if profileDTO.PortfolioItems != nil {
			log.Println("Processando PortfolioItems...")
			if err := tx.Where("perfil_id = ?", perfil.ID).Delete(&models.PortfolioItem{}).Error; err != nil {
				log.Printf("ERRO ao deletar PortfolioItems antigos: %v", err)
				return err
			}
			log.Println("PortfolioItems antigos deletados com sucesso.")
			var portfolioItems []models.PortfolioItem
			for _, itemDTO := range profileDTO.PortfolioItems {
				portfolioItems = append(portfolioItems, models.PortfolioItem{
					PerfilID:   perfil.ID,
					Titulo:     itemDTO.Titulo,
					Descricao:  itemDTO.Descricao,
					ImagemURL:  itemDTO.ImagemURL,
					ProjetoURL: itemDTO.ProjetoURL,
				})
			}
			if len(portfolioItems) > 0 {
				log.Printf("Criando %d novos PortfolioItems...", len(portfolioItems))
				if err := tx.Create(&portfolioItems).Error; err != nil {
					log.Printf("ERRO ao criar novos PortfolioItems: %v", err)
					return err
				}
				log.Println("Novos PortfolioItems criados com sucesso.")
			}
		} else {
			log.Println("Nenhum PortfolioItem para processar (campo nulo no DTO).")
		}
		if profileDTO.SocialLinks != nil {
			log.Println("Processando SocialLinks...")
			if err := tx.Where("perfil_id = ?", perfil.ID).Delete(&models.SocialLink{}).Error; err != nil {
				log.Printf("ERRO ao deletar SocialLinks antigos: %v", err)
				return err
			}
			log.Println("SocialLinks antigos deletados com sucesso.")
			var socialLinks []models.SocialLink
			for _, linkDTO := range profileDTO.SocialLinks {
				socialLinks = append(socialLinks, models.SocialLink{
					PerfilID:   perfil.ID,
					Plataforma: linkDTO.Plataforma,
					URL:        linkDTO.URL,
				})
			}
			if len(socialLinks) > 0 {
				log.Printf("Criando %d novos SocialLinks...", len(socialLinks))
				if err := tx.Create(&socialLinks).Error; err != nil {
					log.Printf("ERRO ao criar novos SocialLinks: %v", err)
					return err
				}
				log.Println("Novos SocialLinks criados com sucesso.")
			}
		} else {
			log.Println("Nenhum SocialLink para processar (campo nulo no DTO).")
		}
		return nil
	})

	if err != nil {
		log.Printf("ERRO na transação de UpdateProfile: %v", err)
		return dto.ProfileResponseDTO{}, err
	}
	log.Println("Transação de UpdateProfile concluída com sucesso. Buscando perfil atualizado...")
	return GetProfileByUserID(userID)
}

// mapProfileToResponseDTO é uma função auxiliar para mapear os modelos para a resposta.
func mapProfileToResponseDTO(user models.User) dto.ProfileResponseDTO {
	var skillsDTO []dto.SkillDTO
	for _, s := range user.Perfil.Skills {
		skillsDTO = append(skillsDTO, dto.SkillDTO{Nome: s.Nome})
	}
	var portfolioItemsDTO []dto.PortfolioItemDTO
	for _, p := range user.Perfil.PortfolioItems {
		portfolioItemsDTO = append(portfolioItemsDTO, dto.PortfolioItemDTO{
			Titulo:     p.Titulo,
			Descricao:  p.Descricao,
			ImagemURL:  p.ImagemURL,
			ProjetoURL: p.ProjetoURL,
		})
	}
	var socialLinksDTO []dto.SocialLinkDTO
	for _, l := range user.Perfil.SocialLinks {
		socialLinksDTO = append(socialLinksDTO, dto.SocialLinkDTO{
			Plataforma: l.Plataforma,
			URL:        l.URL,
		})
	}
	return dto.ProfileResponseDTO{
		UsuarioID:      user.ID,
		NomeUsuario:    user.Nome,
		EmailUsuario:   user.Email,
		Titulo:         user.Perfil.Titulo,
		Bio:            user.Perfil.Bio,
		AvatarURL:      user.Perfil.AvatarURL,
		Skills:         skillsDTO,
		PortfolioItems: portfolioItemsDTO,
		SocialLinks:    socialLinksDTO,
	}
}

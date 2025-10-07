package models

import "gorm.io/gorm"

type Perfil struct {
	gorm.Model
	UsuarioID      uint            `json:"usuario_id"`
	Role           string          `gorm:"type:varchar(50);not null" json:"role"`
	Titulo         string          `gorm:"size:100" json:"titulo"`
	Bio            string          `gorm:"type:text" json:"bio"`
	AvatarURL      string          `json:"avatar_url"`
	Skills         []*Skill        `gorm:"many2many:perfil_skills;" json:"skills"`
	PortfolioItems []PortfolioItem `gorm:"foreignKey:PerfilID" json:"portfolio_items"`
	SocialLinks    []SocialLink    `gorm:"foreignKey:PerfilID" json:"social_links"`
}

type Skill struct {
	gorm.Model
	Nome string `gorm:"size:50;unique;not null" json:"nome"`
}

type PortfolioItem struct {
	gorm.Model
	PerfilID   uint   `json:"perfil_id"`
	Titulo     string `gorm:"size:100;not null" json:"titulo"`
	Descricao  string `json:"descricao"`
	ImagemURL  string `json:"imagem_url"`
	ProjetoURL string `json:"projeto_url"`
}

type SocialLink struct {
	gorm.Model
	PerfilID   uint   `json:"perfil_id"`
	Plataforma string `gorm:"size:50;not null" json:"plataforma"`
	URL        string `gorm:"not null" json:"url"`
}

package models

import "gorm.io/gorm"

// Perfil agora contém todas as informações da vitrine do usuário.
type Perfil struct {
	gorm.Model
	UsuarioID      uint            `json:"usuario_id"`
	Role           string          `gorm:"type:varchar(50);not null" json:"role"`
	Titulo         string          `gorm:"size:100" json:"titulo"` // Ex: "Desenvolvedor Go Sênior"
	Bio            string          `gorm:"type:text" json:"bio"`
	AvatarURL      string          `json:"avatar_url"`
	Skills         []*Skill        `gorm:"many2many:perfil_skills;" json:"skills"`
	PortfolioItems []PortfolioItem `gorm:"foreignKey:PerfilID" json:"portfolio_items"`
	SocialLinks    []SocialLink    `gorm:"foreignKey:PerfilID" json:"social_links"`
}

// Skill representa uma habilidade (ex: Go, Docker).
// A relação com Perfil é ManyToMany, pois um perfil pode ter várias skills
// e uma skill pode estar em vários perfis.
type Skill struct {
	gorm.Model
	Nome string `gorm:"size:50;unique;not null" json:"nome"`
}

// PortfolioItem é um item do portfólio do usuário.
// A relação é OneToMany (um Perfil tem muitos PortfolioItems).
type PortfolioItem struct {
	gorm.Model
	PerfilID   uint   `json:"perfil_id"`
	Titulo     string `gorm:"size:100;not null" json:"titulo"`
	Descricao  string `json:"descricao"`
	ImagemURL  string `json:"imagem_url"`
	ProjetoURL string `json:"projeto_url"`
}

// SocialLink é um link para uma rede social.
// A relação é OneToMany (um Perfil tem muitos SocialLinks).
type SocialLink struct {
	gorm.Model
	PerfilID   uint   `json:"perfil_id"`
	Plataforma string `gorm:"size:50;not null" json:"plataforma"` // "GitHub", "LinkedIn", etc.
	URL        string `gorm:"not null" json:"url"`
}

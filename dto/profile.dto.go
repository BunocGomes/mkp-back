package dto

// DTOs para os sub-recursos do perfil
type SkillDTO struct {
	Nome string `json:"nome" binding:"required"`
}
type PortfolioItemDTO struct {
	Titulo     string `json:"titulo" binding:"required"`
	Descricao  string `json:"descricao"`
	ImagemURL  string `json:"imagem_url"`
	ProjetoURL string `json:"projeto_url"`
}
type SocialLinkDTO struct {
	Plataforma string `json:"plataforma" binding:"required"`
	URL        string `json:"url" binding:"required,url"`
}

// UpdateProfileDTO é usado para atualizar o perfil do usuário.
type UpdateProfileDTO struct {
	Titulo         *string            `json:"titulo"`
	Bio            *string            `json:"bio"`
	AvatarURL      *string            `json:"avatar_url"`
	Skills         []SkillDTO         `json:"skills"`
	PortfolioItems []PortfolioItemDTO `json:"portfolio_items"`
	SocialLinks    []SocialLinkDTO    `json:"social_links"`
}

// ProfileResponseDTO é a resposta completa do perfil.
type ProfileResponseDTO struct {
	UsuarioID      uint               `json:"usuario_id"`
	NomeUsuario    string             `json:"nome_usuario"`
	EmailUsuario   string             `json:"email_usuario"`
	Titulo         string             `json:"titulo"`
	Bio            string             `json:"bio"`
	AvatarURL      string             `json:"avatar_url"`
	Skills         []SkillDTO         `json:"skills"`
	PortfolioItems []PortfolioItemDTO `json:"portfolio_items"`
	SocialLinks    []SocialLinkDTO    `json:"social_links"`
}

package dto

import "time"

// CreatePropostaDTO define a estrutura para um freelancer criar uma nova proposta.
type CreatePropostaDTO struct {
	ProjetoID     uint    `json:"projeto_id" binding:"required"`
	ValorProposto float64 `json:"valor_proposto" binding:"required"`
	Descricao     string  `json:"descricao" binding:"required"`
	PrazoEstimado int     `json:"prazo_estimado" binding:"required"`
}

// PropostaResponseDTO define a estrutura da resposta da API para uma proposta.
type PropostaResponseDTO struct {
	ID            uint      `json:"id"`
	ValorProposto float64   `json:"valor_proposto"`
	Descricao     string    `json:"descricao"`
	PrazoEstimado int       `json:"prazo_estimado"`
	Status        string    `json:"status"`
	ProjetoID     uint      `json:"projeto_id"`
	CreatedAt     time.Time `json:"created_at"`

	// Incluímos informações básicas do freelancer para facilitar a vida do frontend
	Freelancer struct {
		ID   uint   `json:"id"`
		Nome string `json:"nome"`
	} `json:"freelancer"`
}

type UpdatePropostaDTO struct {
	ValorProposto *float64 `json:"valor_proposto"`
	Descricao     *string  `json:"descricao"`
	PrazoEstimado *int     `json:"prazo_estimado"`
}

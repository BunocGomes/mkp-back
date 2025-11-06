package dto

import "time"

// CreateAvaliacaoDTO é o que o usuário envia para criar uma avaliação.
type CreateAvaliacaoDTO struct {
	ContratoID uint   `json:"contrato_id" binding:"required"`
	Nota       int    `json:"nota" binding:"required,min=1,max=5"`
	Comentario string `json:"comentario"`
}

// AvaliacaoResponseDTO é o que a API retorna ao listar avaliações.
type AvaliacaoResponseDTO struct {
	ID         uint      `json:"id"`
	Nota       int       `json:"nota"`
	Comentario string    `json:"comentario"`
	ContratoID uint      `json:"contrato_id"`
	CreatedAt  time.Time `json:"created_at"`

	Avaliador struct {
		ID   uint   `json:"id"`
		Nome string `json:"nome"`
	} `json:"avaliador"`
}

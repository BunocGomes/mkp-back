package dto

import "time"

// CreateProjectDTO define a estrutura para criar um novo projeto.
type CreateProjectDTO struct {
	Titulo    string  `json:"titulo" binding:"required"`
	Descricao string  `json:"descricao" binding:"required"`
	Orcamento float64 `json:"orcamento" binding:"required"`
}

// UpdateProjectDTO define a estrutura para atualizar um projeto existente.
type UpdateProjectDTO struct {
	Titulo    string  `json:"titulo,omitempty"`
	Descricao string  `json:"descricao,omitempty"`
	Status    string  `json:"status,omitempty"` // Permitir atualização do status
	Orcamento float64 `json:"orcamento,omitempty"`
}

// ProjectResponseDTO define a estrutura da resposta da API para um projeto.
type ProjectResponseDTO struct {
	ID        uint      `json:"id"`
	Titulo    string    `json:"titulo"`
	Descricao string    `json:"descricao"`
	Status    string    `json:"status"`
	Orcamento float64   `json:"orcamento"`
	EmpresaID uint      `json:"empresa_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

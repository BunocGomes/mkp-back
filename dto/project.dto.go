package dto

import "time"

type CreateProjectDTO struct {
	Titulo    string  `json:"titulo" binding:"required"`
	Descricao string  `json:"descricao" binding:"required"`
	Orcamento float64 `json:"orcamento" binding:"required"`
}

type UpdateProjectDTO struct {
	Titulo    string  `json:"titulo,omitempty"`
	Descricao string  `json:"descricao,omitempty"`
	Status    string  `json:"status,omitempty"`
	Orcamento float64 `json:"orcamento,omitempty"`
}

type ProjectResponseDTO struct {
	ID          uint      `json:"id"`
	Titulo      string    `json:"titulo"`
	Descricao   string    `json:"descricao"`
	Status      string    `json:"status"`
	Orcamento   float64   `json:"orcamento"`
	NomeEmpresa string    `json:"nome_empresa"`
	EmpresaID   uint      `json:"empresa_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

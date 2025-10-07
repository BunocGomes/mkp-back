package dto

import "time"

// ContratoResponseDTO define a estrutura da resposta da API para um contrato.
type ContratoResponseDTO struct {
	ID              uint      `json:"id"`
	ValorFinal      float64   `json:"valor_final"`
	DataInicio      time.Time `json:"data_inicio"`
	DataFimPrevista time.Time `json:"data_fim_prevista"`
	Status          string    `json:"status"`
	ProjetoID       uint      `json:"projeto_id"`

	ProjetoTitulo  string `json:"projeto_titulo"`  // Campo extra para conveniência
	EmpresaNome    string `json:"empresa_nome"`    // Campo extra para conveniência
	FreelancerNome string `json:"freelancer_nome"` // Campo extra para conveniência
}

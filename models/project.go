package models

import "gorm.io/gorm"

// Projeto representa um trabalho ou tarefa postada por uma empresa.
type Projeto struct {
	gorm.Model
	Titulo    string  `gorm:"size:255;not null"`
	Descricao string  `gorm:"type:text;not null"`
	Status    string  `gorm:"size:50;default:'aberto'"` // Ex: aberto, em_andamento, concluido, cancelado
	Orcamento float64 `gorm:"type:decimal(10,2)"`

	// Chave Estrangeira para a Empresa que criou o projeto
	EmpresaID uint    `gorm:"not null"`
	Empresa   Empresa `gorm:"foreignKey:EmpresaID"`
}

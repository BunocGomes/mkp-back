package models

import "gorm.io/gorm"

// Proposta representa a oferta de um freelancer para um projeto específico.
type Proposta struct {
	gorm.Model
	ValorProposto float64 `gorm:"type:decimal(10,2);not null"`
	Descricao     string  `gorm:"type:text;not null"`        // "Cover Letter" do freelancer
	PrazoEstimado int     `gorm:"not null"`                  // Prazo em dias
	Status        string  `gorm:"size:50;default:'enviada'"` // Ex: enviada, visualizada, aceita, recusada

	// Chave Estrangeira para o Projeto
	ProjetoID uint    `gorm:"not null"`
	Projeto   Projeto `gorm:"foreignKey:ProjetoID"`

	// Chave Estrangeira para o Freelancer (que é um User)
	FreelancerID uint `gorm:"not null"`
	Freelancer   User `gorm:"foreignKey:FreelancerID"`
}

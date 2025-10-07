package models

import "gorm.io/gorm"

type Proposta struct {
	gorm.Model
	ValorProposto float64 `gorm:"type:decimal(10,2);not null"`
	Descricao     string  `gorm:"type:text;not null"`
	PrazoEstimado int     `gorm:"not null"`
	Status        string  `gorm:"size:50;default:'enviada'"`

	ProjetoID uint    `gorm:"not null"`
	Projeto   Projeto `gorm:"foreignKey:ProjetoID"`

	FreelancerID uint `gorm:"not null"`
	Freelancer   User `gorm:"foreignKey:FreelancerID"`
}

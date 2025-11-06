package models

import (
	"time"

	"gorm.io/gorm"
)

type Contrato struct {
	gorm.Model
	ValorFinal            float64   `gorm:"type:decimal(10,2);not null"`
	DataInicio            time.Time `gorm:"not null"`
	DataFimPrevista       time.Time `gorm:"not null"`
	Status                string    `gorm:"size:50;default:'ativo'"`
	AvaliacaoEmpresaID    *uint     
	AvaliacaoFreelancerID *uint     

	ProjetoID uint    `gorm:"not null"`
	Projeto   Projeto `gorm:"foreignKey:ProjetoID"`

	EmpresaID uint    `gorm:"not null"`
	Empresa   Empresa `gorm:"foreignKey:EmpresaID"`

	FreelancerID uint `gorm:"not null"`
	Freelancer   User `gorm:"foreignKey:FreelancerID"`

	PropostaID uint     `gorm:"unique;not null"`
	Proposta   Proposta `gorm:"foreignKey:PropostaID"`
}

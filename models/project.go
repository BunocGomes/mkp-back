package models

import "gorm.io/gorm"

type Projeto struct {
	gorm.Model
	Titulo    string  `gorm:"size:255;not null"`
	Descricao string  `gorm:"type:text;not null"`
	Status    string  `gorm:"size:50;default:'aberto'"`
	Orcamento float64 `gorm:"type:decimal(10,2)"`

	EmpresaID uint    `gorm:"not null"`
	Empresa   Empresa `gorm:"foreignKey:EmpresaID"`
}

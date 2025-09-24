package models

import "gorm.io/gorm"

// Empresa representa a entidade de uma companhia na plataforma.
type Empresa struct {
	gorm.Model
	NomeFantasia string `gorm:"size:255;not null" json:"nome_fantasia"`
	RazaoSocial  string `gorm:"size:255;not null" json:"razao_social"`
	CNPJ         string `gorm:"size:18;unique;not null" json:"cnpj"`

	Usuarios []User `gorm:"foreignKey:EmpresaID" json:"usuarios"`
}

package models

import "gorm.io/gorm"

// Avaliacao armazena a nota e o comentário de um usuário para outro
type Avaliacao struct {
	gorm.Model
	Nota      int    `gorm:"not null"` // Nota de 1 a 5
	Comentario string `gorm:"type:text"`

	ContratoID      uint `gorm:"not null"`
	AvaliadorUserID uint `gorm:"not null"` // Quem escreveu a avaliação
	AvaliadoUserID  uint `gorm:"not null"` // Quem recebeu a avaliação

	Contrato      Contrato
	AvaliadorUser User `gorm:"foreignKey:AvaliadorUserID"`
	AvaliadoUser  User `gorm:"foreignKey:AvaliadoUserID"`
}
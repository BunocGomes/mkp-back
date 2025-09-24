package models

import (
	"github.com/BunocGomes/mkp-back/utils"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Nome      string `json:"nome"`
	Email     string `gorm:"unique" json:"email"`
	EmpresaID *uint  `gorm:"index" json:"empresa_id"`
	Perfil    Perfil `gorm:"foreignKey:UsuarioID" json:"perfil"`
	Digest    Digest `gorm:"foreignKey:UserID" json:"digest"`
}

type Digest struct {
	UserID       uint   `gorm:"primaryKey"`
	Password     string `gorm:"-" json:"-"`
	PasswordHash string `gorm:"not null" json:"-"`
}

func (d *Digest) BeforeSave(tx *gorm.DB) (err error) {
	if d.Password != "" {
		hashedPassword, err := utils.HashPassword(d.Password)
		if err != nil {
			return err
		}
		d.PasswordHash = hashedPassword
	}
	return
}

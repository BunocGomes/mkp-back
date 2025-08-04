package dto

import "time"

type RegisterUserDTO struct {
	Nome     string `json:"nome" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required"`
}

type UpdateUserDTO struct {
	Nome  *string `json:"nome"`
	Email *string `json:"email"`
}

type PerfilResponseDTO struct {
	ID   uint   `json:"id"`
	Role string `json:"role"`
}

type UserResponseDTO struct {
	ID        uint              `json:"id"`
	Nome      string            `json:"nome"`
	Email     string            `json:"email"`
	CreatedAt time.Time         `json:"created_at"`
	Perfil    PerfilResponseDTO `json:"perfil"`
}

type LoginDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponseDTO struct {
	Token string `json:"token"`
}

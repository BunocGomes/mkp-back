package dto

import "time"

// RegisterEmpresaDTO é o DTO para o request de criação de uma conta de empresa.
type RegisterEmpresaDTO struct {
	// Dados do Usuário Admin
	NomeUsuario     string `json:"nome_usuario" binding:"required"`
	EmailUsuario    string `json:"email_usuario" binding:"required,email"`
	PasswordUsuario string `json:"password_usuario" binding:"required,min=6"`

	// Dados da Empresa
	NomeFantasia string `json:"nome_fantasia" binding:"required"`
	RazaoSocial  string `json:"razao_social" binding:"required"`
	CNPJ         string `json:"cnpj" binding:"required"`
}

// EmpresaResponseDTO é a resposta após a criação bem-sucedida.
type EmpresaResponseDTO struct {
	ID           uint            `json:"id"`
	NomeFantasia string          `json:"nome_fantasia"`
	CNPJ         string          `json:"cnpj"`
	CreatedAt    time.Time       `json:"created_at"`
	AdminUser    UserResponseDTO `json:"admin_user"`
}

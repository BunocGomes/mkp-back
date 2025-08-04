package helper

import (
	"errors"

	"github.com/BunocGomes/mkp-back/utils"
)

func IsSuperAdmin(tokenString string) (bool, error) {
	claims, err := utils.VerifyJWT(tokenString)
	if err != nil {
		return false, errors.New("token inválido")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return false, errors.New("não foi possível identificar o papel do usuário no token")
	}

	return role == "superadmin", nil
}

func IsEmpresaOrHigher(tokenString string) (bool, error) {
	claims, err := utils.VerifyJWT(tokenString)
	if err != nil {
		return false, errors.New("token inválido")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return false, errors.New("não foi possível identificar o papel do usuário no token")
	}

	return role == "empresa" || role == "superadmin", nil
}

func IsFreeLancerOrHigher(tokenString string) (bool, error) {
	claims, err := utils.VerifyJWT(tokenString)
	if err != nil {
		return false, errors.New("token inválido")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return false, errors.New("não foi possível identificar o papel do usuário no token")
	}

	return role == "freelancer" || role == "superadmin", nil
}

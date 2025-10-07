package utils

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtSecret []byte

func init() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("A variável de ambiente JWT_SECRET não foi definida.")
	}
	jwtSecret = []byte(secret)
}

func GenerateJWT(userID uint, role string, empresaID *uint) (string, error) {
	expirationTime := time.Now().Add(365 * 24 * time.Hour)

	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     expirationTime.Unix(),
	}

	if empresaID != nil {
		claims["empresa_id"] = *empresaID
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

var ErrInvalidToken = errors.New("invalid token")

func VerifyJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

func GetIDFromToken(tokenString string) (uint, error) {
	claims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil ||
		!token.Valid {
		return 0, errors.New("token inválido")
	}

	userID, ok := (*claims)["user_id"].(float64)
	if !ok {
		return 0, errors.New("token inválido")
	}

	return uint(userID), nil
}

func GetUserFromToken(tokenString string) (uint, string, error) {
	claims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return 0, "", errors.New("token inválido")
	}

	userID, ok := (*claims)["user_id"].(float64)
	if !ok {
		return 0, "", errors.New("token inválido")
	}

	role, ok := (*claims)["role"].(string)
	if !ok {
		return 0, "", errors.New("papel do usuário não encontrado")
	}

	return uint(userID), role, nil
}

func GetRoleFromToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		role, ok := claims["role"].(string)
		if !ok {
			return "", errors.New("role not found in token")
		}
		return role, nil
	}

	return "", errors.New("invalid token")
}

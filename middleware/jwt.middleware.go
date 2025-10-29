package middleware

import (
	"net/http"
	"strings"

	"log"

	"github.com/BunocGomes/mkp-back/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// 1. Tenta pegar o token do Header "Authorization" (para rotas HTTP normais)
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		}

		// 2. Se não achou no header, tenta pegar do Query Param "token" (para WebSockets)
		if tokenString == "" {
			tokenString = c.Query("token")
		}

		// 3. Se não achou em lugar nenhum, bloqueia
		if tokenString == "" {
			log.Println("AuthMiddleware: Token não encontrado (nem no header, nem na query)")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token de autorização não fornecido"})
			return
		}

		// 4. Se chegou aqui, temos um token. Vamos validá-lo.
		claims, err := utils.VerifyJWT(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			return
		}

		// 5. Extrai os dados do token e coloca no contexto
		userID, ok := claims["user_id"].(float64)
		role, ok2 := claims["role"].(string)
		if !ok || !ok2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token com claims inválidas"})
			return
		}

		c.Set("userId", uint(userID))
		c.Set("role", role)

		// Verifica se a claim "empresa_id" existe no token e a define no contexto
		if empresaIDClaim, exists := claims["empresa_id"]; exists {
			empresaID, ok := empresaIDClaim.(float64)
			if ok {
				c.Set("empresaId", uint(empresaID))
			}
		}

		c.Next()
	}
}

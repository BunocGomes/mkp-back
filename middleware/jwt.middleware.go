package middleware

import (
	"net/http"
	"strings"

	"github.com/BunocGomes/mkp-back/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Cabeçalho de autorização não fornecido"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Formato do token de autorização inválido"})
			return
		}

		claims, err := utils.VerifyJWT(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			return
		}

		// --- LÓGICA ATUALIZADA AQUI ---
		userID := uint(claims["user_id"].(float64))
		role := claims["role"].(string)

		c.Set("userId", userID)
		c.Set("role", role)

		// Verifica se a claim "empresa_id" existe no token e a define no contexto
		if empresaIDClaim, exists := claims["empresa_id"]; exists {
			empresaID := uint(empresaIDClaim.(float64))
			c.Set("empresaId", empresaID) // O controller vai ler daqui!
		}
		// ---------------------------------

		c.Next()
	}
}

package helper

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func ExtractToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if len(authHeader) > 7 && strings.ToUpper(authHeader[:7]) == "BEARER " {
		return authHeader[7:]
	}
	return authHeader
}

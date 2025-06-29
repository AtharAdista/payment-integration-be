package middleware

import (
	"net/http"
	"payment/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(maker *utils.JWTMaker) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := maker.VerifyToken(token)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.Set("user", claims)
		c.Next()
	}
}

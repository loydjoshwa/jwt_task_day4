package middleware

import (
	"myapp/utils"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "missing token"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(401, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		token, err := utils.VerifyAccessToken(parts[1])
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.JSON(401, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// --- SAFETY CHECKS (Panic Fix) ---

		// 1. user_id check
		uid, ok := claims["user_id"].(float64)
		if !ok {
			c.JSON(401, gin.H{"error": "user_id not found in token"})
			c.Abort()
			return
		}

		// 2. role check
		role, ok := claims["role"].(string)
		if !ok {
			role = "user" // Default role set cheyyunnu
		}

		c.Set("user_id", uint(uid))
		c.Set("role", role)
		c.Next()
	}
}
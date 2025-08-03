package middleware

import (
	"net/http"
	"strings"

	"github.com/fisiopet/bp/internal/domain/auth"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware cria um middleware de autenticação JWT
func AuthMiddleware(jwtService auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extrai o token do header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Authorization header required",
				"message": "Token not provided",
			})
			c.Abort()
			return
		}

		// Verifica se o header tem o formato "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid authorization header format",
				"message": "Expected format: Bearer <token>",
			})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		// Valida o token
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			status := http.StatusUnauthorized
			message := "Invalid token"

			if err == auth.ErrExpiredToken {
				message = "Token expired"
			}

			c.JSON(status, gin.H{
				"error":   "Authentication failed",
				"message": message,
			})
			c.Abort()
			return
		}

		// Adiciona as informações do usuário ao contexto
		c.Set("userID", claims.UserID)
		c.Set("userEmail", claims.Email)
		c.Set("userRole", claims.Role)

		c.Next()
	}
}

// RoleMiddleware cria um middleware para verificar roles específicos
func RoleMiddleware(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("userRole")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "User role not found",
				"message": "Authentication required",
			})
			c.Abort()
			return
		}

		role := userRole.(string)
		hasPermission := false

		for _, requiredRole := range requiredRoles {
			if role == requiredRole {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Insufficient permissions",
				"message": "You don't have permission to access this resource",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

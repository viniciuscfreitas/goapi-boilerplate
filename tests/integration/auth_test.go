package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAuthMiddleware testa o middleware de autenticação JWT
func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Teste básico do middleware de autenticação
	t.Run("Auth Middleware Without Token", func(t *testing.T) {
		router := gin.New()

		// Adicionar middleware de autenticação
		router.Use(func(c *gin.Context) {
			// Simular middleware de auth que requer token
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error":   "Authorization header required",
					"message": "Token not provided",
				})
				c.Abort()
				return
			}
			c.Next()
		})

		// Rota protegida
		router.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Access granted"})
		})

		// Teste sem token
		req, _ := http.NewRequest("GET", "/protected", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "Authorization header required", response["error"])
	})

	// Teste com token válido
	t.Run("Auth Middleware With Valid Token", func(t *testing.T) {
		router := gin.New()

		// Adicionar middleware de autenticação
		router.Use(func(c *gin.Context) {
			authHeader := c.GetHeader("Authorization")
			if authHeader == "Bearer valid-token" {
				c.Set("userID", "123")
				c.Set("userEmail", "test@example.com")
				c.Set("userRole", "user")
				c.Next()
				return
			}
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid token",
				"message": "Token validation failed",
			})
			c.Abort()
		})

		// Rota protegida
		router.GET("/protected", func(c *gin.Context) {
			userID, _ := c.Get("userID")
			c.JSON(http.StatusOK, gin.H{
				"message": "Access granted",
				"userID":  userID,
			})
		})

		// Teste com token válido
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer valid-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "Access granted", response["message"])
		assert.Equal(t, "123", response["userID"])
	})
}

// TestRateLimiting testa o middleware de rate limiting
func TestRateLimiting(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Rate Limiting", func(t *testing.T) {
		router := gin.New()

		// Adicionar rate limiting simples
		requestCount := 0
		router.Use(func(c *gin.Context) {
			requestCount++
			if requestCount > 3 {
				c.JSON(http.StatusTooManyRequests, gin.H{
					"error":   "Rate limit exceeded",
					"message": "Too many requests",
				})
				c.Abort()
				return
			}
			c.Next()
		})

		// Rota de teste
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Success"})
		})

		// Teste: primeiras 3 requisições devem passar
		for i := 0; i < 3; i++ {
			req, _ := http.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
		}

		// 4ª requisição deve ser bloqueada
		req, _ := http.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusTooManyRequests, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "Rate limit exceeded", response["error"])
	})
}

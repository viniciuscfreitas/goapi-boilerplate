package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const requestIDKey = "requestID"

// Logger cria um middleware de logging para Gin
func Logger(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Gera um ID único para cada requisição
		requestID := uuid.New().String()
		c.Set(requestIDKey, requestID) // Adiciona o ID ao contexto do Gin

		// Cria um logger filho com o contexto da requisição
		reqLog := log.With(
			"request_id", requestID,
		)
		
		// Processa a requisição
		c.Next()

		// Quando a requisição termina, loga as informações
		latency := time.Since(start)

		reqLog.Info("Request handled",
			"status_code", c.Writer.Status(),
			"method", c.Request.Method,
			"path", path,
			"query", query,
			"ip_address", c.ClientIP(),
			"latency_ms", float64(latency.Milliseconds()),
			"user_agent", c.Request.UserAgent(),
		)
	}
}

// GetRequestID retorna o ID da requisição do contexto
func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get(requestIDKey); exists {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
} 
package router

import (
	"log/slog"

	"github.com/fisiopet/bp/internal/infrastructure/http/handlers"
	"github.com/fisiopet/bp/internal/infrastructure/http/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRouter configura as rotas da aplicação
func SetupRouter(userHandler *handlers.UserHandler, log *slog.Logger) *gin.Engine {
	router := gin.New() // Use gin.New() para ter mais controle sobre os middlewares

	// Middleware de logging (deve ser o primeiro)
	router.Use(middleware.Logger(log))
	
	// Middleware de recuperação de pânico
	router.Use(gin.Recovery())

	// Middleware de CORS
	router.Use(corsMiddleware())

	// Grupo de rotas da API
	api := router.Group("/api/v1")
	{
		// Rotas de autenticação
		auth := api.Group("/auth")
		{
			auth.POST("/login", userHandler.Login)
		}

		// Rotas de usuários
		users := api.Group("/users")
		{
			users.POST("", userHandler.CreateUser)
			users.GET("", userHandler.ListUsers)
			users.GET("/email", userHandler.GetUserByEmail)
			users.GET("/:id", userHandler.GetUserByID)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}
	}

	// Rota de health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"message": "API is running",
		})
	})

	return router
}

// corsMiddleware configura o CORS para a API
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
} 
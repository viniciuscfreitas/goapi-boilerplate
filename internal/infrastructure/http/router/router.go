package router

import (
	"log/slog"

	"go-api-boilerplate/internal/domain/auth"
	"go-api-boilerplate/internal/infrastructure/http/handlers"
	"go-api-boilerplate/internal/infrastructure/http/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter configura as rotas da aplicação
func SetupRouter(userHandler *handlers.UserHandler, jwtService auth.JWTService, log *slog.Logger) *gin.Engine {
	router := gin.New() // Use gin.New() para ter mais controle sobre os middlewares

	// Middleware de logging (deve ser o primeiro)
	router.Use(middleware.Logger(log))

	// Middleware de recuperação de pânico
	router.Use(gin.Recovery())

	// Middleware de segurança
	securityConfig := middleware.SecurityConfig{
		CORSOrigins: []string{"*"}, // Em produção, especificar domínios específicos
		RateLimit:   100,           // 100 requests por segundo por IP
	}

	// Middleware de CORS seguro
	router.Use(middleware.CORSMiddleware(securityConfig))

	// Middleware de rate limiting
	router.Use(middleware.RateLimitMiddleware(securityConfig))

	// Middleware de headers de segurança
	router.Use(middleware.SecurityHeadersMiddleware())

	// Middleware de request ID para rastreabilidade
	router.Use(middleware.RequestIDMiddleware())

	// Grupo de rotas da API
	api := router.Group("/api/v1")
	{
		// Rotas de autenticação (públicas)
		auth := api.Group("/auth")
		{
			auth.POST("/login", userHandler.Login)
			auth.POST("/register", userHandler.CreateUser) // Endpoint público para registro
		}

		// Rotas de usuários (protegidas por autenticação)
		users := api.Group("/users")
		users.Use(middleware.AuthMiddleware(jwtService)) // Aplica autenticação em todas as rotas de usuários
		{
			// Rotas que requerem autenticação básica
			users.GET("", userHandler.ListUsers)
			users.GET("/email", userHandler.GetUserByEmail)
			users.GET("/:id", userHandler.GetUserByID)

			// Rotas que requerem role de admin
			adminRoutes := users.Group("")
			adminRoutes.Use(middleware.RoleMiddleware("admin"))
			{
				adminRoutes.POST("", userHandler.CreateUser)
				adminRoutes.PUT("/:id", userHandler.UpdateUser)
				adminRoutes.DELETE("/:id", userHandler.DeleteUser)
			}
		}
	}

	// Rota de health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "API is running",
		})
	})

	// Rota para o arquivo swagger.json (fora do grupo /swagger para evitar conflito)
	router.GET("/swagger.json", func(c *gin.Context) {
		c.File("./docs/swagger.json")
	})

	// Rota do Swagger UI com URL externa
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("http://localhost:8080/swagger.json")))

	return router
}

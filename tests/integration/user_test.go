package integration

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fisiopet/bp/internal/domain/auth"
	"github.com/fisiopet/bp/internal/domain/user"
	"github.com/fisiopet/bp/internal/infrastructure/http/handlers"
	"github.com/fisiopet/bp/internal/infrastructure/repository"
	"github.com/fisiopet/bp/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "github.com/lib/pq"
)

// setupTestDB cria uma conexão de teste com o banco
func setupTestDB(t *testing.T) *sql.DB {
	// Usar banco de teste separado
	dsn := "host=localhost port=5433 user=postgres password=secret dbname=boilerplate_test sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	require.NoError(t, err)
	
	// Testar conexão
	err = db.Ping()
	require.NoError(t, err)
	
	return db
}

// setupTestRouter cria um router de teste com dependências reais
func setupTestRouter(t *testing.T) *gin.Engine {
	db := setupTestDB(t)
	
	// Limpar banco antes de cada teste
	_, err := db.Exec("DELETE FROM users")
	require.NoError(t, err)
	
	// Inicializar dependências
	userRepo := repository.NewPostgresUserRepository(db)
	jwtService := auth.NewJWTService("test-secret", 24*time.Hour)
	userUseCase := usecase.NewUserUseCase(userRepo, jwtService)
	userHandler := handlers.NewUserHandler(userUseCase)
	
	// Configurar router básico para testes
	router := gin.New()
	router.Use(gin.Recovery())
	
	// Configurar rotas
	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", userHandler.Login)
		}
		
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
	
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"message": "API is running",
		})
	})
	
	return router
}

// TestUserCRUD testa o fluxo completo CRUD de usuários
func TestUserCRUD(t *testing.T) {
	router := setupTestRouter(t)
	
	// Teste 1: Criar usuário
	t.Run("Create User", func(t *testing.T) {
		userData := handlers.CreateUserRequest{
			Email:    "test@example.com",
			Password: "password123",
			Name:     "Test User",
			Role:     "user",
		}
		
		jsonData, _ := json.Marshal(userData)
		req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusCreated, w.Code)
		
		var response user.User
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		
		assert.NotEmpty(t, response.ID)
		assert.Equal(t, userData.Email, response.Email)
		assert.Equal(t, userData.Name, response.Name)
		assert.Equal(t, string(userData.Role), string(response.Role))
		assert.True(t, response.IsActive)
	})
	
	// Teste 2: Buscar usuário por email
	t.Run("Get User by Email", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/users/email?email=test@example.com", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response user.User
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		
		assert.Equal(t, "test@example.com", response.Email)
	})
	
	// Teste 3: Login
	t.Run("Login", func(t *testing.T) {
		loginData := handlers.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}
		
		jsonData, _ := json.Marshal(loginData)
		req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		
		// Verificar se retorna token
		assert.Contains(t, response, "token")
		assert.NotEmpty(t, response["token"])
	})
	
	// Teste 4: Listar usuários
	t.Run("List Users", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/users", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		
		assert.Contains(t, response, "users")
		assert.Contains(t, response, "total")
		assert.Greater(t, response["total"], float64(0))
	})
}

// TestUserValidation testa validações de entrada
func TestUserValidation(t *testing.T) {
	router := setupTestRouter(t)
	
	t.Run("Invalid Email", func(t *testing.T) {
		userData := handlers.CreateUserRequest{
			Email:    "invalid-email",
			Password: "password123",
			Name:     "Test User",
			Role:     "user",
		}
		
		jsonData, _ := json.Marshal(userData)
		req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	
	t.Run("Weak Password", func(t *testing.T) {
		userData := handlers.CreateUserRequest{
			Email:    "test@example.com",
			Password: "123", // Senha muito curta
			Name:     "Test User",
			Role:     "user",
		}
		
		jsonData, _ := json.Marshal(userData)
		req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	
	t.Run("Invalid Role", func(t *testing.T) {
		userData := handlers.CreateUserRequest{
			Email:    "test@example.com",
			Password: "password123",
			Name:     "Test User",
			Role:     "invalid_role",
		}
		
		jsonData, _ := json.Marshal(userData)
		req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// TestAuthentication testa cenários de autenticação
func TestAuthentication(t *testing.T) {
	router := setupTestRouter(t)
	
	// Criar usuário primeiro
	userData := handlers.CreateUserRequest{
		Email:    "auth@example.com",
		Password: "password123",
		Name:     "Auth User",
		Role:     "user",
	}
	
	jsonData, _ := json.Marshal(userData)
	req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
	
	t.Run("Valid Login", func(t *testing.T) {
		loginData := handlers.LoginRequest{
			Email:    "auth@example.com",
			Password: "password123",
		}
		
		jsonData, _ := json.Marshal(loginData)
		req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
	})
	
	t.Run("Invalid Password", func(t *testing.T) {
		loginData := handlers.LoginRequest{
			Email:    "auth@example.com",
			Password: "wrongpassword",
		}
		
		jsonData, _ := json.Marshal(loginData)
		req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
	
	t.Run("User Not Found", func(t *testing.T) {
		loginData := handlers.LoginRequest{
			Email:    "nonexistent@example.com",
			Password: "password123",
		}
		
		jsonData, _ := json.Marshal(loginData)
		req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
} 
package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/fisiopet/bp/internal/domain/user"
	"github.com/fisiopet/bp/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UserHandler implementa os handlers HTTP para usuários
type UserHandler struct {
	userUseCase *usecase.UserUseCase
}

// NewUserHandler cria uma nova instância de UserHandler
func NewUserHandler(userUseCase *usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

// CreateUserRequest representa a requisição de criação de usuário
type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

// UpdateUserRequest representa a requisição de atualização de usuário
type UpdateUserRequest struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty" binding:"omitempty,email"`
	Role  *string `json:"role,omitempty"`
}

// LoginRequest representa a requisição de login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// CreateUser cria um novo usuário
// @Summary Criar usuário
// @Description Cria um novo usuário no sistema
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "Dados do usuário"
// @Success 201 {object} user.User
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request data",
			Message: err.Error(),
		})
		return
	}

	// Validar role
	role, err := h.validateRole(req.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid role",
			Message: err.Error(),
		})
		return
	}

	// Criar usuário
	input := usecase.CreateUserInput{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
		Role:     role,
	}

	output, err := h.userUseCase.CreateUser(c.Request.Context(), input)
	if err != nil {
		status, message := h.mapErrorToHTTPStatus(err)
		c.JSON(status, ErrorResponse{
			Error:   "Failed to create user",
			Message: message,
		})
		return
	}

	c.JSON(http.StatusCreated, output.User)
}

// GetUserByID busca um usuário pelo ID
// @Summary Buscar usuário por ID
// @Description Busca um usuário específico pelo ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "ID do usuário"
// @Success 200 {object} user.User
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	// 1. Obtenha o ID da URL e valide-o
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid user ID",
			Message: "User ID is required",
		})
		return
	}

	// 2. Valide se é um UUID válido
	_, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid user ID",
			Message: "User ID must be a valid UUID",
		})
		return
	}

	// 3. Chame o caso de uso
	input := usecase.GetUserByIDInput{ID: idStr}
	output, err := h.userUseCase.GetUserByID(c.Request.Context(), input)

	// 4. ESTE É O BLOCO MAIS IMPORTANTE: Trate o erro PRIMEIRO
	if err != nil {
		// Usa a função centralizada para mapear o erro de domínio para um status HTTP
		status, message := h.mapErrorToHTTPStatus(err)
		c.JSON(status, ErrorResponse{
			Error:   "Failed to get user",
			Message: message,
		})
		return // Encerra a execução aqui!
	}

	// 5. Se não houve erro, retorne o sucesso
	c.JSON(http.StatusOK, output.User)
}

// GetUserByEmail busca um usuário pelo email
// @Summary Buscar usuário por email
// @Description Busca um usuário específico pelo email
// @Tags users
// @Accept json
// @Produce json
// @Param email query string true "Email do usuário"
// @Success 200 {object} user.User
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/email [get]
func (h *UserHandler) GetUserByEmail(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid email",
			Message: "Email is required",
		})
		return
	}

	input := usecase.GetUserByEmailInput{Email: email}
	output, err := h.userUseCase.GetUserByEmail(c.Request.Context(), input)
	if err != nil {
		status, message := h.mapErrorToHTTPStatus(err)
		c.JSON(status, ErrorResponse{
			Error:   "Failed to get user",
			Message: message,
		})
		return
	}

	c.JSON(http.StatusOK, output.User)
}

// UpdateUser atualiza um usuário existente
// @Summary Atualizar usuário
// @Description Atualiza os dados de um usuário existente
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "ID do usuário"
// @Param user body UpdateUserRequest true "Dados para atualização"
// @Success 200 {object} user.User
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	// 1. Obtenha o ID da URL e valide-o
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid user ID",
			Message: "User ID is required",
		})
		return
	}

	// 2. Valide se é um UUID válido
	_, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid user ID",
			Message: "User ID must be a valid UUID",
		})
		return
	}

	// 3. Decodifique o corpo da requisição JSON
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request data",
			Message: err.Error(),
		})
		return
	}

	// 4. Prepare o input para o caso de uso
	input := usecase.UpdateUserInput{ID: idStr}

	// 5. Mapear campos opcionais
	if req.Name != nil {
		input.Name = req.Name
	}
	if req.Email != nil {
		input.Email = req.Email
	}
	if req.Role != nil {
		role, err := h.validateRole(*req.Role)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Invalid role",
				Message: err.Error(),
			})
			return
		}
		input.Role = &role
	}

	// 6. Chame o caso de uso
	output, err := h.userUseCase.UpdateUser(c.Request.Context(), input)

	// 7. ESTE É O BLOCO MAIS IMPORTANTE: Trate o erro PRIMEIRO
	if err != nil {
		// Usa a função centralizada para mapear o erro de domínio para um status HTTP
		status, message := h.mapErrorToHTTPStatus(err)
		c.JSON(status, ErrorResponse{
			Error:   "Failed to update user",
			Message: message,
		})
		return // Encerra a execução aqui!
	}

	// 8. Se não houve erro, retorne o sucesso
	c.JSON(http.StatusOK, output.User)
}

// DeleteUser remove um usuário
// @Summary Deletar usuário
// @Description Remove um usuário do sistema
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "ID do usuário"
// @Success 204 "No Content"
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	// 1. Obtenha o ID da URL e valide-o
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid user ID",
			Message: "User ID is required",
		})
		return
	}

	// 2. Valide se é um UUID válido
	_, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid user ID",
			Message: "User ID must be a valid UUID",
		})
		return
	}

	// 3. Chame o caso de uso
	input := usecase.DeleteUserInput{ID: idStr}
	err = h.userUseCase.DeleteUser(c.Request.Context(), input)

	// 4. ESTE É O BLOCO MAIS IMPORTANTE: Trate o erro PRIMEIRO
	if err != nil {
		// Usa a função centralizada para mapear o erro de domínio para um status HTTP
		status, message := h.mapErrorToHTTPStatus(err)
		c.JSON(status, ErrorResponse{
			Error:   "Failed to delete user",
			Message: message,
		})
		return // Encerra a execução aqui!
	}

	// 5. Se não houve erro, retorne o sucesso
	c.Status(http.StatusNoContent)
}

// ListUsers lista usuários com paginação
// @Summary Listar usuários
// @Description Lista usuários com paginação
// @Tags users
// @Accept json
// @Produce json
// @Param offset query int false "Offset para paginação" default(0)
// @Param limit query int false "Limite de registros" default(10)
// @Success 200 {object} usecase.ListUsersOutput
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	offsetStr := c.DefaultQuery("offset", "0")
	limitStr := c.DefaultQuery("limit", "10")

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid offset",
			Message: "Offset must be a positive integer",
		})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid limit",
			Message: "Limit must be a positive integer",
		})
		return
	}

	input := usecase.ListUsersInput{
		Offset: offset,
		Limit:  limit,
	}

	output, err := h.userUseCase.ListUsers(c.Request.Context(), input)
	if err != nil {
		status, message := h.mapErrorToHTTPStatus(err)
		c.JSON(status, ErrorResponse{
			Error:   "Failed to list users",
			Message: message,
		})
		return
	}

	c.JSON(http.StatusOK, output)
}

// Login autentica um usuário
// @Summary Login
// @Description Autentica um usuário no sistema
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Credenciais de login"
// @Success 200 {object} user.User
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request data",
			Message: err.Error(),
		})
		return
	}

	input := usecase.AuthenticateUserInput{
		Email:    req.Email,
		Password: req.Password,
	}

	output, err := h.userUseCase.AuthenticateUser(c.Request.Context(), input)
	if err != nil {
		status, message := h.mapErrorToHTTPStatus(err)
		c.JSON(status, ErrorResponse{
			Error:   "Authentication failed",
			Message: message,
		})
		return
	}

	c.JSON(http.StatusOK, output.User)
}

// ErrorResponse representa uma resposta de erro padronizada
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// validateRole valida se o role fornecido é válido
func (h *UserHandler) validateRole(roleStr string) (user.Role, error) {
	role := user.Role(roleStr)
	switch role {
	case user.RoleAdmin, user.RoleUser, user.RoleGuest:
		return role, nil
	default:
		return "", user.ErrInvalidRole
	}
}

// mapErrorToHTTPStatus mapeia erros do domínio para códigos HTTP
func (h *UserHandler) mapErrorToHTTPStatus(err error) (int, string) {
	// Verifica se é um erro do domínio usando errors.Is
	if errors.Is(err, user.ErrInvalidRole) {
		return http.StatusBadRequest, "Invalid role"
	}
	if errors.Is(err, user.ErrUserNotFound) {
		return http.StatusNotFound, "User not found"
	}
	if errors.Is(err, user.ErrUserAlreadyExists) {
		return http.StatusConflict, "User already exists"
	}
	if errors.Is(err, user.ErrInvalidPassword) {
		return http.StatusUnauthorized, "Invalid password"
	}
	if errors.Is(err, user.ErrUserDeactivated) {
		return http.StatusUnauthorized, "User account is deactivated"
	}

	return http.StatusInternalServerError, "Internal server error"
}

package usecase

import (
	"context"
	"fmt"

	"go-api-boilerplate/internal/domain/auth"
	"go-api-boilerplate/internal/domain/repository"
	"go-api-boilerplate/internal/domain/user"
)

// UserUseCase implementa os casos de uso relacionados a usuários
type UserUseCase struct {
	userRepo   repository.UserRepository
	jwtService auth.JWTService
}

// NewUserUseCase cria uma nova instância de UserUseCase
func NewUserUseCase(userRepo repository.UserRepository, jwtService auth.JWTService) *UserUseCase {
	return &UserUseCase{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

// CreateUserInput representa os dados de entrada para criação de usuário
type CreateUserInput struct {
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Name     string    `json:"name"`
	Role     user.Role `json:"role"`
}

// CreateUserOutput representa os dados de saída da criação de usuário
type CreateUserOutput struct {
	User *user.User `json:"user"`
}

// CreateUser cria um novo usuário
func (uc *UserUseCase) CreateUser(ctx context.Context, input CreateUserInput) (*CreateUserOutput, error) {
	// Verifica se o email já existe
	exists, err := uc.userRepo.ExistsByEmail(ctx, input.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	}

	if exists {
		return nil, user.ErrUserAlreadyExists
	}

	// Cria a entidade User
	user, err := user.NewUser(input.Email, input.Password, input.Name, input.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to create user entity: %w", err)
	}

	// Persiste no repositório
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user in repository: %w", err)
	}

	return &CreateUserOutput{User: user}, nil
}

// GetUserByIDInput representa os dados de entrada para busca de usuário por ID
type GetUserByIDInput struct {
	ID string `json:"id"`
}

// GetUserByIDOutput representa os dados de saída da busca de usuário por ID
type GetUserByIDOutput struct {
	User *user.User `json:"user"`
}

// GetUserByID busca um usuário pelo ID
func (uc *UserUseCase) GetUserByID(ctx context.Context, input GetUserByIDInput) (*GetUserByIDOutput, error) {
	userEntity, err := uc.userRepo.GetByID(ctx, input.ID)
	if err != nil {
		// Propaga erros de domínio sem envolver
		if err == user.ErrUserNotFound {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return &GetUserByIDOutput{User: userEntity}, nil
}

// GetUserByEmailInput representa os dados de entrada para busca de usuário por email
type GetUserByEmailInput struct {
	Email string `json:"email"`
}

// GetUserByEmailOutput representa os dados de saída da busca de usuário por email
type GetUserByEmailOutput struct {
	User *user.User `json:"user"`
}

// GetUserByEmail busca um usuário pelo email
func (uc *UserUseCase) GetUserByEmail(ctx context.Context, input GetUserByEmailInput) (*GetUserByEmailOutput, error) {
	userEntity, err := uc.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		// Propaga erros de domínio sem envolver
		if err == user.ErrUserNotFound {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &GetUserByEmailOutput{User: userEntity}, nil
}

// UpdateUserInput representa os dados de entrada para atualização de usuário
type UpdateUserInput struct {
	ID    string     `json:"id"`
	Name  *string    `json:"name,omitempty"`
	Email *string    `json:"email,omitempty"`
	Role  *user.Role `json:"role,omitempty"`
}

// UpdateUserOutput representa os dados de saída da atualização de usuário
type UpdateUserOutput struct {
	User *user.User `json:"user"`
}

// UpdateUser atualiza um usuário existente
func (uc *UserUseCase) UpdateUser(ctx context.Context, input UpdateUserInput) (*UpdateUserOutput, error) {
	// Busca o usuário existente
	dbUser, err := uc.userRepo.GetByID(ctx, input.ID)
	if err != nil {
		// Propaga erros de domínio sem envolver
		if err == user.ErrUserNotFound {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get user for update: %w", err)
	}

	// Atualiza os campos fornecidos
	if input.Name != nil {
		if err := dbUser.UpdateName(*input.Name); err != nil {
			return nil, fmt.Errorf("failed to update name: %w", err)
		}
	}

	if input.Email != nil {
		// Verifica se o novo email já existe (se for diferente do atual)
		if *input.Email != dbUser.Email {
			exists, err := uc.userRepo.ExistsByEmail(ctx, *input.Email)
			if err != nil {
				return nil, fmt.Errorf("failed to check email existence: %w", err)
			}

			if exists {
				return nil, user.ErrUserAlreadyExists
			}
		}

		if err := dbUser.UpdateEmail(*input.Email); err != nil {
			return nil, fmt.Errorf("failed to update email: %w", err)
		}
	}

	if input.Role != nil {
		if err := dbUser.UpdateRole(*input.Role); err != nil {
			return nil, fmt.Errorf("failed to update role: %w", err)
		}
	}

	// Persiste as alterações
	if err := uc.userRepo.Update(ctx, dbUser); err != nil {
		return nil, fmt.Errorf("failed to update user in repository: %w", err)
	}

	return &UpdateUserOutput{User: dbUser}, nil
}

// DeleteUserInput representa os dados de entrada para exclusão de usuário
type DeleteUserInput struct {
	ID string `json:"id"`
}

// DeleteUser remove um usuário
func (uc *UserUseCase) DeleteUser(ctx context.Context, input DeleteUserInput) error {
	// Remove o usuário diretamente - o repositório retornará ErrUserNotFound se não existir
	if err := uc.userRepo.Delete(ctx, input.ID); err != nil {
		// Propaga erros de domínio sem envolver
		if err == user.ErrUserNotFound {
			return err
		}
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// ListUsersInput representa os dados de entrada para listagem de usuários
type ListUsersInput struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

// ListUsersOutput representa os dados de saída da listagem de usuários
type ListUsersOutput struct {
	Users []*user.User `json:"users"`
	Total int64        `json:"total"`
}

// ListUsers lista usuários com paginação
func (uc *UserUseCase) ListUsers(ctx context.Context, input ListUsersInput) (*ListUsersOutput, error) {
	// Valida parâmetros de paginação
	if input.Limit <= 0 {
		input.Limit = 10 // Default limit
	}

	if input.Offset < 0 {
		input.Offset = 0
	}

	// Busca usuários
	users, err := uc.userRepo.List(ctx, input.Offset, input.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	// Conta total de usuários
	total, err := uc.userRepo.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	return &ListUsersOutput{
		Users: users,
		Total: total,
	}, nil
}

// AuthenticateUserInput representa os dados de entrada para autenticação
type AuthenticateUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthenticateUserOutput representa os dados de saída da autenticação
type AuthenticateUserOutput struct {
	User  *user.User `json:"user"`
	Token string     `json:"token"`
}

// AuthenticateUser autentica um usuário
func (uc *UserUseCase) AuthenticateUser(ctx context.Context, input AuthenticateUserInput) (*AuthenticateUserOutput, error) {
	// Busca o usuário pelo email
	userEntity, err := uc.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		// Se o usuário não foi encontrado, retorna erro de domínio
		if err == user.ErrUserNotFound {
			return nil, user.ErrInvalidPassword
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	// Verifica se o usuário está ativo
	if !userEntity.IsActiveUser() {
		return nil, user.ErrUserDeactivated
	}

	// Verifica a senha
	if !userEntity.CheckPassword(input.Password) {
		return nil, user.ErrInvalidPassword
	}

	// Gera o token JWT
	token, err := uc.jwtService.GenerateToken(userEntity.ID, userEntity.Email, string(userEntity.Role))
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &AuthenticateUserOutput{
		User:  userEntity,
		Token: token,
	}, nil
}

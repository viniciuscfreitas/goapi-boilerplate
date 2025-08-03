package repository

import (
	"context"

	"github.com/fisiopet/bp/internal/domain/user"
)

// UserRepository define os contratos para persistência de usuários
type UserRepository interface {
	// Create cria um novo usuário no repositório
	Create(ctx context.Context, user *user.User) error

	// GetByID busca um usuário pelo ID
	GetByID(ctx context.Context, id string) (*user.User, error)

	// GetByEmail busca um usuário pelo email
	GetByEmail(ctx context.Context, email string) (*user.User, error)

	// Update atualiza um usuário existente
	Update(ctx context.Context, user *user.User) error

	// Delete remove um usuário pelo ID
	Delete(ctx context.Context, id string) error

	// List retorna uma lista de usuários com paginação
	List(ctx context.Context, offset, limit int) ([]*user.User, error)

	// Count retorna o total de usuários
	Count(ctx context.Context) (int64, error)

	// ExistsByEmail verifica se existe um usuário com o email fornecido
	ExistsByEmail(ctx context.Context, email string) (bool, error)

	// ExistsByID verifica se existe um usuário com o ID fornecido
	ExistsByID(ctx context.Context, id string) (bool, error)
} 
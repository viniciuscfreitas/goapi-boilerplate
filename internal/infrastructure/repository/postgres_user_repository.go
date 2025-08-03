package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	domainRepo "github.com/fisiopet/bp/internal/domain/repository"
	"github.com/fisiopet/bp/internal/domain/user"
	db "github.com/fisiopet/bp/internal/infrastructure/database"
	"github.com/google/uuid"
)

// PostgresUserRepository implementa UserRepository usando PostgreSQL
type PostgresUserRepository struct {
	db      *sql.DB
	querier *db.Queries
}

// NewPostgresUserRepository cria uma nova instância de PostgresUserRepository
func NewPostgresUserRepository(sqlDB *sql.DB) domainRepo.UserRepository {
	return &PostgresUserRepository{
		db:      sqlDB,
		querier: db.New(sqlDB),
	}
}

// Create cria um novo usuário no repositório
func (r *PostgresUserRepository) Create(ctx context.Context, u *user.User) error {
	// Gera um novo UUID se não existir
	if u.ID == "" {
		u.ID = uuid.New().String()
	}

	// Define timestamps se não existirem
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
	}
	if u.UpdatedAt.IsZero() {
		u.UpdatedAt = time.Now()
	}

	// Insere no banco de dados
	dbUser, err := r.querier.CreateUser(ctx, db.CreateUserParams{
		Email:     u.Email,
		Password:  u.Password,
		Name:      u.Name,
		Role:      string(u.Role),
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	})
	if err != nil {
		return fmt.Errorf("failed to create user in database: %w", err)
	}

	// Atualiza a entidade com os dados do banco
	r.mapDBUserToDomainUser(&dbUser, u)

	return nil
}

// GetByID busca um usuário pelo ID
func (r *PostgresUserRepository) GetByID(ctx context.Context, id string) (*user.User, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	dbUser, err := r.querier.GetUserByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, user.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return r.mapDBUserToDomainUser(&dbUser, nil), nil
}

// GetByEmail busca um usuário pelo email
func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	dbUser, err := r.querier.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, user.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return r.mapDBUserToDomainUser(&dbUser, nil), nil
}

// Update atualiza um usuário existente
func (r *PostgresUserRepository) Update(ctx context.Context, u *user.User) error {
	// Atualiza o timestamp
	u.UpdatedAt = time.Now()

	// Converte string ID para UUID
	userID, err := uuid.Parse(u.ID)
	if err != nil {
		return fmt.Errorf("invalid user ID format: %w", err)
	}

	// Atualiza no banco de dados
	dbUser, err := r.querier.UpdateUser(ctx, db.UpdateUserParams{
		ID:        userID,
		Email:     u.Email,
		Password:  u.Password,
		Name:      u.Name,
		Role:      string(u.Role),
		IsActive:  u.IsActive,
		UpdatedAt: u.UpdatedAt,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return user.ErrUserNotFound
		}
		return fmt.Errorf("failed to update user in database: %w", err)
	}

	// Atualiza a entidade com os dados do banco
	r.mapDBUserToDomainUser(&dbUser, u)

	return nil
}

// Delete remove um usuário pelo ID
func (r *PostgresUserRepository) Delete(ctx context.Context, id string) error {
	userID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid user ID format: %w", err)
	}

	err = r.querier.DeleteUser(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return user.ErrUserNotFound
		}
		return fmt.Errorf("failed to delete user from database: %w", err)
	}

	return nil
}

// List retorna uma lista de usuários com paginação
func (r *PostgresUserRepository) List(ctx context.Context, offset, limit int) ([]*user.User, error) {
	dbUsers, err := r.querier.ListUsers(ctx, db.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list users from database: %w", err)
	}

	users := make([]*user.User, len(dbUsers))
	for i, dbUser := range dbUsers {
		users[i] = r.mapDBUserToDomainUser(&dbUser, nil)
	}

	return users, nil
}

// Count retorna o total de usuários
func (r *PostgresUserRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.querier.CountUsers(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count users in database: %w", err)
	}

	return count, nil
}

// ExistsByEmail verifica se existe um usuário com o email fornecido
func (r *PostgresUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	exists, err := r.querier.ExistsByEmail(ctx, email)
	if err != nil {
		return false, fmt.Errorf("failed to check email existence in database: %w", err)
	}

	return exists, nil
}

// ExistsByID verifica se existe um usuário com o ID fornecido
func (r *PostgresUserRepository) ExistsByID(ctx context.Context, id string) (bool, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return false, fmt.Errorf("invalid user ID format: %w", err)
	}

	exists, err := r.querier.ExistsByID(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to check ID existence in database: %w", err)
	}

	return exists, nil
}

// mapDBUserToDomainUser mapeia um User do banco de dados para a entidade de domínio
func (r *PostgresUserRepository) mapDBUserToDomainUser(dbUser *db.User, domainUser *user.User) *user.User {
	if domainUser == nil {
		domainUser = &user.User{}
	}

	domainUser.ID = dbUser.ID.String()
	domainUser.Email = dbUser.Email
	domainUser.Password = dbUser.Password
	domainUser.Name = dbUser.Name
	domainUser.Role = user.Role(dbUser.Role)
	domainUser.IsActive = dbUser.IsActive
	domainUser.CreatedAt = dbUser.CreatedAt
	domainUser.UpdatedAt = dbUser.UpdatedAt

	return domainUser
}

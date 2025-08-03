package user

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Erros personalizados do domínio
var (
	ErrInvalidRole        = errors.New("invalid role")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidPassword    = errors.New("invalid password")
	ErrUserDeactivated    = errors.New("user account is deactivated")
)

// User representa a entidade de usuário no domínio
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Não exposto na serialização JSON
	Name      string    `json:"name"`
	Role      Role      `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Role representa o papel/permissão do usuário
type Role string

const (
	RoleAdmin  Role = "admin"
	RoleUser   Role = "user"
	RoleGuest  Role = "guest"
)

// NewUser cria uma nova instância de User
func NewUser(email, password, name string, role Role) (*User, error) {
	user := &User{
		Email:     email,
		Name:      name,
		Role:      role,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := user.SetPassword(password); err != nil {
		return nil, err
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	return user, nil
}

// SetPassword define a senha do usuário com hash bcrypt
func (u *User) SetPassword(password string) error {
	if password == "" {
		return errors.New("password cannot be empty")
	}

	if len(password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword verifica se a senha fornecida corresponde à senha do usuário
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// Validate valida os campos da entidade User
func (u *User) Validate() error {
	if u.Email == "" {
		return errors.New("email cannot be empty")
	}

	if u.Name == "" {
		return errors.New("name cannot be empty")
	}

	if u.Password == "" {
		return errors.New("password cannot be empty")
	}

	if !isValidRole(u.Role) {
		return errors.New("invalid role")
	}

	return nil
}

// UpdateName atualiza o nome do usuário
func (u *User) UpdateName(name string) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}

	u.Name = name
	u.UpdatedAt = time.Now()
	return nil
}

// UpdateEmail atualiza o email do usuário
func (u *User) UpdateEmail(email string) error {
	if email == "" {
		return errors.New("email cannot be empty")
	}

	u.Email = email
	u.UpdatedAt = time.Now()
	return nil
}

// UpdateRole atualiza o papel do usuário
func (u *User) UpdateRole(role Role) error {
	if !isValidRole(role) {
		return errors.New("invalid role")
	}

	u.Role = role
	u.UpdatedAt = time.Now()
	return nil
}

// Activate ativa o usuário
func (u *User) Activate() {
	u.IsActive = true
	u.UpdatedAt = time.Now()
}

// Deactivate desativa o usuário
func (u *User) Deactivate() {
	u.IsActive = false
	u.UpdatedAt = time.Now()
}

// isValidRole verifica se o papel é válido
func isValidRole(role Role) bool {
	switch role {
	case RoleAdmin, RoleUser, RoleGuest:
		return true
	default:
		return false
	}
}

// IsAdmin verifica se o usuário é administrador
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// IsActive verifica se o usuário está ativo
func (u *User) IsActiveUser() bool {
	return u.IsActive
} 
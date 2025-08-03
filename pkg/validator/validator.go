package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// CustomValidator implementa validação customizada
type CustomValidator struct {
	validator *validator.Validate
}

// NewCustomValidator cria uma nova instância de CustomValidator
func NewCustomValidator() *CustomValidator {
	v := validator.New()

	// Registra validações customizadas
	v.RegisterValidation("password", validatePassword)
	v.RegisterValidation("role", validateRole)

	return &CustomValidator{
		validator: v,
	}
}

// Validate valida uma struct
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// ValidateVar valida uma variável
func (cv *CustomValidator) ValidateVar(field interface{}, tag string) error {
	return cv.validator.Var(field, tag)
}

// GetValidationErrors retorna erros de validação formatados
func (cv *CustomValidator) GetValidationErrors(err error) []string {
	var errors []string

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			errors = append(errors, formatValidationError(e))
		}
	} else {
		errors = append(errors, err.Error())
	}

	return errors
}

// validatePassword valida se a senha atende aos critérios
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Mínimo 6 caracteres
	if len(password) < 6 {
		return false
	}

	// Deve conter pelo menos uma letra e um número
	hasLetter := false
	hasNumber := false

	for _, char := range password {
		if char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z' {
			hasLetter = true
		}
		if char >= '0' && char <= '9' {
			hasNumber = true
		}
	}

	return hasLetter && hasNumber
}

// validateRole valida se o role é válido
func validateRole(fl validator.FieldLevel) bool {
	role := fl.Field().String()
	validRoles := []string{"admin", "user", "guest"}

	for _, validRole := range validRoles {
		if role == validRole {
			return true
		}
	}

	return false
}

// formatValidationError formata um erro de validação
func formatValidationError(e validator.FieldError) string {
	field := strings.ToLower(e.Field())

	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", field, e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters long", field, e.Param())
	case "password":
		return fmt.Sprintf("%s must be at least 6 characters long and contain both letters and numbers", field)
	case "role":
		return fmt.Sprintf("%s must be one of: admin, user, guest", field)
	default:
		return fmt.Sprintf("%s failed validation: %s", field, e.Tag())
	}
}

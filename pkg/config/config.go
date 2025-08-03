package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config representa a configuração da aplicação
type Config struct {
	Server     ServerConfig     `mapstructure:"server"`
	Database   DatabaseConfig   `mapstructure:"database"`
	Logging    LoggingConfig    `mapstructure:"logging"`
	Security   SecurityConfig   `mapstructure:"security"`
	Environment string          `mapstructure:"environment"`
}

// ServerConfig representa as configurações do servidor
type ServerConfig struct {
	Host         string        `mapstructure:"host"`
	Port         string        `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

// DatabaseConfig representa as configurações do banco de dados
type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            string        `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	Name            string        `mapstructure:"name"`
	SSLMode         string        `mapstructure:"ssl_mode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

// LoggingConfig representa as configurações de logging
type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

// SecurityConfig representa as configurações de segurança
type SecurityConfig struct {
	BcryptCost    int           `mapstructure:"bcrypt_cost"`
	JWTSecret     string        `mapstructure:"jwt_secret"`
	JWTExpiration time.Duration `mapstructure:"jwt_expiration"`
}

// Load carrega a configuração do arquivo e variáveis de ambiente
func Load() (*Config, error) {
	// Configurar Viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// Configurar variáveis de ambiente
	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()

	// Mapear variáveis de ambiente para configurações
	setupEnvMappings()

	// Ler arquivo de configuração
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// Se o arquivo não for encontrado, usar apenas variáveis de ambiente
	}

	// Criar estrutura de configuração
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validar configuração
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// setupEnvMappings configura o mapeamento de variáveis de ambiente
func setupEnvMappings() {
	// Server
	viper.BindEnv("server.host", "APP_SERVER_HOST")
	viper.BindEnv("server.port", "APP_SERVER_PORT")
	viper.BindEnv("server.read_timeout", "APP_SERVER_READ_TIMEOUT")
	viper.BindEnv("server.write_timeout", "APP_SERVER_WRITE_TIMEOUT")
	viper.BindEnv("server.idle_timeout", "APP_SERVER_IDLE_TIMEOUT")

	// Database
	viper.BindEnv("database.host", "APP_DB_HOST")
	viper.BindEnv("database.port", "APP_DB_PORT")
	viper.BindEnv("database.user", "APP_DB_USER")
	viper.BindEnv("database.password", "APP_DB_PASSWORD")
	viper.BindEnv("database.name", "APP_DB_NAME")
	viper.BindEnv("database.ssl_mode", "APP_DB_SSL_MODE")
	viper.BindEnv("database.max_open_conns", "APP_DB_MAX_OPEN_CONNS")
	viper.BindEnv("database.max_idle_conns", "APP_DB_MAX_IDLE_CONNS")
	viper.BindEnv("database.conn_max_lifetime", "APP_DB_CONN_MAX_LIFETIME")

	// Logging
	viper.BindEnv("logging.level", "APP_LOG_LEVEL")
	viper.BindEnv("logging.format", "APP_LOG_FORMAT")
	viper.BindEnv("logging.output", "APP_LOG_OUTPUT")

	// Security
	viper.BindEnv("security.bcrypt_cost", "APP_BCRYPT_COST")
	viper.BindEnv("security.jwt_secret", "APP_JWT_SECRET")
	viper.BindEnv("security.jwt_expiration", "APP_JWT_EXPIRATION")

	// Environment
	viper.BindEnv("environment", "APP_ENV")
}

// Validate valida a configuração
func (c *Config) Validate() error {
	// Validar servidor
	if c.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}

	// Validar banco de dados
	if c.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if c.Database.Port == "" {
		return fmt.Errorf("database port is required")
	}
	if c.Database.User == "" {
		return fmt.Errorf("database user is required")
	}
	if c.Database.Name == "" {
		return fmt.Errorf("database name is required")
	}

	// Validar segurança
	if c.Security.JWTSecret == "" {
		return fmt.Errorf("jwt secret is required")
	}

	return nil
}

// GetDSN retorna a string de conexão do banco de dados
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode)
}

// IsDevelopment retorna true se o ambiente for development
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction retorna true se o ambiente for production
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// IsTesting retorna true se o ambiente for testing
func (c *Config) IsTesting() bool {
	return c.Environment == "testing"
} 
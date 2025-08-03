package database

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/pressly/goose/v3"
	_ "github.com/lib/pq"
)

// Migrator gerencia migrações do banco de dados
type Migrator struct {
	db     *sql.DB
	logger *slog.Logger
}

// NewMigrator cria uma nova instância de Migrator
func NewMigrator(db *sql.DB, logger *slog.Logger) *Migrator {
	return &Migrator{
		db:     db,
		logger: logger,
	}
}

// RunMigrations executa todas as migrações pendentes
func (m *Migrator) RunMigrations(migrationsDir string) error {
	m.logger.Info("Starting database migrations", "dir", migrationsDir)
	
	// Configurar goose
	goose.SetLogger(m.createGooseLogger())
	
	// Executar migrações
	if err := goose.Up(m.db, migrationsDir); err != nil {
		m.logger.Error("Failed to run migrations", "error", err)
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	
	m.logger.Info("Database migrations completed successfully")
	return nil
}

// CreateMigrationsTable cria a tabela de controle de migrações se não existir
func (m *Migrator) CreateMigrationsTable() error {
	m.logger.Info("Creating migrations table")
	
	if err := goose.Up(m.db, "."); err != nil {
		m.logger.Error("Failed to create migrations table", "error", err)
		return fmt.Errorf("failed to create migrations table: %w", err)
	}
	
	m.logger.Info("Migrations table created successfully")
	return nil
}

// GetMigrationStatus retorna o status das migrações
func (m *Migrator) GetMigrationStatus(migrationsDir string) error {
	m.logger.Info("Checking migration status", "dir", migrationsDir)
	
	if err := goose.Status(m.db, migrationsDir); err != nil {
		m.logger.Error("Failed to get migration status", "error", err)
		return fmt.Errorf("failed to get migration status: %w", err)
	}
	
	return nil
}

// Rollback executa rollback da última migração
func (m *Migrator) Rollback(migrationsDir string) error {
	m.logger.Info("Rolling back last migration", "dir", migrationsDir)
	
	if err := goose.Down(m.db, migrationsDir); err != nil {
		m.logger.Error("Failed to rollback migration", "error", err)
		return fmt.Errorf("failed to rollback migration: %w", err)
	}
	
	m.logger.Info("Migration rollback completed successfully")
	return nil
}

// createGooseLogger cria um logger compatível com goose
func (m *Migrator) createGooseLogger() goose.Logger {
	return &gooseLogger{logger: m.logger}
}

// gooseLogger implementa o logger do goose
type gooseLogger struct {
	logger *slog.Logger
}

func (l *gooseLogger) Fatal(v ...interface{}) {
	l.logger.Error("Goose fatal error", "message", fmt.Sprint(v...))
}

func (l *gooseLogger) Fatalf(format string, v ...interface{}) {
	l.logger.Error("Goose fatal error", "message", fmt.Sprintf(format, v...))
}

func (l *gooseLogger) Print(v ...interface{}) {
	l.logger.Info("Goose info", "message", fmt.Sprint(v...))
}

func (l *gooseLogger) Println(v ...interface{}) {
	l.logger.Info("Goose info", "message", fmt.Sprintln(v...))
}

func (l *gooseLogger) Printf(format string, v ...interface{}) {
	l.logger.Info("Goose info", "message", fmt.Sprintf(format, v...))
} 
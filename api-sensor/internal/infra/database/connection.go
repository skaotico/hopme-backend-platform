package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// Config define las variables de conexión necesarias para Postgres
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewConnection establece un pool de conexiones a la base de datos PostgreSQL
func NewConnection(cfg Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error al abrir la base de datos: %w", err)
	}

	// Configuración del Pool de Conexiones
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	// Verificar conexión física
	err = db.Ping()
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("no se pudo conectar a la base de datos: %w", err)
	}

	return db, nil
}

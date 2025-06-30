package database

import (
	"fmt"
	"os"

	"github.com/Noviiich/golang-new-relic-sentry-prometheus/internal/config"
	"github.com/jmoiron/sqlx"
)

// Database обертка для работы с БД
type Database struct {
	DB *sqlx.DB
}

// NewDatabase создает новое подключение к базе данных
func NewDatabase(config config.Database) (*Database, error) {
	var (
		db  *sqlx.DB
		err error
	)

	dbType := config.Type
	if dbType == "" {
		dbType = os.Getenv("DB_TYPE")
	}

	// Создаем подключение в зависимости от типа БД
	switch dbType {
	case "postgres":
		db, err = PostgreSQLConnection(config)
	case "mysql":
		db, err = MySQLConnection(config)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &Database{
		DB: db,
	}, nil
}

// ConnectionURLBuilder строит URL для подключения к БД
func ConnectionURLBuilder(dbType string, config config.Database) (string, error) {
	var url string

	switch dbType {
	case "postgres":
		// URL для PostgreSQL подключения
		url = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			config.Host,
			config.Port,
			config.Username,
			config.Password,
			config.DBName,
			config.SSLMode,
		)
	case "mysql":
		// URL для MySQL подключения
		url = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s",
			config.Username,
			config.Password,
			config.Host,
			config.Port,
			config.DBName,
		)
	default:
		return "", fmt.Errorf("unsupported database type: %s", dbType)
	}

	return url, nil
}

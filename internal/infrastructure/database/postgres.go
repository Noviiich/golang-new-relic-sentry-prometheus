package database

import (
	"fmt"

	"github.com/Noviiich/golang-new-relic-sentry-prometheus/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func PostgreSQLConnection(config config.Database) (*sqlx.DB, error) {
	postgresConnURL, err := ConnectionURLBuilder(config.Type, config)
	if err != nil {
		return nil, err
	}

	// Define database connection for PostgreSQL.
	db, err := sqlx.Connect("pgx", postgresConnURL)
	if err != nil {
		return nil, fmt.Errorf("error, not connected to database, %w", err)
	}

	db.SetMaxOpenConns(config.MaxDBConnections)
	db.SetMaxIdleConns(config.MaxDBIdleConnections)
	db.SetConnMaxLifetime(config.MaxDBLifetime)

	// Try to ping database.
	if err := db.Ping(); err != nil {
		defer db.Close() // close database connection
		return nil, fmt.Errorf("error, not sent ping to database, %w", err)
	}

	return db, nil
}

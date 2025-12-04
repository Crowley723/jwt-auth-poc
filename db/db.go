package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
	logger *slog.Logger
}

// New creates a new database connection
func New(dataSourceName string, logger *slog.Logger) (*DB, error) {
	if err := os.MkdirAll(filepath.Dir(dataSourceName), 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	sqlDB, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if _, err := sqlDB.Exec("PRAGMA foreign_keys = ON"); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	db := &DB{
		DB:     sqlDB,
		logger: logger,
	}

	logger.Info("Database connected successfully", "dsn", dataSourceName)
	return db, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	if db.DB != nil {
		db.logger.Info("Closing database connection")
		return db.DB.Close()
	}
	return nil
}

// Health checks if the database is accessible
func (db *DB) Health() error {
	return db.Ping()
}

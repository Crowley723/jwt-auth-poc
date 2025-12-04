package db

import (
	"embed"
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

//go:embed sqlite_migrations/*.sql
var migrationFiles embed.FS

type Migration struct {
	Version int
	Name    string
	SQL     string
}

// GetMigrations returns all migration files sorted by version
func GetMigrations() ([]Migration, error) {
	entries, err := migrationFiles.ReadDir("sqlite_migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to read migration directory: %w", err)
	}

	var migrations []Migration
	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		parts := strings.SplitN(entry.Name(), "_", 2)
		if len(parts) != 2 {
			continue
		}

		version, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}

		content, err := migrationFiles.ReadFile(filepath.Join("sqlite_migrations", entry.Name()))
		if err != nil {
			return nil, fmt.Errorf("failed to read migration file %s: %w", entry.Name(), err)
		}

		name := strings.TrimSuffix(parts[1], ".sql")

		migrations = append(migrations, Migration{
			Version: version,
			Name:    name,
			SQL:     string(content),
		})
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

// RunMigrations executes all pending migrations
func (db *DB) RunMigrations() error {
	if err := db.createMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	currentVersion, err := db.getCurrentVersion()
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	}

	migrations, err := GetMigrations()
	if err != nil {
		return fmt.Errorf("failed to get migrations: %w", err)
	}

	for _, migration := range migrations {
		if migration.Version <= currentVersion {
			db.logger.Debug("Skipping migration (already applied)",
				"version", migration.Version, "name", migration.Name)
			continue
		}

		db.logger.Info("Running migration",
			"version", migration.Version, "name", migration.Name)

		if err := db.executeMigration(migration); err != nil {
			return fmt.Errorf("failed to execute migration %d (%s): %w",
				migration.Version, migration.Name, err)
		}
	}

	db.logger.Info("All migrations completed successfully")
	return nil
}

func (db *DB) createMigrationsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := db.Exec(query)
	return err
}

func (db *DB) getCurrentVersion() (int, error) {
	query := "SELECT COALESCE(MAX(version), 0) FROM schema_migrations"
	var version int
	err := db.QueryRow(query).Scan(&version)
	return version, err
}

func (db *DB) executeMigration(migration Migration) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if strings.TrimSpace(migration.SQL) != "" {
		if _, err := tx.Exec(migration.SQL); err != nil {
			return fmt.Errorf("failed to execute migration SQL: %w", err)
		}
	}

	query := "INSERT INTO schema_migrations (version, name) VALUES (?, ?)"
	if _, err := tx.Exec(query, migration.Version, migration.Name); err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	return tx.Commit()
}

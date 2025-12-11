package database

import (
	"embed"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/mysql/*.sql
var mysqlMigrations embed.FS

//go:embed migrations/sqlite/*.sql
var sqliteMigrations embed.FS

// runMigrations runs database migrations using golang-migrate
func runMigrations(dbType DBType) error {
	switch dbType {
	case DBTypeMySQL:
		return runMySQLMigrationsV2()
	case DBTypeSQLite:
		return runSQLiteMigrationsV2()
	default:
		return fmt.Errorf("unsupported database type for migrations: %s", dbType)
	}
}

// runMySQLMigrationsV2 runs MySQL migrations using golang-migrate
func runMySQLMigrationsV2() error {
	// Create the source driver from embedded files
	sourceDriver, err := iofs.New(mysqlMigrations, "migrations/mysql")
	if err != nil {
		return fmt.Errorf("failed to create migration source: %w", err)
	}

	// Create the database driver
	dbDriver, err := mysql.WithInstance(DB, &mysql.Config{})
	if err != nil {
		return fmt.Errorf("failed to create MySQL migration driver: %w", err)
	}

	// Create the migrate instance
	m, err := migrate.NewWithInstance("iofs", sourceDriver, "mysql", dbDriver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// Run migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("MySQL migration failed: %w", err)
	}

	version, dirty, _ := m.Version()
	if dirty {
		log.Printf("Warning: MySQL migrations are in dirty state at version %d", version)
	} else {
		log.Printf("MySQL migrations completed (version: %d)", version)
	}

	return nil
}

// runSQLiteMigrationsV2 runs SQLite migrations using golang-migrate
func runSQLiteMigrationsV2() error {
	// Create the source driver from embedded files
	sourceDriver, err := iofs.New(sqliteMigrations, "migrations/sqlite")
	if err != nil {
		return fmt.Errorf("failed to create migration source: %w", err)
	}

	// Create the database driver
	dbDriver, err := sqlite.WithInstance(DB, &sqlite.Config{})
	if err != nil {
		return fmt.Errorf("failed to create SQLite migration driver: %w", err)
	}

	// Create the migrate instance
	m, err := migrate.NewWithInstance("iofs", sourceDriver, "sqlite", dbDriver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// Run migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("SQLite migration failed: %w", err)
	}

	version, dirty, _ := m.Version()
	if dirty {
		log.Printf("Warning: SQLite migrations are in dirty state at version %d", version)
	} else {
		log.Printf("SQLite migrations completed (version: %d)", version)
	}

	return nil
}

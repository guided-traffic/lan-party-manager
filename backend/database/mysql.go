package database

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
)

// MySQLConfig holds MySQL connection configuration
type MySQLConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string

	// TLS configuration
	TLSEnabled    bool
	TLSSkipVerify bool
	TLSCACert     string // Path to CA certificate file

	// Connection pool configuration
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// DefaultMySQLConfig returns a MySQLConfig with sensible defaults
func DefaultMySQLConfig() MySQLConfig {
	return MySQLConfig{
		Host:            "localhost",
		Port:            3306,
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 1 * time.Minute,
	}
}

// initMySQL initializes a MySQL database connection
func initMySQL(cfg MySQLConfig) error {
	// First, try to create the database if it doesn't exist
	if err := ensureMySQLDatabaseExists(cfg); err != nil {
		return fmt.Errorf("failed to ensure database exists: %w", err)
	}

	// Build MySQL DSN
	mysqlCfg := mysql.NewConfig()
	mysqlCfg.User = cfg.User
	mysqlCfg.Passwd = cfg.Password
	mysqlCfg.Net = "tcp"
	mysqlCfg.Addr = fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	mysqlCfg.DBName = cfg.Database
	mysqlCfg.ParseTime = true
	mysqlCfg.Loc = time.UTC
	mysqlCfg.MultiStatements = true
	mysqlCfg.InterpolateParams = true

	// Configure TLS if enabled
	if cfg.TLSEnabled {
		tlsConfig, err := buildTLSConfig(cfg)
		if err != nil {
			return fmt.Errorf("failed to configure TLS: %w", err)
		}

		// Register the TLS config with a unique name
		tlsConfigName := "custom"
		if err := mysql.RegisterTLSConfig(tlsConfigName, tlsConfig); err != nil {
			return fmt.Errorf("failed to register TLS config: %w", err)
		}
		mysqlCfg.TLSConfig = tlsConfigName
	}

	// Build DSN and open connection
	dsn := mysqlCfg.FormatDSN()
	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open MySQL database: %w", err)
	}

	// Configure connection pool
	DB.SetMaxOpenConns(cfg.MaxOpenConns)
	DB.SetMaxIdleConns(cfg.MaxIdleConns)
	DB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	DB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	// Test the connection
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping MySQL database: %w", err)
	}

	// Set database type
	dbType = DBTypeMySQL

	// Log connection info (without password)
	log.Printf("MySQL database initialized: %s@%s:%d/%s (TLS: %v)",
		cfg.User, cfg.Host, cfg.Port, cfg.Database, cfg.TLSEnabled)

	return nil
}

// ensureMySQLDatabaseExists connects without a database and creates it if necessary
func ensureMySQLDatabaseExists(cfg MySQLConfig) error {
	// Build MySQL DSN without database name
	mysqlCfg := mysql.NewConfig()
	mysqlCfg.User = cfg.User
	mysqlCfg.Passwd = cfg.Password
	mysqlCfg.Net = "tcp"
	mysqlCfg.Addr = fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	mysqlCfg.ParseTime = true
	mysqlCfg.Loc = time.UTC
	mysqlCfg.MultiStatements = true

	// Configure TLS if enabled
	if cfg.TLSEnabled {
		tlsConfig, err := buildTLSConfig(cfg)
		if err != nil {
			return fmt.Errorf("failed to configure TLS: %w", err)
		}

		tlsConfigName := "custom-init"
		if err := mysql.RegisterTLSConfig(tlsConfigName, tlsConfig); err != nil {
			// Ignore error if already registered
			if err.Error() != "tls: failed to find any PEM data in certificate input" {
				// Try to use existing config
			}
		}
		mysqlCfg.TLSConfig = tlsConfigName
	}

	dsn := mysqlCfg.FormatDSN()
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open MySQL connection: %w", err)
	}
	defer db.Close()

	// Test the connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping MySQL server: %w", err)
	}

	// Create database if it doesn't exist
	createDBSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", cfg.Database)
	_, err = db.Exec(createDBSQL)
	if err != nil {
		return fmt.Errorf("failed to create database '%s': %w", cfg.Database, err)
	}

	log.Printf("Ensured MySQL database '%s' exists", cfg.Database)
	return nil
}

// buildTLSConfig creates a TLS configuration for MySQL
func buildTLSConfig(cfg MySQLConfig) (*tls.Config, error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: cfg.TLSSkipVerify,
		MinVersion:         tls.VersionTLS12,
	}

	// If a CA certificate is provided, load and use it
	if cfg.TLSCACert != "" {
		caCert, err := os.ReadFile(cfg.TLSCACert)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA certificate: %w", err)
		}

		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("failed to parse CA certificate")
		}

		tlsConfig.RootCAs = caCertPool
		// When using a custom CA, we still want to verify the server cert
		// but allow skip verify for hostname only if configured
		if !cfg.TLSSkipVerify {
			tlsConfig.InsecureSkipVerify = false
		}
	}

	return tlsConfig, nil
}

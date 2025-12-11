package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

// ErrBusy is returned when SQLite is busy after all retries
var ErrBusy = errors.New("database is busy, please try again")

// initSQLite initializes a SQLite database connection
func initSQLite(dbPath string) error {
	// Ensure the directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create database directory: %w", err)
	}

	// Open database connection with optimized settings for concurrent access
	// _journal_mode=WAL enables Write-Ahead Logging for better concurrent writes
	// _busy_timeout=10000 waits up to 10 seconds before returning SQLITE_BUSY
	// _synchronous=NORMAL is a good balance between safety and performance
	// _cache_size=1000 increases the page cache size
	// _foreign_keys=ON enables foreign key constraints
	// _txlock=immediate ensures write transactions get the lock immediately
	dsn := fmt.Sprintf("%s?_journal_mode=WAL&_busy_timeout=10000&_synchronous=NORMAL&_cache_size=1000&_foreign_keys=ON&_txlock=immediate", dbPath)

	var err error
	DB, err = sql.Open("sqlite", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool for SQLite with WAL mode
	// WAL mode allows multiple readers and one writer concurrently
	// We use a small pool to avoid connection overhead while allowing some concurrency
	DB.SetMaxOpenConns(5)                   // Allow multiple connections for concurrent reads
	DB.SetMaxIdleConns(2)                   // Keep some connections warm
	DB.SetConnMaxLifetime(5 * time.Minute)  // Recycle connections periodically
	DB.SetConnMaxIdleTime(1 * time.Minute)  // Close idle connections

	// Test the connection
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Verify WAL mode is enabled
	var journalMode string
	if err := DB.QueryRow("PRAGMA journal_mode").Scan(&journalMode); err != nil {
		log.Printf("Warning: Could not verify journal mode: %v", err)
	} else {
		log.Printf("SQLite journal mode: %s", journalMode)
	}

	// Set database type
	dbType = DBTypeSQLite

	log.Printf("SQLite database initialized: %s", dbPath)
	return nil
}

// isBusyError checks if an error is a SQLite BUSY error
func isBusyError(err error) bool {
	if err == nil {
		return false
	}
	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "busy") || strings.Contains(errStr, "locked")
}

// WithRetry executes a function with retry logic for SQLITE_BUSY errors
// It will retry up to maxRetries times with exponential backoff
// For MySQL, the function is executed without retry logic
func WithRetry(fn func() error) error {
	return WithRetryContext(context.Background(), fn)
}

// WithRetryContext executes a function with retry logic and context support
// For MySQL, the function is executed without retry logic
func WithRetryContext(ctx context.Context, fn func() error) error {
	// For MySQL, no retry needed - just execute the function
	if dbType == DBTypeMySQL {
		return fn()
	}

	// SQLite retry logic
	const maxRetries = 5
	baseDelay := 50 * time.Millisecond

	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		// Check context before each attempt
		if ctx.Err() != nil {
			return ctx.Err()
		}

		lastErr = fn()
		if lastErr == nil {
			return nil
		}

		// Only retry on SQLITE_BUSY errors
		if !isBusyError(lastErr) {
			return lastErr
		}

		// Log retry attempt
		if attempt > 0 {
			log.Printf("SQLite busy, retry attempt %d/%d", attempt+1, maxRetries)
		}

		// Exponential backoff: 50ms, 100ms, 200ms, 400ms, 800ms
		delay := baseDelay * time.Duration(1<<attempt)
		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	log.Printf("SQLite busy after %d retries: %v", maxRetries, lastErr)
	return ErrBusy
}

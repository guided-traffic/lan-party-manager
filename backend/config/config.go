package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	// Server
	Port        string
	FrontendURL string
	BackendURL  string

	// Steam
	SteamAPIKey string

	// JWT
	JWTSecret          string
	JWTExpirationDays  int

	// Credits
	CreditIntervalMinutes int
	CreditMax             int

	// Admin
	AdminSteamIDs []string
}

// Load reads configuration from environment variables
func Load() *Config {
	// Load .env file if it exists (for local development)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := &Config{
		Port:                  getEnv("PORT", "8080"),
		FrontendURL:          getEnv("FRONTEND_URL", "http://localhost:4200"),
		BackendURL:           getEnv("BACKEND_URL", "http://localhost:8080"),
		SteamAPIKey:          getEnv("STEAM_API_KEY", ""),
		JWTSecret:            getEnv("JWT_SECRET", ""),
		JWTExpirationDays:    getEnvAsInt("JWT_EXPIRATION_DAYS", 7),
		CreditIntervalMinutes: getEnvAsInt("CREDIT_INTERVAL_MINUTES", 10),
		CreditMax:            getEnvAsInt("CREDIT_MAX", 10),
		AdminSteamIDs:        getEnvAsStringSlice("ADMIN_STEAM_IDS", []string{}),
	}

	// Validate required configuration
	cfg.validate()

	return cfg
}

// validate checks that all required configuration is present
func (c *Config) validate() {
	if c.SteamAPIKey == "" {
		log.Println("WARNING: STEAM_API_KEY is not set - Steam profile data will not be available")
	}
	if c.JWTSecret == "" {
		log.Fatal("FATAL: JWT_SECRET must be set")
	}
}

// getEnv reads an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvAsInt reads an environment variable as integer or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsStringSlice reads an environment variable as a comma-separated list of strings
func getEnvAsStringSlice(key string, defaultValue []string) []string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		parts := strings.Split(value, ",")
		result := make([]string, 0, len(parts))
		for _, part := range parts {
			trimmed := strings.TrimSpace(part)
			if trimmed != "" {
				result = append(result, trimmed)
			}
		}
		return result
	}
	return defaultValue
}

// IsAdmin checks if the given Steam ID is in the admin list
func (c *Config) IsAdmin(steamID string) bool {
	for _, adminID := range c.AdminSteamIDs {
		if adminID == steamID {
			return true
		}
	}
	return false
}

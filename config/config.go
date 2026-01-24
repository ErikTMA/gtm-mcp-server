package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the GTM MCP Server.
type Config struct {
	// Server configuration
	Port      int
	BaseURL   string
	Transport string // "http" (default) or "stdio"

	// Google OAuth configuration
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURI  string
	GoogleTokenFile    string // Path to cached OAuth token file (for stdio mode)

	// JWT configuration
	JWTSecret string

	// Logging
	LogLevel string

	// GTM restrictions
	// If set, restricts all operations to this account ID only
	GTMAccountID string
	// If set (comma-separated), restricts operations to these container IDs only
	GTMContainerIDs []string

	// Token saving (for stdio mode setup)
	// If set, saves the OAuth token to this path after successful authentication
	SaveTokenPath string
}

// Load reads configuration from environment variables.
// It first attempts to load from .env file if present, then .env.local for overrides.
func Load() (*Config, error) {
	// Load .env file if it exists (ignore error if not found)
	_ = godotenv.Load()
	// Load .env.local for local development overrides (takes precedence)
	_ = godotenv.Overload(".env.local")

	cfg := &Config{
		Port:               getEnvInt("PORT", 8081),
		BaseURL:            getEnv("BASE_URL", "http://localhost:8081"),
		Transport:          getEnv("TRANSPORT", "http"),
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURI:  getEnv("GOOGLE_REDIRECT_URI", ""),
		GoogleTokenFile:    getEnv("GOOGLE_TOKEN_FILE", ""),
		JWTSecret:          getEnv("JWT_SECRET", ""),
		LogLevel:           getEnv("LOG_LEVEL", "info"),
		GTMAccountID:       getEnv("GTM_ACCOUNT_ID", ""),
		GTMContainerIDs:    getEnvList("GTM_CONTAINER_IDS", nil),
		SaveTokenPath:      getEnv("SAVE_TOKEN_PATH", ""),
	}

	// Validation is deferred to when auth is actually needed
	// This allows the server to start and respond to initialize/ping
	// even without OAuth credentials configured

	return cfg, nil
}

// ValidateAuth checks if OAuth credentials are configured.
func (c *Config) ValidateAuth() error {
	if c.GoogleClientID == "" {
		return fmt.Errorf("GOOGLE_CLIENT_ID is required for authentication")
	}
	if c.GoogleClientSecret == "" {
		return fmt.Errorf("GOOGLE_CLIENT_SECRET is required for authentication")
	}
	if c.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required for authentication")
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvList(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		parts := strings.Split(value, ",")
		result := make([]string, 0, len(parts))
		for _, p := range parts {
			trimmed := strings.TrimSpace(p)
			if trimmed != "" {
				result = append(result, trimmed)
			}
		}
		if len(result) > 0 {
			return result
		}
	}
	return defaultValue
}

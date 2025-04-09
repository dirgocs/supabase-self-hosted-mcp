package config

import (
	"os"
	"strconv"
)

// Config represents the application configuration
type Config struct {
	Supabase SupabaseConfig
	Server   ServerConfig
}

// SupabaseConfig contains Supabase connection details
type SupabaseConfig struct {
	URL        string
	Key        string
	AnonKey    string
	JWTSecret  string
	PGConnStr  string
}

// ServerConfig contains server configuration
type ServerConfig struct {
	Port int
	Env  string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	// Default port
	port := 3000
	portStr := os.Getenv("PORT")
	if portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}

	return &Config{
		Supabase: SupabaseConfig{
			URL:        getEnv("SUPABASE_URL", "http://localhost:8000"),
			Key:        getEnv("SUPABASE_KEY", ""),
			AnonKey:    getEnv("SUPABASE_ANON_KEY", ""),
			JWTSecret:  getEnv("SUPABASE_JWT_SECRET", ""),
			PGConnStr:  getEnv("PG_CONNECTION_STRING", ""),
		},
		Server: ServerConfig{
			Port: port,
			Env:  getEnv("GO_ENV", "development"),
		},
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

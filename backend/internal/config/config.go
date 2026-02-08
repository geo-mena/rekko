package config

import (
	"os"
)

type Config struct {
	DatabaseURL     string
	KarenaiAPIURL   string
	KarenaiAPIToken string
	FinnhubAPIKey   string
	ServerPort      string
	MigrationsPath  string
	DBDriver        string
	StaticDir       string
}

func Load() *Config {
	return &Config{
		DatabaseURL:     getEnv("DATABASE_URL", "postgresql://root@localhost:26257/stockdb?sslmode=disable"),
		KarenaiAPIURL:   getEnv("KARENAI_API_URL", "https://api.karenai.click"),
		KarenaiAPIToken: getEnv("KARENAI_AUTH_TOKEN", ""),
		FinnhubAPIKey:   getEnv("FINNHUB_API_KEY", ""),
		ServerPort:      getEnvWithFallback("SERVER_PORT", "PORT", "8080"),
		MigrationsPath:  getEnv("MIGRATIONS_PATH", "./migrations"),
		DBDriver:        getEnv("DB_DRIVER", "cockroachdb"),
		StaticDir:       getEnv("STATIC_DIR", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvWithFallback(primary, fallback, defaultValue string) string {
	if value := os.Getenv(primary); value != "" {
		return value
	}
	if value := os.Getenv(fallback); value != "" {
		return value
	}
	return defaultValue
}

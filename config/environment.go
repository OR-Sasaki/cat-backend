package config

import "os"

var (
	Port   = getEnvOrDefault("PORT", "8080")
	DBPath = getEnvOrDefault("DB_PATH", "db/app.db")
)

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

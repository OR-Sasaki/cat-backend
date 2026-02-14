package config

import (
	"os"
	"strconv"
)

var (
	Port               = fetchString("PORT", "8040")
	DBPath             = fetchString("DB_PATH", "db/app.db")
	JWTSecret          = fetchString("JWT_SECRET", "your-secret-key-change-in-production")
	JWTExpirationHours = fetchInt("JWT_EXPIRATION_HOURS", 24)
)

func fetchString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func fetchInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

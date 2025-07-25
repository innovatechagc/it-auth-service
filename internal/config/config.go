package config

import (
	"os"
	"strconv"
)

type Config struct {
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string
	Port              string
	FirebaseProjectID string
	LogLevel          string
	Environment       string
	RateLimitRPS      int
	RateLimitBurst    int
	JWTSecret         string
	VaultConfig       VaultConfig
}

type VaultConfig struct {
	Address string
	Token   string
	Path    string
}

func LoadConfig() Config {
	return Config{
		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnv("DB_PORT", "5432"),
		DBUser:            getEnv("DB_USER", "postgres"),
		DBPassword:        getEnv("DB_PASSWORD", "postgres"),
		DBName:            getEnv("DB_NAME", "itapp"),
		Port:              getEnv("PORT", "8082"), // Puerto diferente para auth service
		FirebaseProjectID: getEnv("FIREBASE_PROJECT_ID", ""),
		LogLevel:          getEnv("LOG_LEVEL", "info"),
		Environment:       getEnv("ENVIRONMENT", "development"),
		RateLimitRPS:      getEnvAsInt("RATE_LIMIT_RPS", 100),
		RateLimitBurst:    getEnvAsInt("RATE_LIMIT_BURST", 200),
		JWTSecret:         getEnv("JWT_SECRET", "default-secret-change-in-production"),
		VaultConfig: VaultConfig{
			Address: getEnv("VAULT_ADDR", "http://localhost:8200"),
			Token:   getEnv("VAULT_TOKEN", ""),
			Path:    getEnv("VAULT_PATH", "secret/"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
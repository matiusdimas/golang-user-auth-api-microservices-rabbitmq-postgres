package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	DatabaseURL    string 
	JWTSecret      string
	RabbitMQURL    string
	AllowedOrigins string
	IsProduction   bool 
}

func LoadConfig() *Config {
	if env := os.Getenv("ENVIRONMENT"); env != "production" {
		godotenv.Load()
	}

	isProd := strings.ToLower(getEnv("ENVIRONMENT", "development")) == "production"

	config := &Config{
		Port:           getEnv("PORT", "8080"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", ""),
		DBName:         getEnv("DB_NAME", "userdb"),
		DatabaseURL:    getEnv("DATABASE_URL", ""),
		JWTSecret:      getEnv("JWT_SECRET", "your-super-secret-jwt-key"),
		RabbitMQURL:    getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		AllowedOrigins: getEnv("ALLOWED_ORIGINS", "*"),
		IsProduction:   isProd, 
	}

	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func (c *Config) GetDBConnectionString() string {
	if c.DatabaseURL != "" && c.IsProduction {
		if !strings.Contains(c.DatabaseURL, "sslmode=") {
			return c.DatabaseURL + "?sslmode=require"
		}
		return c.DatabaseURL
	}
	sslMode := "disable"
	if c.IsProduction {
		sslMode = "require"
	}

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, sslMode)
}
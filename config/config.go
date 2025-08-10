package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort      string
	JWTSecretKey []byte

	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string

	InitialAdminEmail    string
	InitialAdminPassword string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		log.Fatalf("Invalid DB_PORT: %v", err)
	}

	return &Config{
		AppPort:      getEnv("APP_PORT", "8080"),
		JWTSecretKey: []byte(getEnv("JWT_SECRET_KEY", "default_secret")),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "mydb"),
		DBPort:     dbPort,

		InitialAdminEmail:    getEnv("INITIAL_ADMIN_EMAIL", "admin@shop.com"),
		InitialAdminPassword: getEnv("INITIAL_ADMIN_PASSWORD", "adminpassword"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

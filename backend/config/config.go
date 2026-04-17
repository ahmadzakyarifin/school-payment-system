package config

import (
	"os"

	"github.com/subosito/gotenv"
)

type Config struct {
	DBUser string
	DBPass string
	DBHost string
	DBPort string
	DBName string
	Port   string
	// development | production
	AppEnv string 

	JWTSecret string
}

func ConfigDB() (*Config, error) {
	_ = gotenv.Load()

	cfg := &Config{
		DBUser:    os.Getenv("DB_USER"),
		DBPass:    os.Getenv("DB_PASSWORD"),
		DBHost:    os.Getenv("DB_HOST"),
		DBPort:    os.Getenv("DB_PORT"),
		DBName:    os.Getenv("DB_NAME"),
		Port:      getEnv("PORT", "8080"),
		AppEnv:    getEnv("APP_ENV", "development"),
		JWTSecret: os.Getenv("JWT_SECRET"),
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

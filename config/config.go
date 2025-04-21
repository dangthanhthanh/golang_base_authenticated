// File: config/config.go
package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload" // auto_ Load .env file
)

type Config struct {
	Port      string
	JWTSecret string

	DBHost  string
	DBPort  string
	DBUser  string
	DBPass  string
	DBName  string
	SSLMode string

	RedisHost string
	RedisPort string
	RedisPass string
}

func LoadConfig() Config {

	return Config{
		Port:      os.Getenv("PORT"),
		JWTSecret: os.Getenv("JWT_SECRET"),

		DBHost:  os.Getenv("DATABASE_HOST"),
		DBPort:  os.Getenv("DATABASE_PORT"),
		DBUser:  os.Getenv("DATABASE_USER"),
		DBPass:  os.Getenv("DATABASE_PASSWORD"),
		DBName:  os.Getenv("DATABASE_NAME"),
		SSLMode: os.Getenv("DATABASE_SSL_MODE"),

		RedisHost: os.Getenv("REDIS_HOST"),
		RedisPort: os.Getenv("REDIS_PORT"),
		RedisPass: os.Getenv("REDIS_PASSWORD"),
	}
}

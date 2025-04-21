package db

import (
	"base-app/config"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // PostgreSQL driver
)

var DB *sql.DB

func Connect(cfg config.Config) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBName,
		cfg.SSLMode,
	)
	fmt.Printf("dsn =: %s\n", dsn)

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	err = DB.Ping()
	if err != nil {
		panic(fmt.Sprintf("Failed to ping database: %v", err))
	}
}

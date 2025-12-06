package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/umess/backend/pkg/config"
)

var DB *sql.DB

func Connect() error {
	var err error
	DB, err = sql.Open("postgres", config.AppConfig.DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings for 16k users
	DB.SetMaxOpenConns(100)        // Maximum open connections
	DB.SetMaxIdleConns(20)         // Maximum idle connections
	DB.SetConnMaxLifetime(5 * time.Minute) // Maximum connection lifetime

	fmt.Println("Database connected successfully")
	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}



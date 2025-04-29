package database

import (
	"database/sql"
	"log"
	"voucher-api/internal/config"

	_ "github.com/lib/pq"
)

func InitDB(cfg *config.Config) *sql.DB {
	log.Printf("Connecting to database with URL: %s", cfg.DatabaseURL)
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Successfully connected to database")
	return db
}

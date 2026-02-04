package configs

import (
	"database/sql"

	"github.com/gofiber/fiber/v3/log"
	_ "github.com/lib/pq"
)

func ConnectDatabase(url string) *sql.DB {
	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Info("Connected to the database")
	return db
}

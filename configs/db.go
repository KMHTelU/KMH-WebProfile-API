package configs

import (
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v3/log"
	_ "github.com/lib/pq"
)

func ConnectDatabase(url string) *sql.DB {
	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	// Batas pool koneksi untuk mencegah kehabisan koneksi di production.
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Info("Connected to the database")
	return db
}

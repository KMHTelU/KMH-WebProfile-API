package repositories

import (
	"database/sql"

	"github.com/KMHTelU/KMH-WebProfile-API/internal/generated"
)

type Repository struct {
	DB      *sql.DB
	Queries *generated.Queries
}

func InitializeRepository(db *sql.DB, queries *generated.Queries) *Repository {
	return &Repository{
		DB:      db,
		Queries: queries,
	}
}

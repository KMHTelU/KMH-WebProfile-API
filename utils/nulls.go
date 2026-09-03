package utils

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Helper konversi ke tipe nullable milik sqlc. Nilai kosong (string kosong,
// waktu nol, UUID nil) diperlakukan sebagai NULL, mengikuti konvensi yang sudah
// dipakai di seluruh service.

func NullString(value string) sql.NullString {
	return sql.NullString{String: value, Valid: value != ""}
}

func NullTime(value time.Time) sql.NullTime {
	return sql.NullTime{Time: value, Valid: !value.IsZero()}
}

func NullBool(value bool) sql.NullBool {
	return sql.NullBool{Bool: value, Valid: true}
}

func NullInt32(value int32) sql.NullInt32 {
	return sql.NullInt32{Int32: value, Valid: true}
}

func NullUUID(value uuid.UUID) uuid.NullUUID {
	return uuid.NullUUID{UUID: value, Valid: value != uuid.Nil}
}

// IsForeignKeyViolation mendeteksi error Postgres 23503 (foreign key
// violation) — dipakai untuk menerjemahkan kegagalan delete menjadi pesan
// yang bisa dipahami admin, bukan 500 generik.
func IsForeignKeyViolation(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "23503"
	}
	return false
}

// BoolPtr dan TimePtr mengubah nilai nullable sqlc menjadi pointer agar NULL
// di basis data menjadi null di JSON, bukan false / 0001-01-01.
func BoolPtr(value sql.NullBool) *bool {
	if !value.Valid {
		return nil
	}
	v := value.Bool
	return &v
}

func TimePtr(value sql.NullTime) *time.Time {
	if !value.Valid {
		return nil
	}
	v := value.Time
	return &v
}

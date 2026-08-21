package utils

import (
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
)

// Regresi: kolom JSONB (json.RawMessage) harus tersisip sebagai JSON asli,
// bukan berubah menjadi deretan angka byte (mis. ["a"] → [91,34,97,34,93]).
func TestSanitizeForJSONKeepsRawJSON(t *testing.T) {
	type sample struct {
		Responsibilities json.RawMessage `json:"responsibilities"`
		Programs         json.RawMessage `json:"programs"`
		Empty            json.RawMessage `json:"empty"`
	}

	got := SanitizeForJSON(sample{
		Responsibilities: json.RawMessage(`["health check rutin","dokumentasi"]`),
		Programs:         json.RawMessage(`[{"name":"Bazaar","description":"tahunan"}]`),
		Empty:            nil,
	})

	raw, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("marshal gagal: %v", err)
	}

	var decoded struct {
		Responsibilities []string                 `json:"responsibilities"`
		Programs         []map[string]interface{} `json:"programs"`
		Empty            interface{}              `json:"empty"`
	}
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("hasil sanitasi bukan JSON yang diharapkan: %v (raw: %s)", err, raw)
	}

	if len(decoded.Responsibilities) != 2 || decoded.Responsibilities[0] != "health check rutin" {
		t.Errorf("responsibilities rusak: %s", raw)
	}
	if len(decoded.Programs) != 1 || decoded.Programs[0]["name"] != "Bazaar" {
		t.Errorf("programs rusak: %s", raw)
	}
	if decoded.Empty != nil {
		t.Errorf("raw kosong seharusnya null, dapat: %v", decoded.Empty)
	}
}

func TestSanitizeForJSONUnwrapsNullTypes(t *testing.T) {
	type sample struct {
		Name         sql.NullString `json:"name"`
		IsActive     sql.NullBool   `json:"is_active"`
		LastLoginAt  sql.NullTime   `json:"last_login_at"`
		PhotoMediaID uuid.NullUUID  `json:"photo_media_id"`
		PasswordHash sql.NullString `json:"password_hash"`
		TokenHash    sql.NullString `json:"token_hash"`
		Count        sql.NullInt32  `json:"count"`
	}

	id := uuid.MustParse("5a188cac-c031-404c-ae40-c0f4ed18fc35")
	loggedInAt := time.Date(2026, 8, 3, 10, 0, 0, 0, time.UTC)

	got := SanitizeForJSON([]sample{
		{
			Name:         sql.NullString{String: "Agung", Valid: true},
			IsActive:     sql.NullBool{},
			LastLoginAt:  sql.NullTime{},
			PhotoMediaID: uuid.NullUUID{},
			PasswordHash: sql.NullString{String: "secret-hash", Valid: true},
			TokenHash:    sql.NullString{String: "token", Valid: true},
			Count:        sql.NullInt32{Int32: 3, Valid: true},
		},
		{
			Name:         sql.NullString{String: "KMH Super Admin", Valid: true},
			IsActive:     sql.NullBool{Bool: true, Valid: true},
			LastLoginAt:  sql.NullTime{Time: loggedInAt, Valid: true},
			PhotoMediaID: uuid.NullUUID{UUID: id, Valid: true},
			PasswordHash: sql.NullString{String: "another-hash", Valid: true},
			Count:        sql.NullInt32{},
		},
	})

	raw, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var decoded []map[string]interface{}
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(decoded) != 2 {
		t.Fatalf("len = %d, want 2", len(decoded))
	}

	if decoded[0]["name"] != "Agung" {
		t.Errorf("name[0] = %#v", decoded[0]["name"])
	}
	if decoded[0]["is_active"] != nil {
		t.Errorf("is_active[0] should be null, got %#v", decoded[0]["is_active"])
	}
	if decoded[0]["last_login_at"] != nil {
		t.Errorf("last_login_at[0] should be null, got %#v", decoded[0]["last_login_at"])
	}
	if _, exists := decoded[0]["password_hash"]; exists {
		t.Error("password_hash must be stripped from response")
	}
	if _, exists := decoded[0]["token_hash"]; exists {
		t.Error("token_hash must be stripped from response")
	}
	if decoded[0]["count"] != float64(3) {
		t.Errorf("count[0] = %#v", decoded[0]["count"])
	}

	if decoded[1]["is_active"] != true {
		t.Errorf("is_active[1] = %#v, want true", decoded[1]["is_active"])
	}
	if decoded[1]["photo_media_id"] != id.String() {
		t.Errorf("photo_media_id[1] = %#v", decoded[1]["photo_media_id"])
	}
	if decoded[1]["last_login_at"] != loggedInAt.Format(time.RFC3339Nano) {
		t.Errorf("last_login_at[1] = %#v", decoded[1]["last_login_at"])
	}
	if decoded[1]["count"] != nil {
		t.Errorf("count[1] should be null, got %#v", decoded[1]["count"])
	}
}

func TestSanitizeForJSONHandlesNilSliceAsEmptyArray(t *testing.T) {
	var rows []string
	got := SanitizeForJSON(rows)
	raw, _ := json.Marshal(got)
	if string(raw) != "[]" {
		t.Errorf("nil slice should become [], got %s", raw)
	}
}

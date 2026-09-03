-- name: InsertEvent :one
INSERT INTO events (id, title, slug, description, event_type, start_time, end_time, location, google_maps_url, registration_url, cover_media_id, status, is_published, division_id, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
RETURNING *;

-- name: UpdateEvent :one
UPDATE events
SET title = $2,
    slug = $3,
    description = $4,
    event_type = $5,
    start_time = $6,
    end_time = $7,
    location = $8,
    google_maps_url = $9,
    registration_url = $10,
    cover_media_id = $11,
    status = $12,
    is_published = $13,
    division_id = $14,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteEvent :exec
DELETE FROM events
WHERE id = $1;

-- name: ClearEventsDivision :exec
-- Melepas tautan divisi dari semua event (dipakai sebelum divisi dihapus).
UPDATE events
SET division_id = NULL,
    updated_at = NOW()
WHERE division_id = $1;

-- name: GetEventByID :one
SELECT events.*, media.*, divisions.name AS division_name, divisions.slug AS division_slug
FROM events
LEFT JOIN media ON events.cover_media_id = media.id
LEFT JOIN divisions ON events.division_id = divisions.id
WHERE events.id = $1;

-- name: ListEvents :many
SELECT events.*, media.*, divisions.name AS division_name, divisions.slug AS division_slug
FROM events
LEFT JOIN media ON events.cover_media_id = media.id
LEFT JOIN divisions ON events.division_id = divisions.id
ORDER BY events.start_time DESC
LIMIT $1 OFFSET $2;

-- name: GetEventBySlug :one
SELECT * FROM events
WHERE slug = $1;

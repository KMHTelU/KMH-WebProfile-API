-- name: InsertDivision :one
INSERT INTO divisions (id, name, slug, subtitle, description, responsibilities, programs, coordinator_id, is_active)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: UpdateDivision :one
UPDATE divisions
SET name = $2,
    slug = $3,
    subtitle = $4,
    description = $5,
    responsibilities = $6,
    programs = $7,
    coordinator_id = $8,
    is_active = $9,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateDivisionIcon :exec
UPDATE divisions
SET icon_media_id = $2,
    updated_at = NOW()
WHERE id = $1;

-- name: GetDivisionByID :one
SELECT *
FROM divisions
LEFT JOIN media ON divisions.icon_media_id = media.id
LEFT JOIN members ON divisions.coordinator_id = members.id
WHERE divisions.id = $1;

-- name: GetAllDivisions :many
SELECT *
FROM divisions
LEFT JOIN media ON divisions.icon_media_id = media.id
LEFT JOIN members ON divisions.coordinator_id = members.id
ORDER BY divisions.name ASC;

-- name: GetDivisionBySlug :one
SELECT *
FROM divisions
WHERE slug = $1;

-- name: DeleteDivision :exec
DELETE FROM divisions
WHERE id = $1;

-- name: ClearDivisionCoordinator :exec
-- Melepas seorang anggota dari posisi koordinator di semua divisi
-- (dipakai sebelum anggota tersebut dihapus).
UPDATE divisions
SET coordinator_id = NULL,
    updated_at = NOW()
WHERE coordinator_id = $1;

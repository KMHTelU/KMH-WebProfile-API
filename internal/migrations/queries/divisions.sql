-- name: InsertDivision :one
INSERT INTO divisions (id, name, slug, description, icon_media_id, coordinator_id, is_active)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateDivision :one
UPDATE divisions
SET name = $2,
    slug = $3,
    description = $4,
    icon_media_id = $5,
    coordinator_id = $6,
    is_active = $7,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetDivisionByID :one
SELECT *
FROM divisions 
INNER JOIN media ON divisions.icon_media_id = media.id
INNER JOIN members ON divisions.coordinator_id = members.id
WHERE divisions.id = $1;

-- name: GetAllDivisions :many
SELECT *
FROM divisions
INNER JOIN media ON divisions.icon_media_id = media.id
INNER JOIN members ON divisions.coordinator_id = members.id
ORDER BY divisions.name ASC;

-- name: DeleteDivision :exec
DELETE FROM divisions
WHERE id = $1;

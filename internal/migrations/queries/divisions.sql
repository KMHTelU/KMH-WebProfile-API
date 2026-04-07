-- name: InsertDivision :one
INSERT INTO divisions (id, name, slug, description, coordinator_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateDivision :one
UPDATE divisions
SET name = $2,
    slug = $3,
    description = $4,
    coordinator_id = $5,
    is_active = $6,
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

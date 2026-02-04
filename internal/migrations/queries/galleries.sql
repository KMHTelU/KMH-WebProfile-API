-- name: InsertGallery :one
INSERT INTO galleries (
  id,
  title,
  description,
  event_id,
  is_public
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: UpdateGallery :one
UPDATE galleries
SET
  title = $2,
  description = $3,
  event_id = $4,
  is_public = $5
WHERE id = $1
RETURNING *;

-- name: DeleteGallery :exec
DELETE FROM galleries
WHERE id = $1;

-- name: SelectAllGalleries :many
SELECT *
FROM galleries
INNER JOIN events ON galleries.event_id = events.id
ORDER BY galleries.created_at DESC
LIMIT $1 OFFSET $2;

-- name: SelectGalleryByID :one
SELECT *
FROM galleries
INNER JOIN events ON galleries.event_id = events.id
WHERE galleries.id = $1;
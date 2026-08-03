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
SELECT galleries.*, events.title AS event_title, events.slug AS event_slug, events.start_time AS event_start_time
FROM galleries
LEFT JOIN events ON galleries.event_id = events.id
WHERE (sqlc.narg('event_id')::uuid IS NULL OR galleries.event_id = sqlc.narg('event_id')::uuid)
  AND (sqlc.narg('is_public')::boolean IS NULL OR galleries.is_public = sqlc.narg('is_public')::boolean)
ORDER BY galleries.created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: SelectGalleryByID :one
SELECT galleries.*, events.title AS event_title, events.slug AS event_slug, events.start_time AS event_start_time
FROM galleries
LEFT JOIN events ON galleries.event_id = events.id
WHERE galleries.id = $1;

-- name: SelectGalleryCategories :many
SELECT
  events.id         AS event_id,
  events.title      AS event_title,
  events.slug       AS event_slug,
  events.start_time AS event_start_time,
  COUNT(galleries.id) AS gallery_count
FROM events
INNER JOIN galleries ON galleries.event_id = events.id
GROUP BY events.id, events.title, events.slug, events.start_time
ORDER BY events.start_time DESC NULLS LAST;

-- name: InsertMedia :one
INSERT INTO media (
  id,
  file_name,
  file_type,
  mime_type,
  file_size,
  url,
  alt_text,
  caption,
  uploaded_by
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: DeleteMedia :exec
DELETE FROM media
WHERE id = $1;

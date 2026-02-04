-- name: InsertContactMessage :one
INSERT INTO contact_messages (
  id,
  name,
  email,
  subject,
  message,
  is_read
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: MarkContactMessageAsRead :one
UPDATE contact_messages
SET
  is_read = TRUE
WHERE id = $1
RETURNING *;

-- name: DeleteContactMessage :exec
DELETE FROM contact_messages
WHERE id = $1;

-- name: ListContactMessages :many
SELECT *
FROM contact_messages
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: SelectContactMessageByID :one
SELECT *
FROM contact_messages
WHERE id = $1;
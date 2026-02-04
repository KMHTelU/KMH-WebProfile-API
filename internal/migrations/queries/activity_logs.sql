-- name: InsertActivityLog :one
INSERT INTO activity_logs (
  id,
  user_id,
  action,
  entity,
  entity_id,
  ip_address,
  user_agent
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: DeleteActivityLog :exec
DELETE FROM activity_logs
WHERE id = $1;

-- name: SelectActivityLogsByUserID :many
SELECT *
FROM activity_logs
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: SelectActivityLogByID :one
SELECT *
FROM activity_logs
WHERE id = $1;

-- name: ListActivityLogs :many
SELECT *
FROM activity_logs
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

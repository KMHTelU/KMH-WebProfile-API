-- name: GetUserByID :one
SELECT *
FROM users
INNER JOIN roles ON users.role_id = roles.id
WHERE users.id = $1;

-- name: GetUserByEmail :one
SELECT *
FROM users
INNER JOIN roles ON users.role_id = roles.id
WHERE users.email = $1;

-- name: GetUsers :many
SELECT *
FROM users
INNER JOIN roles ON users.role_id = roles.id
ORDER BY users.name ASC
LIMIT $1 OFFSET $2;

-- name: CreateUser :one
INSERT INTO users (id, name, email, password_hash, role_id, is_active, last_login_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET name = $2,
    email = $3,
    password_hash = $4,
    role_id = $5,
    is_active = $6,
    last_login_at = $7,
    updated_at = NOW()
WHERE users.id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE users.id = $1;
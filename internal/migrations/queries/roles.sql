-- name: GetRoleByID :one
SELECT id, name, description
FROM roles
WHERE id = $1;

-- name: GetAllRoles :many
SELECT id, name, description
FROM roles
ORDER BY name ASC;

-- name: InsertRole :one
INSERT INTO roles (id, name, description)
VALUES ($1, $2, $3)
RETURNING id, name, description;

-- name: UpdateRole :one
UPDATE roles
SET name = $2,
    description = $3
WHERE id = $1
RETURNING id, name, description;

-- name: DeleteRole :exec
DELETE FROM roles
WHERE id = $1;

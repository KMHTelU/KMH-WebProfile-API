-- name: InsertMember :one
INSERT INTO members (id, name, npm, email, phone, bio, instagram_url, period_start, period_end)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: UpdateMember :one
UPDATE members
SET name = $2,
    email = $3,
    phone = $4,
    bio = $5,
    instagram_url = $6,
    period_start = $7,
    period_end = $8,
    is_active = $9,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateMemberPhoto :exec
UPDATE members
SET photo_media_id = $2,
    updated_at = NOW()
WHERE id = $1;

-- name: GetMemberByID :one
SELECT *
FROM members
INNER JOIN media ON members.photo_media_id = media.id
WHERE members.id = $1;

-- name: GetAllMembers :many
SELECT *
FROM members
INNER JOIN media ON members.photo_media_id = media.id
ORDER BY members.name ASC
LIMIT $1 OFFSET $2;

-- name: DeleteMember :exec
DELETE FROM members
WHERE members.id = $1;
-- name: InsertMember :one
INSERT INTO members (id, name, nim, email, phone, bio, instagram_url, period_start, period_end, is_active)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: UpdateMember :one
UPDATE members
SET name = $2,
    nim = $3,
    email = $4,
    phone = $5,
    bio = $6,
    instagram_url = $7,
    period_start = $8,
    period_end = $9,
    is_active = $10,
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
LEFT JOIN media ON members.photo_media_id = media.id
WHERE members.id = $1;

-- name: GetAllMembers :many
SELECT *
FROM members
LEFT JOIN media ON members.photo_media_id = media.id
ORDER BY members.name ASC
LIMIT $1 OFFSET $2;

-- name: GetMemberByNIM :one
SELECT *
FROM members
WHERE nim = $1;

-- name: DeleteMember :exec
DELETE FROM members
WHERE members.id = $1;

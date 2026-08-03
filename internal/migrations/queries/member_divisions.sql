-- name: InsertMemberDivision :one
INSERT INTO member_divisions (id, member_id, division_id, role_title)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateMemberDivision :one
UPDATE member_divisions
SET member_id = $2,
    division_id = $3,
    role_title = $4
WHERE id = $1
RETURNING *;

-- name: GetMemberDivisionByID :one
SELECT *
FROM member_divisions
WHERE id = $1;

-- name: GetMemberDivisionByPair :one
SELECT *
FROM member_divisions
WHERE member_id = $1 AND division_id = $2;

-- name: GetMemberDivisionsByMemberID :many
SELECT *
FROM member_divisions
INNER JOIN divisions ON member_divisions.division_id = divisions.id
WHERE member_id = $1;

-- name: GetMemberDivisionsByDivisionID :many
SELECT *
FROM member_divisions
INNER JOIN members ON member_divisions.member_id = members.id
WHERE division_id = $1;

-- name: DeleteMemberDivision :exec
DELETE FROM member_divisions
WHERE id = $1;

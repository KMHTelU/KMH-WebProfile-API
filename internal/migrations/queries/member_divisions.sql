-- name: InsertMemberDivision :one
INSERT INTO member_divisions (member_id, division_id, role_title)
VALUES ($1, $2, $3)
RETURNING *;

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
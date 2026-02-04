-- name: InsertBlogTag :one
INSERT INTO blog_tags (
  id,
  name,
  slug
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: SelectBlogTagByID :one
SELECT *
FROM blog_tags
WHERE id = $1;

-- name: UpdateBlogTag :one
UPDATE blog_tags
SET
  name = $2,
  slug = $3
WHERE id = $1
RETURNING *;

-- name: DeleteBlogTag :exec
DELETE FROM blog_tags
WHERE id = $1;

-- name: ListBlogTags :many
SELECT *
FROM blog_tags
ORDER BY name ASC
LIMIT $1 OFFSET $2;
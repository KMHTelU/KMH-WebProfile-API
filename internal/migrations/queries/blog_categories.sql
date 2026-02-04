-- name: InsertBlogCategory :one
INSERT INTO blog_categories (id, name, slug)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateBlogCategory :one
UPDATE blog_categories
SET name = $2,
    slug = $3
WHERE id = $1
RETURNING *;

-- name: DeleteBlogCategory :exec
DELETE FROM blog_categories
WHERE id = $1;

-- name: GetBlogCategoryByID :one
SELECT * FROM blog_categories
WHERE id = $1;

-- name: ListBlogCategories :many
SELECT * FROM blog_categories
ORDER BY name ASC
LIMIT $1 OFFSET $2;
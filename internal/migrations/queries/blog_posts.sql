-- name: InsertBlogPost :one
INSERT INTO blog_posts (
  id,
  title,
  slug,
  excerpt,
  content,
  category_id,
  featured_media_id,
  author_id,
  status,
  published_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING *;

-- name: SelectBlogPostByID :one
SELECT *
FROM blog_posts
INNER JOIN users ON blog_posts.author_id = users.id
INNER JOIN blog_categories ON blog_posts.category_id = blog_categories.id
INNER JOIN media ON blog_posts.featured_media_id = media.id
WHERE blog_posts.id = $1;

-- name: UpdateBlogPost :one
UPDATE blog_posts
SET
  title = $2,
  slug = $3,
  excerpt = $4,
  content = $5,
  category_id = $6,
  featured_media_id = $7,
  author_id = $8,
  status = $9,
  published_at = $10,
  updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteBlogPost :exec
DELETE FROM blog_posts
WHERE id = $1;

-- name: ListBlogPosts :many
SELECT *
FROM blog_posts
INNER JOIN users ON blog_posts.author_id = users.id
INNER JOIN blog_categories ON blog_posts.category_id = blog_categories.id
INNER JOIN media ON blog_posts.featured_media_id = media.id
ORDER BY published_at DESC
LIMIT $1 OFFSET $2;

-- name: ListBlogPostsByCategory :many
SELECT *
FROM blog_posts
INNER JOIN users ON blog_posts.author_id = users.id
INNER JOIN blog_categories ON blog_posts.category_id = blog_categories.id
INNER JOIN media ON blog_posts.featured_media_id = media.id
WHERE category_id = $1
ORDER BY published_at DESC
LIMIT $2 OFFSET $3;
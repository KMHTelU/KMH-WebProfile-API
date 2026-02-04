-- name: InsertBlogPostTag :one
INSERT INTO blog_post_tags (
  post_id,
  tag_id
) VALUES (
  $1, $2
)
RETURNING *;

-- name: DeleteBlogPostTag :exec
DELETE FROM blog_post_tags
WHERE post_id = $1 AND tag_id = $2;

-- name: ListTagsByBlogPostID :many
SELECT bt.*
FROM blog_tags bt
INNER JOIN blog_post_tags bpt ON bt.id = bpt.tag_id
WHERE bpt.post_id = $1
ORDER BY bt.name ASC
LIMIT $2 OFFSET $3;

-- name: ListBlogPostsByTagID :many
SELECT bp.*
FROM blog_posts bp
INNER JOIN blog_post_tags bpt ON bp.id = bpt.post_id
WHERE bpt.tag_id = $1
ORDER BY bp.title ASC
LIMIT $2 OFFSET $3;

-- name: CountBlogPostsByTagID :one
SELECT COUNT(*)
FROM blog_post_tags
WHERE tag_id = $1;

-- name: CountTagsByBlogPostID :one
SELECT COUNT(*)
FROM blog_post_tags
WHERE post_id = $1;
-- name: InsertHomepageBanner :one
INSERT INTO homepage_banners (
  id,
  title,
  subtitle,
  media_id,
  cta_text,
  cta_url,
  is_active,
  start_date,
  end_date
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: UpdateHomepageBanner :one
UPDATE homepage_banners
SET title = $2,
    subtitle = $3,
    media_id = $4,
    cta_text = $5,
    cta_url = $6,
    is_active = $7,
    start_date = $8,
    end_date = $9
WHERE id = $1
RETURNING *;

-- name: GetHomepageBannerByID :one
SELECT *
FROM homepage_banners
WHERE id = $1;

-- name: DeleteHomepageBanner :exec
DELETE FROM homepage_banners
WHERE id = $1;

-- name: SelectAllHomepageBanners :many
SELECT *
FROM homepage_banners
LEFT JOIN media ON homepage_banners.media_id = media.id
ORDER BY start_date DESC
LIMIT $1 OFFSET $2;

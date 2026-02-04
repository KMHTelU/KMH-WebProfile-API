-- name: InsertGalleryItem :one
INSERT INTO gallery_items (
  id,
  gallery_id,
  media_id,
  sort_order
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: SelectGalleryItemsByGalleryID :many
SELECT *
FROM gallery_items
INNER JOIN media ON gallery_items.media_id = media.id
WHERE gallery_id = $1
ORDER BY sort_order ASC;

-- name: DeleteGalleryItem :exec
DELETE FROM gallery_items
WHERE id = $1;

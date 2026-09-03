-- name: InsertDivision :one
INSERT INTO divisions (id, name, slug, subtitle, description, division_type, responsibilities, programs, coordinator_id, is_active)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: UpdateDivision :one
UPDATE divisions
SET name = $2,
    slug = $3,
    subtitle = $4,
    description = $5,
    division_type = $6,
    responsibilities = $7,
    programs = $8,
    coordinator_id = $9,
    is_active = $10,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateDivisionIcon :exec
UPDATE divisions
SET icon_media_id = $2,
    updated_at = NOW()
WHERE id = $1;

-- name: GetDivisionByID :one
SELECT *
FROM divisions
LEFT JOIN media ON divisions.icon_media_id = media.id
LEFT JOIN members ON divisions.coordinator_id = members.id
WHERE divisions.id = $1;

-- name: GetAllDivisions :many
-- member_count = jumlah anggota unik divisi (penugasan member_divisions
-- digabung koordinator, tanpa dobel) — dipakai kartu daftar divisi publik.
SELECT divisions.*, media.*, members.*,
  (
    SELECT COUNT(*)
    FROM (
      SELECT md.member_id
      FROM member_divisions md
      WHERE md.division_id = divisions.id
      UNION
      SELECT divisions.coordinator_id
      WHERE divisions.coordinator_id IS NOT NULL
    ) AS team
  ) AS member_count
FROM divisions
LEFT JOIN media ON divisions.icon_media_id = media.id
LEFT JOIN members ON divisions.coordinator_id = members.id
ORDER BY divisions.name ASC;

-- name: GetDivisionBySlug :one
SELECT *
FROM divisions
WHERE slug = $1;

-- name: DeleteDivision :exec
DELETE FROM divisions
WHERE id = $1;

-- name: ClearDivisionCoordinator :exec
-- Melepas seorang anggota dari posisi koordinator di semua divisi
-- (dipakai sebelum anggota tersebut dihapus).
UPDATE divisions
SET coordinator_id = NULL,
    updated_at = NOW()
WHERE coordinator_id = $1;

-- Query untuk halaman struktur organisasi (organization tree) publik.
-- Foto orang diambil dari foto member yang sudah ada (members.photo_media_id),
-- jadi tidak perlu mekanisme unggah foto terpisah untuk tree.

-- name: SelectOrgTreeDivisions :many
SELECT
  divisions.id,
  divisions.name,
  divisions.slug,
  divisions.subtitle,
  divisions.description,
  divisions.responsibilities,
  members.id   AS coordinator_id,
  members.name AS coordinator_name,
  members.nim  AS coordinator_nim,
  media.url    AS coordinator_photo_url
FROM divisions
LEFT JOIN members ON divisions.coordinator_id = members.id
LEFT JOIN media ON members.photo_media_id = media.id
WHERE COALESCE(divisions.is_active, true) = true
ORDER BY divisions.name ASC;

-- name: SelectOrgTreeAssignments :many
SELECT
  member_divisions.id,
  member_divisions.role_title,
  member_divisions.created_at,
  divisions.id   AS division_id,
  divisions.slug AS division_slug,
  members.id     AS member_id,
  members.name   AS member_name,
  members.nim    AS member_nim,
  media.url      AS photo_url
FROM member_divisions
INNER JOIN divisions ON member_divisions.division_id = divisions.id
INNER JOIN members ON member_divisions.member_id = members.id
LEFT JOIN media ON members.photo_media_id = media.id
WHERE COALESCE(divisions.is_active, true) = true
  AND COALESCE(members.is_active, true) = true
ORDER BY member_divisions.created_at ASC;

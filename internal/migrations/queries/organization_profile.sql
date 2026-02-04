-- name: InsertOrganizationProfile :one
INSERT INTO organization_profile (id, name, short_name, description, vision, mission, history, logo_media_id, address, email, phone, instagram_url, youtube_url, website_url)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
RETURNING *;

-- name: UpdateOrganizationProfile :one
UPDATE organization_profile
SET name = $2,
    short_name = $3,
    description = $4,
    vision = $5,
    mission = $6,
    history = $7,
    logo_media_id = $8,
    address = $9,
    email = $10,
    phone = $11,
    instagram_url = $12,
    youtube_url = $13,
    website_url = $14,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetOrganizationProfile :one
SELECT *
FROM organization_profile
INNER JOIN media ON organization_profile.logo_media_id = media.id
WHERE organization_profile.id = $1;

-- name: DeleteOrganizationProfile :exec
DELETE FROM organization_profile
WHERE organization_profile.id = $1;
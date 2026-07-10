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
    address = $8,
    email = $9,
    phone = $10,
    instagram_url = $11,
    youtube_url = $12,
    website_url = $13,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateOrganizationProfileLogo :one
UPDATE organization_profile
SET logo_media_id = $2,
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
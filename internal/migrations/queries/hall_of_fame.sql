-- Query Hall of Fame (arsip museum 3D). Read publik satu tembakan +
-- CRUD per entitas untuk admin.

-- ── Generations ──

-- name: SelectHofGenerations :many
SELECT * FROM hof_generations
ORDER BY year_start ASC, sort_order ASC;

-- name: InsertHofGeneration :one
INSERT INTO hof_generations (id, name, year_start, year_end, description, milestones, accent, sort_order)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: UpdateHofGeneration :one
UPDATE hof_generations
SET name = $2,
    year_start = $3,
    year_end = $4,
    description = $5,
    milestones = $6,
    accent = $7,
    sort_order = $8,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteHofGeneration :exec
DELETE FROM hof_generations WHERE id = $1;

-- ── People ──

-- name: SelectHofPeople :many
SELECT hof_people.*, media.url AS photo_url
FROM hof_people
LEFT JOIN media ON hof_people.photo_media_id = media.id
ORDER BY hof_people.sort_order ASC, hof_people.name ASC;

-- name: InsertHofPerson :one
INSERT INTO hof_people (id, generation_id, name, role, study_program, biography, contributions, legacy, quote, fields, photo_media_id, sort_order)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING *;

-- name: UpdateHofPerson :one
UPDATE hof_people
SET generation_id = $2,
    name = $3,
    role = $4,
    study_program = $5,
    biography = $6,
    contributions = $7,
    legacy = $8,
    quote = $9,
    fields = $10,
    photo_media_id = $11,
    sort_order = $12,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteHofPerson :exec
DELETE FROM hof_people WHERE id = $1;

-- ── Achievements ──

-- name: SelectHofAchievements :many
SELECT * FROM hof_achievements
ORDER BY year ASC, created_at ASC;

-- name: InsertHofAchievement :one
INSERT INTO hof_achievements (id, person_id, title, category, year, organization, result, description)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: UpdateHofAchievement :one
UPDATE hof_achievements
SET person_id = $2,
    title = $3,
    category = $4,
    year = $5,
    organization = $6,
    result = $7,
    description = $8,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteHofAchievement :exec
DELETE FROM hof_achievements WHERE id = $1;

-- ── Timeline organisasi ──

-- name: SelectHofTimelineEvents :many
SELECT * FROM hof_timeline_events
ORDER BY sort_order ASC, created_at ASC;

-- name: InsertHofTimelineEvent :one
INSERT INTO hof_timeline_events (id, year_label, title, description, sort_order)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateHofTimelineEvent :one
UPDATE hof_timeline_events
SET year_label = $2,
    title = $3,
    description = $4,
    sort_order = $5
WHERE id = $1
RETURNING *;

-- name: DeleteHofTimelineEvent :exec
DELETE FROM hof_timeline_events WHERE id = $1;

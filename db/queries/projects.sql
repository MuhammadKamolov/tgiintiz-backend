-- name: ListActiveProjectsByLang :many
SELECT
    id,
    lang,
    slug,
    title,
    category,
    location,
    completed_year,
    short_desc,
    full_desc,
    image_url,
    images,
    is_active,
    client,
    created_at
FROM projects
WHERE lang = $1 AND is_active = true
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetProjectBySlugAndLang :one
SELECT
    id,
    lang,
    slug,
    title,
    category,
    location,
    completed_year,
    short_desc,
    full_desc,
    image_url,
    images,
    is_active,
    client,
    created_at
FROM projects
WHERE slug = $1 and lang = $2 AND is_active = true;

-- name: GetProjectByID :one
SELECT 
    id,
    lang,
    slug,
    title,
    category,
    location,
    completed_year,
    short_desc,
    full_desc,
    image_url,
    images,
    is_active,
    client,
    created_at
FROM projects
WHERE id = $1;

-- name: CountActiveProjectsByLang :one
SELECT COUNT(*)
FROM projects
WHERE lang = $1 AND is_active = true;

-- name: ListProjectServices :many
SELECT 
    s.id,
    s.lang,
    s.slug,
    s.title,
    s.short_desc,
    s.image_url,
    s.created_at
FROM services s
JOIN project_services ps ON ps.service_id = s.id
WHERE ps.project_id = $1
ORDER BY ps.created_at ASC, s.id ASC;

-- name: ListProjectsByService :many
SELECT
    p.id,
    p.lang,
    p.slug,
    p.title,
    p.category,
    p.location,
    p.completed_year,
    p.short_desc,
    p.full_desc,
    p.image_url,
    p.images,
    p.is_active,
    p.client,
    p.created_at
FROM projects p 
JOIN project_services ps ON ps.project_id = p.id
WHERE ps.service_id = $1 AND p.lang = $2 AND p.is_active = true
ORDER BY p.created_at DESC;

-- name: CreateProject :one
INSERT INTO projects(
    lang,
    slug,
    title,
    category,
    location,
    completed_year,
    short_desc,
    full_desc,
    image_url,
    images,
    is_active,
    client
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING
    id,
    lang, 
    slug,
    title,
    category,
    location,
    completed_year,
    short_desc,
    full_desc,
    image_url,
    images,
    is_active,
    client,
    created_at;

-- name: UpdateProject :one
UPDATE projects
SET
    slug = $2,
    title = $3,
    category = $4,
    location = $5,
    completed_year = $6,
    short_desc = $7,
    full_desc = $8,
    image_url = $9,
    images = $10,
    is_active = $11,
    client = $12
WHERE id = $1
RETURNING
    id,
    lang,
    slug,
    title,
    category,
    location,
    completed_year,
    short_desc,
    full_desc,
    image_url,
    images,
    is_active,
    client,
    created_at;

-- name: DeleteProject :exec
DELETE FROM projects
WHERE id = $1;

-- name: AddServiceToProject :exec
INSERT INTO project_services(project_id, service_id)
VALUES($1, $2)
ON CONFLICT DO NOTHING;

-- name: DeleteServiceFromProject :exec
DELETE FROM project_services
WHERE project_id = $1 AND service_id = $2;

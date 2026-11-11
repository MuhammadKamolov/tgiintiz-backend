-- name: ListServicesByLang :many
SELECT
    id,
    lang,
    slug,
    title,
    short_desc,
    image_url,
    created_at    
FROM services
WHERE is_active = true AND lang = $1
ORDER BY created_at DESC;

-- name: GetServiceBySlugAndLang :one
SELECT 
    id,
    lang,
    slug,
    title,
    short_desc,
    image_url,
    standards,
    equipment,
    is_active,
    created_at
FROM services
WHERE slug = $1 AND lang = $2;

-- name: GetServiceByID :one
SELECT
    id,
    lang,
    slug,
    title,
    short_desc,
    full_desc,
    image_url,
    standards,
    equipment,
    is_active,
    created_at
FROM services
WHERE id = $1;

-- name: CreateService :one
INSERT INTO services (
    lang,
    slug,
    title,
    short_desc,
    full_desc,
    image_url,
    standards,
    equipment,
    is_active
) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING
    id,
    lang,
    slug,
    title,
    short_desc,
    full_desc,
    image_url,
    standards,
    equipment,
    is_active,
    created_at;

-- name: UpdateService :one
UPDATE services
SET
    title = $2,
    short_desc = $3,
    full_desc = $4,
    image_url = $5,
    standards = $6,
    equipment = $7,
    is_active = $8
WHERE id = $1
RETURNING 
    id,
    lang,
    slug,
    title,
    short_desc,
    full_desc,
    image_url,
    standards,
    equipment,
    is_active,
    created_at;

-- name: DeleteService :exec
DELETE FROM services
WHERE id = $1;

-- name: AddServiceToRequest :exec
INSERT INTO request_services (request_id, service_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: ListServicesByRequestID :many
SELECT
    s.id,
    s.lang,
    s.slug,
    s.title,
    s.short_desc,
    s.image_url
FROM services s
JOIN request_services rs ON rs.service_id = s.id
WHERE rs.request_id = $1
ORDER BY s.id ASC;

-- name: RemoveServicesByRequestID :exec
DELETE FROM request_services
WHERE request_id = $1;
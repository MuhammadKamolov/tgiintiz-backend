-- name: CreateSurveyRequest :one
INSERT INTO survey_requests(user_id, title, location, description)
VALUES ($1, $2, $3, $4)
RETURNING
    id,
    user_id,
    title,
    location,
    description,
    status,
    created_at,
    updated_at;

-- name: ListSurveyRequestsByUserID :many
SELECT
    id,
    user_id,
    title,
    location,
    description,
    status,
    created_at,
    updated_at
FROM survey_requests
WHERE user_id = $1 
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetSurveyRequestByID :one
SELECT
    id,
    user_id,
    title,
    location,
    description,
    status,
    created_at,
    updated_at
FROM survey_requests
WHERE id = $1;

-- name: ListSurveyRequests :many
SELECT
    id,
    user_id,
    title,
    location,
    description,
    status,
    created_at,
    updated_at
FROM survey_requests
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateSurveyRequestStatus :one
UPDATE survey_requests
SET status = $1, updated_at = NOW()
WHERE id = $2
RETURNING 
    id,
    user_id,
    title,
    location,
    description,
    status,
    created_at,
    updated_at;

-- name: DeleteSurveyRequest :exec
DELETE FROM survey_requests
WHERE id = $1;

-- name: ListSurveyRequestsByStatus :many
SELECT
    id,
    user_id,
    title,
    location,
    description,
    status,
    created_at,
    updated_at
FROM survey_requests
WHERE status = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

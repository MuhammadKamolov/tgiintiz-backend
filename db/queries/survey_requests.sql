-- name: CreateSurveyRequest :one
INSERT INTO survey_requests(user_id, title, location, description)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ListSurveyRequestsByUserID :many
SELECT *
FROM survey_requests
WHERE user_id = $1 ORDER BY created_at DESC;

-- name: GetSurveyRequestByID :one
SELECT *
FROM survey_requests
WHERE id = $1;

-- name: ListSurveyRequests :many
SELECT *
FROM survey_requests
ORDER BY created_at DESC;

-- name: UpdateSurveyRequestStatus :one
UPDATE survey_requests
SET status = $1, updated_at = NOW()
WHERE id = $2
RETURNING *;

-- name: DeleteSurveyRequest :exec
DELETE FROM survey_requests
WHERE id = $1;

-- name: ListSurveyRequestsByStatus :many
SELECT *
FROM survey_requests
WHERE status = $1 ORDER BY created_at DESC;

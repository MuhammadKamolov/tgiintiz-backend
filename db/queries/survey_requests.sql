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
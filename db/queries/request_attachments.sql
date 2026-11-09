-- name: CreateRequestAttachment :one
INSERT INTO request_attachments(request_id, file_name, file_url, file_size) 
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ListAttachmentsByRequestID :many
SELECT * FROM request_attachments
WHERE request_id = $1
ORDER BY created_at ASC;

-- name: GetAttachmentByID :one
SELECT * FROM request_attachments
WHERE id = $1;

-- name: DeleteRequestAttachment :exec
DELETE FROM request_attachments
WHERE id = $1;
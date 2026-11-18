-- name: CreateUser :one
INSERT INTO users(email, password_hash, role_id)
VALUES ($1, $2, $3)
RETURNING
    id,
    email,
    password_hash,
    created_at,
    role_id;

-- name: GetUserByEmail :one
SELECT
    id,
    email,
    password_hash,
    created_at,
    role_id
FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT
    id,
    email,
    password_hash,
    created_at,
    role_id
FROM users
WHERE id = $1;
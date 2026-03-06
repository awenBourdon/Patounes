-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
WHERE deleted_at IS NULL
ORDER BY created_at DESC;

-- name: CreateUser :one
INSERT INTO users (name, email, password)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET name = $2, email = $3, updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteUser :exec
UPDATE users
SET deleted_at = NOW()
WHERE id = $1;
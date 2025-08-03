-- name: CreateUser :one
INSERT INTO users (
    email, password, name, role, is_active, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: UpdateUser :one
UPDATE users SET
    email = COALESCE($2, email),
    password = COALESCE($3, password),
    name = COALESCE($4, name),
    role = COALESCE($5, role),
    is_active = COALESCE($6, is_active),
    updated_at = $7
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: ListUsers :many
SELECT * FROM users 
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;

-- name: ExistsByEmail :one
SELECT EXISTS(SELECT 1 FROM users WHERE email = $1);

-- name: ExistsByID :one
SELECT EXISTS(SELECT 1 FROM users WHERE id = $1); 
-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, login, password, email)
VALUES (
  $1,
  NOW(),
  NOW(),
  $2,
  $3,
  $4
)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

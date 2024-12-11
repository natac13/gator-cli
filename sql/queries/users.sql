-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
  $1,
  $2,
  $3,
  $4
)
RETURNING *;

-- name: GetUser :one
SELECT * from users
WHERE $1 = name;

-- name: GetUserById :one
SELECT * from users
WHERE $1 = id;

-- name: DeleteAllUsers :exec
DELETE from users;

-- name: GetUsers :many
SELECT * from users;

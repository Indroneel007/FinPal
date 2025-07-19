-- name: CreateUser :one
INSERT INTO users (
  username, hashed_password, full_name, email, salary
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;


-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: UpdateUserPassword :one
UPDATE users
SET hashed_password = $2
WHERE username = $1
RETURNING *;
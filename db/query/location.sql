-- name: CreateLocation :exec
INSERT INTO locations (username, address, latitude, longitude)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetLocationByUsername :one
SELECT * FROM locations WHERE username = $1 LIMIT 1;

-- name: UpdateLocationByUsername :exec
UPDATE locations
SET address = $2,
    latitude = $3,
    longitude = $4
WHERE username = $1
RETURNING *;
-- name: CreateGroup :one
INSERT INTO groups (
  group_name, currency, type
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: ListGroups :many
SELECT * FROM groups
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: ListGroupsByUser :many
SELECT g.* FROM groups g
JOIN accounts a ON g.id = a.group_id
WHERE a.owner = $1
ORDER BY g.id
LIMIT $2 OFFSET $3;

-- name: GetGroup :one
SELECT * FROM groups
WHERE id = $1 LIMIT 1;

-- name: GetGroupMembers :many
SELECT a.* FROM accounts a
JOIN groups g ON a.group_id = g.id
WHERE g.id = $1
ORDER BY a.id
LIMIT $2 OFFSET $3;

-- name: AcceptGroupInvitation :one
UPDATE accounts
SET has_accepted = true
WHERE id = $1 AND group_id = $2
RETURNING *;

-- name: DeleteGroupMember :exec
DELETE FROM accounts
WHERE id = $1 AND group_id = $2;

-- name: UpdateGroupName :one
UPDATE groups
SET group_name = $2
WHERE id = $1
RETURNING *;

-- name: DeleteGroup :exec
DELETE FROM groups
WHERE id = $1;

-- name: GetGroupTransactionHistory :many
SELECT
    t.id AS transfer_id,
    t.amount,
    fa.owner AS from_username,
    ta.owner AS to_username,
    t.created_at,
    t.group_id
FROM
    transfers t
    JOIN accounts fa ON t.from_account_id = fa.id
    JOIN accounts ta ON t.to_account_id = ta.id
WHERE
    t.group_id = $1
ORDER BY
    t.created_at DESC;
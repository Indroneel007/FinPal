-- name: CreateAccount :one
INSERT INTO accounts (
  owner, balance, currency, type
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: CreateAccountWithGroup :one
INSERT INTO accounts (
  owner, balance, currency, type, group_id, has_accepted
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: ListAccounts :many
SELECT * FROM accounts
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: ListAccountsByOwner :many
SELECT * FROM accounts
WHERE owner = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

-- name: GetAccountByGroupIDAndOwner :one
SELECT * FROM accounts
WHERE group_id = $1 AND owner = $2 LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: UpdateAcount :one
UPDATE accounts
  set balance = $2
WHERE id = $1
RETURNING *;

-- name: UpdateAccountGroup :one
UPDATE accounts
  set group_id = NULL
WHERE id = $1
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = $1;

-- name: GetAccountListByOwnerAndType :many
SELECT * FROM accounts
WHERE owner = $1 AND type = $2
ORDER BY id
LIMIT $3 OFFSET $4;

-- name: GetTotalByOwnerAndType :many
SELECT
  owner,
  type,
  SUM(balance) AS total_balance
FROM
  accounts
WHERE
  owner = $1
GROUP BY
  owner, type;
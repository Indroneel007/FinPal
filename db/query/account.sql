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

-- name: GetAccountByOwnerCurrencyType :one
SELECT * FROM accounts
WHERE owner = $1 AND currency = $2 AND type = $3 LIMIT 1;

-- name: GetAccountByOwnerCurrencyTypeGroupID :one
SELECT * FROM accounts
WHERE owner = $1 AND currency = $2 AND type = $3 AND group_id = $4 LIMIT 1;

-- name: ListTransactedUsersWithTotals :many
SELECT
    other_user::text AS username,
    COALESCE(SUM(CASE WHEN sub.from_account_id = sub.a_id THEN sub.amount END), 0)::bigint AS total_sent,
    COALESCE(SUM(CASE WHEN sub.to_account_id = sub.a_id THEN sub.amount END), 0)::bigint AS total_received
FROM (
    SELECT
        CASE
            WHEN t.from_account_id = a.id THEN a2.owner
            ELSE a.owner
        END AS other_user,
        t.from_account_id,
        t.to_account_id,
        t.amount,
        a.id AS a_id
    FROM transfers t
    JOIN accounts a ON a.owner = $1
    JOIN accounts a2 ON a2.id = t.to_account_id
    WHERE t.from_account_id = a.id OR t.to_account_id = a.id
) sub
GROUP BY other_user
ORDER BY other_user
LIMIT $2 OFFSET $3;
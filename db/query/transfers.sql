-- name: CreateTransfer :one
INSERT INTO transfers (from_account_id, to_account_id, amount, group_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetTransfer :one
SELECT * FROM transfers
WHERE id = $1 LIMIT 1;

-- name: ListTransfers :many
SELECT * FROM transfers
WHERE from_account_id = $1 OR to_account_id = $1
ORDER BY id DESC
LIMIT $2 OFFSET $3;

-- name: ListTransfersBetweenAccounts :many
SELECT 
  t.amount, 
  t.created_at, 
  t.from_account_id, 
  t.to_account_id,
  a.type
FROM transfers t
JOIN accounts a ON a.id = t.from_account_id
WHERE (t.from_account_id = ANY($1::bigint[]) AND t.to_account_id = ANY($2::bigint[]))
ORDER BY t.created_at DESC;
-- name: CreateAccount :one
INSERT INTO ACCOUNT (Owner, Balance, Currency, CreatedAt)
VALUES ($1, $2, $3, $4)
RETURNING *;


-- name: GetAccount :one
SELECT * FROM ACCOUNT
WHERE id = $1 LIMIT 1;

-- name: ListAccounts :many
SELECT * FROM ACCOUNT
WHERE owner = $1 
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: AddAccountBalance :one
update account 
set Balance = Balance + sqlc.arg(amount) 
where id = sqlc.arg(id) 
RETURNING *;

-- name: UpdateAccount :one
update account
set Balance = $2
where id = $1
RETURNING *;


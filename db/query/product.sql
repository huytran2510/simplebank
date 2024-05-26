-- name: CreateProduct :one
INSERT INTO Product (
    ProductId,
    ProductName,
    Price,
    StockQuantity
) VALUES (
    $1, $2, $3, $4
) 
RETURNING *;

-- name: GetProduct :one
SELECT * FROM Product
Where ProductId =$1 LIMIT 1;

-- name: UpdateProduct :one
UPDATE Product
SET ProductName = $1
WHERE ProductId = $2
RETURNING *;

-- name: DeleteProduct :one
DELETE 
FROM Product
WHERE ProductId = $1
RETURNING *;


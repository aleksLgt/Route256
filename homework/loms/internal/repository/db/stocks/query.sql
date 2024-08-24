-- name: GetStocks :many
SELECT * FROM stocks
WHERE id = any (sqlc.slice('ids'));

-- name: ReserveStock :exec
UPDATE stocks
SET reserved = $1
WHERE sku = $2;

-- name: GetStock :one
SELECT * FROM stocks
WHERE sku = $1
LIMIT 1;

-- name: RemoveReserveStock :exec
UPDATE stocks
SET reserved = $1 AND total_count = $2
WHERE id = $3;

-- name: CancelReserveStocks :exec
UPDATE stocks
SET reserved = $1
WHERE sku = any (sqlc.slice('skus'));

-- name: CreateStock :exec
INSERT INTO stocks (sku, total_count, reserved)
SELECT $1, $2, $3
WHERE NOT EXISTS (SELECT 1 FROM stocks WHERE sku = $1);

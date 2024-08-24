-- name: CreateOrder :one
INSERT INTO orders(user_id, status)
VALUES ($1, $2)
RETURNING id;

-- name: CreateOrderItem :exec
INSERT INTO order_items(order_id, sku, count)
VALUES ($1, $2, $3);

-- name: SetOrderStatus :exec
UPDATE orders
SET status = $1
WHERE id = $2;

-- name: GetOrder :one
SELECT * FROM orders
WHERE id = $1
LIMIT 1;

-- name: GetOrderItems :many
SELECT * FROM order_items
WHERE order_id = $1;

-- name: CreateOutboxOrderEvent :exec
INSERT INTO outbox_order_events(order_id, event_type)
VALUES ($1, $2);

-- name: GetUnsentOutboxOrderEvents :many
SELECT * FROM outbox_order_events
WHERE was_sent = false
LIMIT $1;

-- name: MarkAsSentOutboxOrderEvent :exec
UPDATE outbox_order_events
SET was_sent = true
WHERE id = $1;

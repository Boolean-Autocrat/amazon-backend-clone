-- name: CreateOrder :one
INSERT INTO orders (user_id, product_id, quantity) VALUES ($1, $2, $3) RETURNING *;

-- name: GetOrder :one
SELECT A.*, B.name AS product_name, B.price AS product_price FROM orders A INNER JOIN products B ON A.product_id = B.id WHERE A.id = $1;

-- name: GetOrders :many
SELECT * FROM orders WHERE user_id = $1;

-- name: ChangeOrderStatus :one
UPDATE orders SET status = $1 WHERE id = $2 RETURNING *;

-- name: UpdateProductStock :exec
UPDATE products SET stock = $1 WHERE id = $2;
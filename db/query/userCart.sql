-- name: CreateUserCart :one
INSERT INTO user_cart (user_id, product_id, quantity) VALUES ($1, $2, $3) RETURNING *;

-- name: GetUserCart :one
SELECT * FROM user_cart WHERE id = $1;

-- name: GetUserCarts :many
SELECT * FROM user_cart WHERE user_id = $1;

-- name: DeleteUserCart :exec
DELETE FROM user_cart WHERE id = $1;

-- name: UpdateUserCart :exec
UPDATE user_cart SET quantity = $1 WHERE id = $2;

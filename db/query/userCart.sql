-- name: AddToCart : one
INSERT INTO user_cart (user_id, product_id, quantity)
VALUES ($1, $2, $3)
RETURNING *;

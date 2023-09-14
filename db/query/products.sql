-- name: CreateProduct :one
INSERT INTO products 
(
    name, price, description, image, category, stock
) VALUES 
(
    $1 , $2 , $3 , $4 , $5 , $6
) RETURNING *;

-- name: GetProduct :one
SELECT * FROM products WHERE id = $1;

-- name: GetProducts :many
SELECT * FROM products LIMIT $1 OFFSET $2;

-- name: UpdateProduct :one
UPDATE products SET name = $1, price = $2, description = $3, image = $4, category = $5, stock = $6 WHERE id = $7 RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products WHERE id = $1;

-- name: GetProductByName :one
SELECT * FROM products WHERE name = $1;

-- name: SearchProducts :many
SELECT * FROM products WHERE name ILIKE $1 LIMIT $2 OFFSET $3;

-- name: ListProductCategories :many
SELECT DISTINCT category FROM products;

-- name: GetProductsByCategory :many
SELECT * FROM products WHERE category = $1 LIMIT $2 OFFSET $3;

-- name: AddImage :one
UPDATE products SET image = $1 WHERE id = $2 RETURNING *;

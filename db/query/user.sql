-- name: CreateUser :one
INSERT INTO users 
(
    username, password, email, phone_num
) VALUES 
(
    $1 , $2 , $3 , $4
) RETURNING ID, username, email, phone_num;

-- name: GetUser :one
SELECT ID, username, email, phone_num FROM users 
WHERE id = $1 LIMIT 1;

-- name: UpdateUser :exec
UPDATE users SET
    username = $1,
    email = $2,
    phone_num = $3
WHERE id = $4;

-- name: ChangePassword :exec
UPDATE users SET
    password = $1
WHERE id = $2;

-- name: GetPassword :one
SELECT password FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserID :one
SELECT id FROM users
WHERE username = $1 LIMIT 1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
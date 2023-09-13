// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0
// source: userCart.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createUserCart = `-- name: CreateUserCart :one
INSERT INTO user_cart (user_id, product_id, quantity) VALUES ($1, $2, $3) RETURNING id, user_id, product_id, quantity
`

type CreateUserCartParams struct {
	UserID    uuid.UUID `json:"userId"`
	ProductID uuid.UUID `json:"productId"`
	Quantity  int32     `json:"quantity"`
}

func (q *Queries) CreateUserCart(ctx context.Context, arg CreateUserCartParams) (UserCart, error) {
	row := q.db.QueryRowContext(ctx, createUserCart, arg.UserID, arg.ProductID, arg.Quantity)
	var i UserCart
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ProductID,
		&i.Quantity,
	)
	return i, err
}

const deleteUserCart = `-- name: DeleteUserCart :exec
DELETE FROM user_cart WHERE id = $1
`

func (q *Queries) DeleteUserCart(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteUserCart, id)
	return err
}

const getUserCart = `-- name: GetUserCart :one
SELECT id, user_id, product_id, quantity FROM user_cart WHERE id = $1
`

func (q *Queries) GetUserCart(ctx context.Context, id uuid.UUID) (UserCart, error) {
	row := q.db.QueryRowContext(ctx, getUserCart, id)
	var i UserCart
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ProductID,
		&i.Quantity,
	)
	return i, err
}

const getUserCarts = `-- name: GetUserCarts :many
SELECT id, user_id, product_id, quantity FROM user_cart WHERE user_id = $1
`

func (q *Queries) GetUserCarts(ctx context.Context, userID uuid.UUID) ([]UserCart, error) {
	rows, err := q.db.QueryContext(ctx, getUserCarts, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []UserCart
	for rows.Next() {
		var i UserCart
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.ProductID,
			&i.Quantity,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

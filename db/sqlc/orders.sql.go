// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0
// source: orders.sql

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const changeOrderStatus = `-- name: ChangeOrderStatus :one
UPDATE orders SET status = $1 WHERE id = $2 RETURNING id, status, user_id, product_id, quantity, created_at
`

type ChangeOrderStatusParams struct {
	Status string    `json:"status"`
	ID     uuid.UUID `json:"id"`
}

func (q *Queries) ChangeOrderStatus(ctx context.Context, arg ChangeOrderStatusParams) (Order, error) {
	row := q.db.QueryRowContext(ctx, changeOrderStatus, arg.Status, arg.ID)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.Status,
		&i.UserID,
		&i.ProductID,
		&i.Quantity,
		&i.CreatedAt,
	)
	return i, err
}

const createOrder = `-- name: CreateOrder :one
INSERT INTO orders (user_id, product_id, quantity) VALUES ($1, $2, $3) RETURNING id, status, user_id, product_id, quantity, created_at
`

type CreateOrderParams struct {
	UserID    uuid.UUID `json:"userId"`
	ProductID uuid.UUID `json:"productId"`
	Quantity  int32     `json:"quantity"`
}

func (q *Queries) CreateOrder(ctx context.Context, arg CreateOrderParams) (Order, error) {
	row := q.db.QueryRowContext(ctx, createOrder, arg.UserID, arg.ProductID, arg.Quantity)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.Status,
		&i.UserID,
		&i.ProductID,
		&i.Quantity,
		&i.CreatedAt,
	)
	return i, err
}

const getOrder = `-- name: GetOrder :one
SELECT a.id, a.status, a.user_id, a.product_id, a.quantity, a.created_at, B.name AS product_name, B.price AS product_price FROM orders A INNER JOIN products B ON A.product_id = B.id WHERE A.id = $1
`

type GetOrderRow struct {
	ID           uuid.UUID `json:"id"`
	Status       string    `json:"status"`
	UserID       uuid.UUID `json:"userId"`
	ProductID    uuid.UUID `json:"productId"`
	Quantity     int32     `json:"quantity"`
	CreatedAt    time.Time `json:"createdAt"`
	ProductName  string    `json:"productName"`
	ProductPrice int32     `json:"productPrice"`
}

func (q *Queries) GetOrder(ctx context.Context, id uuid.UUID) (GetOrderRow, error) {
	row := q.db.QueryRowContext(ctx, getOrder, id)
	var i GetOrderRow
	err := row.Scan(
		&i.ID,
		&i.Status,
		&i.UserID,
		&i.ProductID,
		&i.Quantity,
		&i.CreatedAt,
		&i.ProductName,
		&i.ProductPrice,
	)
	return i, err
}

const getOrders = `-- name: GetOrders :many
SELECT id, status, user_id, product_id, quantity, created_at FROM orders WHERE user_id = $1
`

func (q *Queries) GetOrders(ctx context.Context, userID uuid.UUID) ([]Order, error) {
	rows, err := q.db.QueryContext(ctx, getOrders, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Order
	for rows.Next() {
		var i Order
		if err := rows.Scan(
			&i.ID,
			&i.Status,
			&i.UserID,
			&i.ProductID,
			&i.Quantity,
			&i.CreatedAt,
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

const updateProductStock = `-- name: UpdateProductStock :exec
UPDATE products SET stock = $1 WHERE id = $2
`

type UpdateProductStockParams struct {
	Stock int32     `json:"stock"`
	ID    uuid.UUID `json:"id"`
}

func (q *Queries) UpdateProductStock(ctx context.Context, arg UpdateProductStockParams) error {
	_, err := q.db.ExecContext(ctx, updateProductStock, arg.Stock, arg.ID)
	return err
}

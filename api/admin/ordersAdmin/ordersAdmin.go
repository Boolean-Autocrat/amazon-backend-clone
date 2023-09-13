package ordersadmin

import (
	"context"
	"net/http"
	db "postman/amzn/db/sqlc"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Service struct {
	queries *db.Queries
}

func NewService(queries *db.Queries) *Service {
	return &Service{queries: queries}
}

func (s *Service) RegisterHandlers(router *gin.Engine) {
	router.GET("/admin/order/:id", s.AdminGetOrderHandler)
	router.GET("/admin/orders", s.AdminListOrdersHandler)
	router.PATCH("/admin/order/:id", s.AdminStatusOrderHandler)
}

type AdminOrderApi struct {
	ID           uuid.UUID `json:"id"`
	Status       string    `json:"status"`
	UserID       uuid.UUID `json:"userId"`
	ProductID    uuid.UUID `json:"productId"`
	Quantity     int32     `json:"quantity"`
	CreatedAt    time.Time `json:"createdAt"`
}

type AdminGetOrderApi struct {
	ID           uuid.UUID `json:"id"`
	Status       string    `json:"status"`
	UserID       uuid.UUID `json:"userId"`
	ProductID    uuid.UUID `json:"productId"`
	Quantity     int32     `json:"quantity"`
	CreatedAt    time.Time `json:"createdAt"`
	ProductName  string    `json:"productName"`
	ProductPrice int32     `json:"productPrice"`
}

type apiOrderStatusParams struct {
	Status string    `json:"status"`
	ID     uuid.UUID `json:"id"`
}

func fromDB(order db.GetOrderRow) *AdminGetOrderApi {
	return &AdminGetOrderApi{
		ID : order.ID,
		Status : order.Status,
		UserID : order.UserID,
		ProductID : order.ProductID,
		Quantity : order.Quantity,
		CreatedAt : order.CreatedAt,
		ProductName : order.ProductName,
		ProductPrice : order.ProductPrice,
	}
}

func fromDBList(order db.Order) *AdminOrderApi {
	return &AdminOrderApi{
		ID : order.ID,
		Status : order.Status,
		UserID : order.UserID,
		ProductID : order.ProductID,
		Quantity : order.Quantity,
		CreatedAt : order.CreatedAt,
	}
}

func (s *Service) AdminGetOrderHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := s.queries.GetOrder(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, fromDB(order))
}

func (s *Service) AdminListOrdersHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	orders, err := s.queries.GetOrders(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var ordersApi []*AdminOrderApi
	for _, order := range orders {
		ordersApi = append(ordersApi, fromDBList(order))
	}

	c.JSON(http.StatusOK, ordersApi)
}

func (s *Service) AdminStatusOrderHandler(c *gin.Context) {
	var req apiOrderStatusParams;
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	arg := db.ChangeOrderStatusParams{
		ID:     id,
		Status: req.Status,
	}

	order, err := s.queries.ChangeOrderStatus(context.Background(), arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}
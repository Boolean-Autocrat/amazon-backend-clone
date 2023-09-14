package orders

import (
	"context"
	"net/http"
	db "postman/amzn/db/sqlc"

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
	router.POST("/orders", s.CreateOrderHandler)
	router.GET("/orders", s.ListOrdersHandler)
	router.GET("/orders/:id", s.GetOrderHandler)
	router.PATCH("/orders/:id", s.UpdateOrderStatusHandler)
}

type apiOrder struct {
	UserID    uuid.UUID `json:"userId"`
	ProductID uuid.UUID `json:"productId"`
	Quantity  int32     `json:"quantity"`
}

type userID struct {
	UserID uuid.UUID `json:"userId"`
}

func (s *Service) CreateOrderHandler(c *gin.Context) {
	var req apiOrder
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Quantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "quantity must be positive"})
		return
	}

	_, err := s.queries.GetUser(context.Background(), req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}
	// verify product exists and check if product is in stock
	product, err := s.queries.GetProduct(context.Background(), req.ProductID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product not found"})
		return
	}
	productQuantity := product.Stock
	if productQuantity < req.Quantity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product out of stock"})
		return
	}


	arg := db.CreateOrderParams{
		UserID:    req.UserID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}

	order, err := s.queries.CreateOrder(context.Background(), arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (s *Service) ListOrdersHandler(c *gin.Context) {
	var req userID
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	orders, err := s.queries.GetOrders(context.Background(), req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (s *Service) GetOrderHandler(c *gin.Context) {
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

	c.JSON(http.StatusOK, order)
}

type apiOrderStatusParams struct {
	Status string    `json:"status"`
	ID     uuid.UUID `json:"id"`
}

func (s *Service) UpdateOrderStatusHandler(c *gin.Context) {
	var req apiOrderStatusParams
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
package orders

import (
	db "postman/amzn/db/sqlc"

	"github.com/gin-gonic/gin"
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
	router.PUT("/orders/:id", s.UpdateOrderHandler)
	router.DELETE("/orders/:id", s.DeleteOrderHandler)
}

type apiOrder struct {}

func (s *Service) CreateOrderHandler(c *gin.Context) {}

func (s *Service) ListOrdersHandler(c *gin.Context) {}

func (s *Service) GetOrderHandler(c *gin.Context) {}

func (s *Service) UpdateOrderHandler(c *gin.Context) {}

func (s *Service) DeleteOrderHandler(c *gin.Context) {}
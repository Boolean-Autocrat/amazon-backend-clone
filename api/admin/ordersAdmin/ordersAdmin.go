package ordersadmin

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
	router.GET("/admin/order/:id", s.GetOrderHandler)
	router.DELETE("/admin/order/:id", s.DeleteOrderHandler)
}

func (s *Service) GetOrderHandler(c *gin.Context) {}

func (s *Service) DeleteOrderHandler(c *gin.Context) {}
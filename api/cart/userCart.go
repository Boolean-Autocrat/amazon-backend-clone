package cart

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
	router.POST("/cart", s.CreateCartHandler)
	router.GET("/cart/:id", s.GetCartHandler)
	router.PATCH("/cart/:id", s.UpdateCartHandler)
}

type apiCart struct{}

func (s *Service) CreateCartHandler(c *gin.Context) {}

func (s *Service) GetCartHandler(c *gin.Context) {}

func (s *Service) UpdateCartHandler(c *gin.Context) {}
package usersadmin

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
	router.GET("/admin/user/:id", s.AdminUserHandler)
	router.DELETE("/admin/user/:id", s.DeleteUserHandler)
}

func (s *Service) AdminUserHandler(c *gin.Context) {}

func (s *Service) DeleteUserHandler(c *gin.Context) {}
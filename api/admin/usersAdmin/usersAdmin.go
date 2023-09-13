package usersadmin

import (
	"context"
	"database/sql"
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
	router.GET("/admin/user/:id", s.AdminUserHandler)
	router.GET("/admin/username/:username", s.AdminGetUserByUsernameHandler)
	router.DELETE("/admin/user/:id", s.DeleteUserHandler)
}

func (s *Service) AdminGetUserByUsernameHandler(c *gin.Context) {
	username := c.Param("username")
	user, err := s.queries.GetUserID(context.Background(), username)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (s *Service) AdminUserHandler(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	user, err := s.queries.GetUser(context.Background(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (s *Service) DeleteUserHandler(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err := s.queries.DeleteUser(context.Background(), id); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "pq: update or delete on table \"users\" violates foreign key constraint \"orders_user_id_fkey\" on table \"orders\"" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user has pending orders"})
			return
		}

		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}
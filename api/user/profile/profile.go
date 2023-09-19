package profile

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
	router.GET("/profile", s.GetUser)
	router.PUT("/profile/edit", s.UpdateUser)
	router.DELETE("/profile", s.DeleteUser)
}

type returnUser struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	PhoneNum string    `json:"phoneNum"`
}

func fromGetDB(user db.GetUserRow) *returnUser {
	return &returnUser{
		ID: 	   user.ID,
		Username:  user.Username,
		Email:     user.Email,
		PhoneNum:  user.PhoneNum,
	}
}

func (s *Service) GetUser(c *gin.Context) {
	// Parse request
	idStr, _ := c.Get("userID")
	id, err := uuid.Parse(idStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user by username
	user, err := s.queries.GetUser(context.Background(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	// Build response
	response := fromGetDB(user)
	c.IndentedJSON(http.StatusOK, response)
}

func (s *Service) UpdateUser(c *gin.Context) {
	// Parse request
	idStr, _ := c.Get("userID")
	id, err := uuid.Parse(idStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var request returnUser
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ValidateUpdateRequest(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update user
	params := db.UpdateUserParams{
		Username: request.Username,
		Email:    request.Email,
		PhoneNum: request.PhoneNum,
		ID:       id,
	}
	if err := s.queries.UpdateUser(context.Background(), params); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	// Build response
	c.Status(http.StatusOK)
}

func (s *Service) DeleteUser(c *gin.Context) {
	idStr, _ := c.Get("userID")
	id, err := uuid.Parse(idStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Delete user
	if err := s.queries.DeleteUser(context.Background(), id); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	// return 200 OK
	c.Status(http.StatusOK)
}
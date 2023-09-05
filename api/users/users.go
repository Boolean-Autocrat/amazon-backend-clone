package users

import (
	"context"
	"database/sql"
	"net/http"
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
	router.POST("/user/register", s.CreateUser)
	router.GET("/user/:id", s.GetUser)
	router.PUT("/user/:id", s.UpdateUser)
	router.DELETE("/user/:id", s.DeleteUser)
	router.POST("/user/login", s.LoginUser)
	router.POST("/user/logout", s.LogoutUser)
	router.POST("/user/changepwd", s.ChangePassword)
}

type apiUser struct {
	ID        int64          `json:"id"`
	Username  string         `json:"username"`
	Password  string         `pg:"-" binding:"required,min=7,max=32"`
    Email     string		 `json:"email"`
	PhoneNum  string         `json:"phoneNum"`
}

type loginUser struct {
	Username  string         `json:"username"`
	Password  string         `json:"password"`
}

type returnUser struct {
	ID        int64          `json:"id"`
	Username  string         `json:"username"`
	Email     string		 `json:"email"`
	PhoneNum  string         `json:"phoneNum"`
}

type changePwd struct {
	ID        int64          `json:"id"`
	Password  string         `json:"password"`
	OldPwd    string		 `json:"oldPwd"`
}

func fromCreateDB(user db.CreateUserRow) *returnUser {
	return &returnUser{
		ID: 	   user.ID,
		Username:  user.Username,
		Email:     user.Email,
		PhoneNum:  user.PhoneNum,
	}
}

func fromGetDB(user db.GetUserRow) *returnUser {
	return &returnUser{
		ID: 	   user.ID,
		Username:  user.Username,
		Email:     user.Email,
		PhoneNum:  user.PhoneNum,
	}
}

type pathParameters struct {
	ID int64 `uri:"id" binding:"required"`
}

func (s *Service) CreateUser(c *gin.Context) {
	// Parse request
	var request apiUser
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate request
	if err := ValidateUserRequest(request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	hash := hashAndSalt([]byte(request.Password))

	if hash == "Error hashing password" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error hashing password"})
		return
	}

	// Create user
	params := db.CreateUserParams{
		Username: request.Username,
		Password: hash,
		Email: request.Email,
		PhoneNum: request.PhoneNum,
	}
	user, err := s.queries.CreateUser(context.Background(), params)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	// Build response
	response := fromCreateDB(user)
	c.IndentedJSON(http.StatusCreated, response)
}

func (s *Service) GetUser(c *gin.Context) {
	// Parse request
	var pathParams pathParameters
	if err := c.ShouldBindUri(&pathParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user by username
	user, err := s.queries.GetUser(context.Background(), pathParams.ID)
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

func (s *Service) DeleteUser(c *gin.Context) {
	// Parse request
	var pathParams pathParameters
	if err := c.ShouldBindUri(&pathParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Delete user
	if err := s.queries.DeleteUser(context.Background(), pathParams.ID); err != nil {
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

func (s *Service) UpdateUser(c *gin.Context) {
	// Parse request
	var pathParams pathParameters
	if err := c.ShouldBindUri(&pathParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var request apiUser
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ValidateUserRequest(request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	// Update user
	params := db.UpdateUserParams{
		Username: request.Username,
		Email:    request.Email,
		PhoneNum: request.PhoneNum,
		ID:   pathParams.ID,
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

func (s *Service) LoginUser(c *gin.Context) {
	var request loginUser
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := s.queries.GetUserID(context.Background(), request.Username)
	currentPassword, _ := s.queries.GetPassword(context.Background(), userID)

	if !comparePasswords(currentPassword, []byte(request.Password)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect username or password"})
		return
	}

	// Build response
	c.Status(http.StatusOK)
}

func (s *Service) LogoutUser(c *gin.Context) {
}

func (s *Service) ChangePassword(c *gin.Context) {
		// Parse request
	var request changePwd
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	previousPassword, _ := s.queries.GetPassword(context.Background(), request.ID)

	if !comparePasswords(previousPassword, []byte(request.OldPwd)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect old password"})
		return
	}
	hash := hashAndSalt([]byte(request.Password))

	params := db.ChangePasswordParams{
		Password: hash,
		ID: request.ID,
	}
	if err := s.queries.ChangePassword(context.Background(), params); err != nil {
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
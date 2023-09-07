package userAuth

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
	router.POST("/auth/register", s.CreateUser)
	router.POST("/auth/login", s.LoginUser)
	router.POST("/auth/logout", s.LogoutUser)
	router.POST("/auth/changepassword", s.ChangePassword)

	router.GET("/profile/:id", s.GetUser)
	router.PUT("/profile/edit/:id", s.UpdateUser)
	router.DELETE("/profile/:id", s.DeleteUser)
}

type apiUser struct {
	ID        uuid.UUID      `json:"id"`
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
	ID        uuid.UUID      `json:"id"`
	Username  string         `json:"username"`
	Email     string		 `json:"email"`
	PhoneNum  string         `json:"phoneNum"`
}

type changePwd struct {
	ID        uuid.UUID      `json:"id"`
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

func (s *Service) CreateUser(c *gin.Context) {
	var request apiUser
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Server side validation (client side validation lol)
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
	idStr := c.Param("id")
    id, err := uuid.Parse(idStr)
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

func (s *Service) DeleteUser(c *gin.Context) {
    // explicit conversion from string to uuid.UUID
    idStr := c.Param("id")
    id, err := uuid.Parse(idStr)
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

func (s *Service) UpdateUser(c *gin.Context) {
	// Parse request
	idStr := c.Param("id")
    id, err := uuid.Parse(idStr)
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
		ID:   id,
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
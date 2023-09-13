package productsadmin

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
	router.POST("/admin/products/create", s.CreateProductHandler)
	router.PUT("/admin/product/:id", s.UpdateProductHandler)
	router.DELETE("/admin/product/:id", s.DeleteProductHandler)
}

type apiProduct struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name" binding:"required"`
	Price       int32     `json:"price" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Image       string    `json:"image" binding:"required"`
	Category    string    `json:"category" binding:"required"`
	Stock       int32     `json:"stock" binding:"required"`
}

func fromDB(product db.Product) *apiProduct {
	return &apiProduct{
		ID:          product.ID,
		Name:        product.Name,
		Price:       product.Price,
		Description: product.Description,
		Image:       product.Image,
		Category:    product.Category,
		Stock:       product.Stock,
	}
}

func (s *Service) CreateProductHandler(c *gin.Context) {
	var request apiProduct

	if err := c.ShouldBindJSON(&request); err != nil {
		if err := request.ValidateProductUpdateRequest(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := request.ValidateProductUpdateRequest(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := s.queries.CreateProduct(context.Background(), db.CreateProductParams{
		Name:        request.Name,
		Price:       request.Price,
		Description: request.Description,
		Image:       request.Image,
		Category:    request.Category,
		Stock:       request.Stock,
	})
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	response := fromDB(product)
	c.IndentedJSON(http.StatusCreated, response)
}

func (s *Service) UpdateProductHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var request apiProduct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := s.queries.UpdateProduct(context.Background(), db.UpdateProductParams{
		ID:          id,
		Name:        request.Name,
		Price:       request.Price,
		Description: request.Description,
		Image:       request.Image,
		Category:    request.Category,
		Stock:       request.Stock,
	})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	response := fromDB(product)
	c.IndentedJSON(http.StatusOK, response)
}

func (s *Service) DeleteProductHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	delerr := s.queries.DeleteProduct(context.Background(), id)
	if delerr != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
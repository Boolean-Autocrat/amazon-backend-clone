package products

import (
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
	router.POST("/product/create", s.CreateProduct)
	router.GET("/product/:id", s.GetProduct)
	router.PUT("/product/:id", s.UpdateProduct)
	router.DELETE("/product/:id", s.DeleteProduct)
}

type apiProduct struct {
	Name        string `json:"name"`
	Price       int32  `json:"price"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Category    string `json:"category"`
	Stock       int32  `json:"stock"`
}

func fromDB(product db.Product) *apiProduct {
	return &apiProduct{
		Name:        product.Name,
		Price:       product.Price,
		Description: product.Description,
		Image:       product.Image,
		Category:    product.Category,
		Stock:       product.Stock,
	}
}

type pathParameters struct {
	ID int64 `uri:"id" binding:"required"`
}

func (s *Service) CreateProduct(c *gin.Context) {
	var request apiProduct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := s.queries.CreateProduct(c, db.CreateProductParams{
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

func (s *Service) GetProduct(c *gin.Context) {
	var params pathParameters
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := s.queries.GetProduct(c, params.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	response := fromDB(product)
	c.IndentedJSON(http.StatusOK, response)
}

func (s *Service) UpdateProduct(c *gin.Context) {
	var params pathParameters
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var request apiProduct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := s.queries.UpdateProduct(c, db.UpdateProductParams{
		ID:          params.ID,
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

func (s *Service) DeleteProduct(c *gin.Context) {
	var params pathParameters
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.queries.DeleteProduct(c, params.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
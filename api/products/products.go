package products

import (
	"database/sql"
	"net/http"
	db "postman/amzn/db/sqlc"
	"strconv"

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
	router.GET("/products", s.ListProductsHandler)
	router.GET("/products/:id", s.GetProductHandler)
	router.GET("/products/search", s.SearchProductsHandler)
	router.GET("/products/categories", s.ListProductCategoriesHandler)
	router.GET("/products/category/:category", s.GetProductCategoryHandler)

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
	Category    string 	  `json:"category" binding:"required"`
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

type pathParameters struct {
	ID uuid.UUID `uri:"id" binding:"required"`
}

func (s *Service) ListProductsHandler(c *gin.Context) {
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.Atoi(c.Query("offset"))
	if limit == 0 || offset == 0 {
		limit = 10
		offset = 0
	}
	params := db.GetProductsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}
	products, err := s.queries.GetProducts(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response []apiProduct
	for _, product := range products {
		response = append(response, *fromDB(product))
	}

	c.JSON(http.StatusOK, response)
}

func (s *Service) SearchProductsHandler(c *gin.Context) {
	name := c.Query("q")
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.Atoi(c.Query("offset"))
	if limit == 0 || offset == 0 {
		limit = 10
		offset = 0
	}
	params := db.SearchProductsParams{
		Name: "%" + name + "%",
		Limit: int32(limit),
		Offset: int32(offset),
	}
	products, err := s.queries.SearchProducts(c, params)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response []apiProduct
	for _, product := range products {
		response = append(response, *fromDB(product))
	}
	c.JSON(http.StatusOK, response)
}

func (s *Service) ListProductCategoriesHandler(c *gin.Context) {}

func (s *Service) GetProductCategoryHandler(c *gin.Context) {}

func (s *Service) CreateProductHandler(c *gin.Context) {
	var request apiProduct
	

	if err := c.ShouldBindJSON(&request); err != nil {
		if err := request.ValidateProductRequest(); err != nil {
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

func (s *Service) GetProductHandler(c *gin.Context) {
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

	product, err := s.queries.UpdateProduct(c, db.UpdateProductParams{
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

	delerr := s.queries.DeleteProduct(c, id)
	if delerr != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
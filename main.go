package main

import (
	"log"
	"os"
	ordersadmin "postman/amzn/api/admin/ordersAdmin"
	productsadmin "postman/amzn/api/admin/productsAdmin"
	usersadmin "postman/amzn/api/admin/usersAdmin"
	"postman/amzn/api/orders"
	"postman/amzn/api/products"
	"postman/amzn/api/user/auth"
	"postman/amzn/api/user/profile"
	db "postman/amzn/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }
}

func main() {
	postgres, err := db.NewPostgres(os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"))
	if err != nil {
		log.Fatal(err.Error())
	}

	queries := db.New(postgres.DB)
	
	// Instantiating services

	// Admin Services
	adminOrdersService := ordersadmin.NewService(queries)
	adminUsersService := usersadmin.NewService(queries)
	adminProductsService := productsadmin.NewService(queries)
	
	// Auth + Profile
	authService := auth.NewService(queries)
	profileService := profile.NewService(queries) 
	
	// Products
	productService := products.NewService(queries)

	// Orders
	ordersService := orders.NewService(queries)

	// Registering service handlers to the Gin router
	router := gin.Default()
	adminOrdersService.RegisterHandlers(router)
	adminUsersService.RegisterHandlers(router)
	adminProductsService.RegisterHandlers(router)

	authService.RegisterHandlers(router)
	profileService.RegisterHandlers(router)
	
	productService.RegisterHandlers(router)

	ordersService.RegisterHandlers(router)

	router.Run()
}
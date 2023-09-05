package main

import (
	"log"
	"os"
	"postman/amzn/api/users"
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

	// Instantiate the user service
	queries := db.New(postgres.DB)
	authorService := users.NewService(queries)

	// Register our service handlers to the router
	router := gin.Default()
	authorService.RegisterHandlers(router)

	// Start the server
	router.Run()
}
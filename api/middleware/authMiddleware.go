package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.FullPath() == "/auth/register" || c.FullPath() == "/auth/login" {
			c.Next()
			return
		}
		secretKey := "rust-goat"
		tokenStr, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			fmt.Println(err)
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, nil
			}
		
			return []byte(secretKey), nil
		})
		
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// check if token is expired
			expired := claims["exp"].(float64)
			if expired < float64(time.Now().Unix()) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
				fmt.Println(err)
				c.Abort()
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			fmt.Println(err)
			c.Abort()
			return
		}
		c.Next()
	}
}

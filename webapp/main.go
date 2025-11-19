package main

import (
	"cs3380/database"
	"cs3380/middleware"
	"cs3380/routes"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Hello, World!")

	if err := database.ConnectDatabase(); err != nil {
		fmt.Println("Failed to connect to database:", err)
		return
	}

	fmt.Println("Connected to database successfully!")

	r := gin.Default()

	apiGroup := r.Group("/api")
	{
		apiGroup.POST("/register", routes.RegisterUser)
		apiGroup.GET("/login", routes.LoginUser)

		secureGroup := apiGroup.Group("/secure")
		{
			// Middleware for authentication
			secureGroup.Use(middleware.AuthorizeRequest)
			secureGroup.GET("/ping", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "pong",
				})
			})
		}
	}
	if err := r.Run(); err != nil {
		log.Fatalf("Failed to start server: %s", err.Error())
	}
}

package main

import (
	"cs3380/database"
	"cs3380/middleware"
	"cs3380/routes"
	"fmt"
	"log"

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
	gin.SetMode(gin.DebugMode)

	apiGroup := r.Group("/api")
	{
		apiGroup.POST("/register", routes.RegisterUser)
		apiGroup.GET("/login", routes.LoginUser)

		secureGroup := apiGroup.Group("/secure")
		{
			// Middleware for authentication
			secureGroup.Use(middleware.AuthorizeRequest)
			secureGroup.POST("/create-family", routes.CreateFamily)
			secureGroup.POST("/join-family", routes.JoinFamily)
			secureGroup.POST("/leave-family", routes.LeaveFamily)
			secureGroup.GET("/get-family", routes.GetFamily)

			secureGroup.POST("/create-publisher", routes.CreatePublisher)
			secureGroup.POST("/join-publisher", routes.JoinPublisher)
			secureGroup.POST("/leave-publisher", routes.LeavePublisher)
		}
	}
	if err := r.Run(); err != nil {
		log.Fatalf("Failed to start server: %s", err.Error())
	}
}

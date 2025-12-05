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

	// API routes first
	apiGroup := r.Group("/api")
	{
		apiGroup.POST("/register", routes.RegisterUser)
		apiGroup.POST("/login", routes.LoginUser)

		secureGroup := apiGroup.Group("/secure")
		{
			// Middleware for authentication
			secureGroup.Use(middleware.AuthorizeRequest)
			secureGroup.POST("/create-family", routes.CreateFamily)
			secureGroup.POST("/join-family", routes.JoinFamily)
			secureGroup.POST("/leave-family", routes.LeaveFamily)
			secureGroup.GET("/get-family", routes.GetFamily)
			secureGroup.POST("/add-user-to-family", routes.AddUserToFamily)
			secureGroup.GET("/get-family-games", routes.GetFamilyGames)

			secureGroup.POST("/add-to-cart", routes.AddToCart)
			secureGroup.POST("/remove-from-cart", routes.RemoveFromCart)
			secureGroup.GET("/get-cart", routes.GetCartItems)

			secureGroup.POST("/purchase", routes.PurchaseItems)
			secureGroup.POST("/request-refund", routes.RequestRefund)
			secureGroup.POST("/approve-refund", routes.ApproveRefund)

			secureGroup.POST("/add-to-wishlist", routes.AddToWishlist)
			secureGroup.POST("/remove-from-wishlist", routes.RemoveFromWishlist)

			secureGroup.POST("/create-publisher", routes.CreatePublisher)
			secureGroup.POST("/join-publisher", routes.JoinPublisher)
			secureGroup.POST("/leave-publisher", routes.LeavePublisher)

			secureGroup.POST("/create-game", routes.CreateGame)
			secureGroup.POST("/remove-game", routes.RemoveGame)
			secureGroup.GET("/find-games", routes.FindGames)
			secureGroup.GET("/library", routes.GetLibrary)
			secureGroup.GET("/game-discounts/:gameId", routes.GetGameDiscounts)
		}
	}

	if err := r.Run(); err != nil {
		log.Fatalf("Failed to start server: %s", err.Error())
	}
}

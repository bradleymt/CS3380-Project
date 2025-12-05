package routes

import (
	"cs3380/database"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// POST /api/secure/create-game
func CreateGame(c *gin.Context) {
	var requestBody database.Game
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(requestBody.Name) == 0 ||
		len(requestBody.Currency) == 0 ||
		requestBody.Price < 0 ||
		len(requestBody.Genre) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body supplied"})
		return
	}

	var user database.User
	if userID, ok := c.Get("user_id"); !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find associated user"})
		return
	} else {
		if err := database.DB.First(&user, "id = ?", userID.(uint64)).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find associated user"})
			return
		}
	}

	if user.PublisherID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user does not belong to a publisher"})
		return
	}

	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		if res := tx.Model(&database.Publisher{}).Where("id = ?", user.PublisherID).First(nil); res.RowsAffected != 1 {
			return fmt.Errorf("failed to find user publisher")
		}

		if res := tx.Model(&database.Game{}).Where("publisher_id = ? AND name = ?", user.PublisherID, requestBody.Name).First(nil); res.RowsAffected != 0 {
			return fmt.Errorf("publisher cannot publish the same game")
		}

		requestBody.PublisherID = *user.PublisherID
		return tx.Create(&requestBody).Error
	}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to create game: %v", err.Error())})
		return
	}
}

// POST /api/secure/remove-game
func RemoveGame(c *gin.Context) {
	var requestBody struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(requestBody.Name) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body supplied"})
		return
	}

	var user database.User
	if userID, ok := c.Get("user_id"); !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find associated user"})
		return
	} else {
		if err := database.DB.First(&user, "id = ?", userID.(uint64)).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find associated user"})
			return
		}
	}

	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		var game database.Game
		if err := tx.Find(&game, "publisher_id = ? AND name = ?", user.PublisherID, requestBody.Name).Error; err != nil {
			return fmt.Errorf("failed to find game with publisher and name supplied")
		}

		return tx.Delete(&game).Error
	}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to remove game"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "remove game"})
}

// add discount
// POST /api/secure/create-discount

// GET /api/secure/find-games?name=&publisher=&currency=&price=&genre=
func FindGames(c *gin.Context) {
	name := c.Query("name")
	publisher := c.Query("publisher")
	currency := c.Query("currency")
	priceStr := c.Query("price")
	genre := c.Query("genre")

	db := database.DB.Model(&database.Game{})

	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}

	if currency != "" {
		db = db.Where("currency = ?", currency)
	}

	if genre != "" {
		db = db.Where("genre = ?", genre)
	}

	if priceStr != "" {
		var price float64
		if _, err := fmt.Sscanf(priceStr, "%f", &price); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid price"})
			return
		}
		db = db.Where("price = ?", price)
	}

	if publisher != "" {
		var pub database.Publisher
		if err := database.DB.Where("studio_name = ?", publisher).First(&pub).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to lookup publisher"})
			return
		}
		db = db.Where("publisher_id = ?", pub.ID)
	}

	var games []database.Game
	if err := db.Find(&games).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query games"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"games": games})
}

// GET /api/secure/library
func GetLibrary(c *gin.Context) {
	var user database.User
	if userID, ok := c.Get("user_id"); !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find associated user"})
		return
	} else {
		if err := database.DB.First(&user, "id = ?", userID.(uint64)).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find associated user"})
			return
		}
	}

	// Get all purchases for this user with game info
	var purchases []database.Purchase
	if err := database.DB.Preload("Game").Where("user_id = ?", user.ID).Find(&purchases).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query purchases"})
		return
	}

	// Extract games from purchases
	var games []database.Game
	for _, purchase := range purchases {
		games = append(games, purchase.Game)
	}

	c.JSON(http.StatusOK, gin.H{"games": games})
}

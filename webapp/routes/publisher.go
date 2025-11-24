package routes

import (
	"cs3380/database"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// POST /api/secure/create-publisher
func CreatePublisher(c *gin.Context) {
	var requestBody database.Publisher
	if err := c.ShouldBindJSON(&requestBody); err != nil {
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

	if user.PublisherID != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user already belongs to a publisher"})
		return
	}

	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Model(&database.Publisher{}).Where("studio_name = ?", requestBody.StudioName).Count(&count).Error; err != nil {
			return fmt.Errorf("failed to check for existing publisher: %v", err)
		}
		if count > 0 {
			return fmt.Errorf("publisher with that name already exists")
		}
		if err := tx.Create(&requestBody).Error; err != nil {
			return fmt.Errorf("failed to create publisher")
		}

		if err := tx.Model(&database.User{}).Where("id = ?", user.ID).Update("publisher_id", requestBody.ID).Error; err != nil {
			return fmt.Errorf("failed to update publisher column for user: %v", err)
		}

		return nil
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to create publisher: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "added user to new publisher"})
}

// POST /api/secure/leave-publisher
func LeavePublisher(c *gin.Context) {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "user is not part of a publisher"})
		return
	}

	if err := database.DB.Model(&database.User{}).Where("id = ?", user.ID).Update("publisher_id", nil).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update publisher information for user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "removed user from publisher"})
}

// POST /api/secure/join-publisher
func JoinPublisher(c *gin.Context) {
	var requestBody struct {
		StudioName string `json:"studio_name" binding:"required,min=8,max=32"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
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

	if user.PublisherID != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user is already part of another publisher"})
		return
	}

	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		var publisher database.Publisher
		if result := tx.Find(&publisher, "studio_name = ?", requestBody.StudioName); result.RowsAffected == 0 {
			return fmt.Errorf("publisher not found")
		}

		if err := tx.Model(&database.User{}).Where("id = ?", user.ID).Update("publisher_id", publisher.ID).Error; err != nil {
			return fmt.Errorf("failed to update publisher column for user: %v", err)
		}

		return nil
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to add user to publisher: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "added user to publisher"})
}

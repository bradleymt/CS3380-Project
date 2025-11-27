package routes

import (
	"cs3380/database"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// POST /api/secure/add-to-cart
func AddToCart(c *gin.Context) {
	var requestBody struct {
		GameID uint `json:"game_id"`
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

	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		if res := tx.Model(&database.Game{}).Where("id = ?", requestBody.GameID).Find(nil); res.RowsAffected != 1 {
			return fmt.Errorf("failed to find game")
		}

		var cartItem database.CartItem
		cartItem.GameID = requestBody.GameID
		cartItem.UserID = user.ID

		return tx.Create(&cartItem).Error
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to add item to cart: %v", err)})
	}

	c.JSON(http.StatusOK, gin.H{"result": "Added item to cart!"})
}

// POST /api/secure/remove-from-cart
func RemoveFromCart(c *gin.Context) {
	var requestBody struct {
		GameID uint `json:"game_id"`
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

	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		var cartItem database.CartItem
		if err := tx.Find(&cartItem, "user_id = ? AND game_id = ?", user.ID, requestBody.GameID).Error; err != nil {
			return fmt.Errorf("failed to find cart item")
		}

		return tx.Delete(&cartItem).Error
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to remove item from cart: %v", err)})
	}

	c.JSON(http.StatusOK, gin.H{"result": "removed item to cart!"})
}

// GET /api/secure/get-cart
func GetCartItems(c *gin.Context) {
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

	var result []database.Game
	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		var cartItem []database.CartItem
		if err := tx.Preload("Game").Find(&cartItem, "user_id = ?", user.ID).Error; err != nil {
			return err
		}

		for _, item := range cartItem {
			result = append(result, item.Game)
		}

		return nil
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to find cart items: %v", err)})
	}

	c.JSON(http.StatusOK, gin.H{"cart_items": result})
}

// POST /api/secure/purchase
func PurchaseItems(c *gin.Context) {
	var requestBody struct {
		GameIDs       []uint `json:"game_ids"`
		Street        string `json:"street"`
		City          string `json:"city"`
		State         string `json:"state"`
		ZipCode       string `json:"zip_code"`
		Country       string `json:"country"`
		PaymentMethod string `json:"payment_method" binding:"required"`
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

	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		var cartItems []database.CartItem
		var err error
		if len(requestBody.GameIDs) > 0 {
			err = tx.Preload("Game").Find(&cartItems, "user_id = ? AND game_id IN ?", user.ID, requestBody.GameIDs).Error
		} else {
			err = tx.Preload("Game").Find(&cartItems, "user_id = ?", user.ID).Error
		}
		if err != nil {
			return err
		}

		if len(cartItems) == 0 {
			return fmt.Errorf("no items to purchase")
		}

		for _, item := range cartItems {
			purchase := database.Purchase{
				UserID:        user.ID,
				GameID:        item.GameID,
				Street:        requestBody.Street,
				City:          requestBody.City,
				State:         requestBody.State,
				ZipCode:       requestBody.ZipCode,
				Country:       requestBody.Country,
				PaymentMethod: requestBody.PaymentMethod,
			}

			if err := tx.Create(&purchase).Error; err != nil {
				return err
			}

			// add to library
			libItem := database.LibraryItem{UserID: user.ID, GameID: item.GameID}
			if err := tx.Create(&libItem).Error; err != nil {
				return err
			}

			// remove from cart
			if err := tx.Delete(&database.CartItem{}, item.ID).Error; err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to complete purchase: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "Purchase completed"})
}

// POST /api/secure/request-refund
func RequestRefund(c *gin.Context) {
	var requestBody struct {
		PurchaseID uint   `json:"purchase_id" binding:"required"`
		Reason     string `json:"reason"`
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

	var purchase database.Purchase
	if err := database.DB.First(&purchase, "id = ? AND user_id = ?", requestBody.PurchaseID, user.ID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid purchase id supplied"})
		return
	}

	refund := database.Refund{PurchaseID: purchase.ID, Reason: requestBody.Reason, Approved: false}
	if err := database.DB.Create(&refund).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to create refund request: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "refund requested", "refund_id": refund.ID})
}

// POST /api/secure/approve-refund
func ApproveRefund(c *gin.Context) {
	var requestBody struct {
		RefundID uint `json:"refund_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body supplied"})
		return
	}

	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		var refund database.Refund
		if err := tx.Preload("Purchase").First(&refund, "id = ?", requestBody.RefundID).Error; err != nil {
			return err
		}

		if refund.Approved {
			return fmt.Errorf("refund already approved")
		}

		// mark refund approved
		if err := tx.Model(&refund).Update("approved", true).Error; err != nil {
			return err
		}

		// remove purchase
		if err := tx.Delete(&database.Purchase{}, refund.PurchaseID).Error; err != nil {
			return err
		}

		// remove library entry if present
		if err := tx.Where("user_id = ? AND game_id = ?", refund.Purchase.UserID, refund.Purchase.GameID).Delete(&database.LibraryItem{}).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to approve refund: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "refund approved"})
}

// POST /api/secure/add-to-wishlist
func AddToWishlist(c *gin.Context) {
	var requestBody struct {
		GameID uint `json:"game_id"`
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

	if res := database.DB.Model(&database.Game{}).Where("id = ?", requestBody.GameID).Find(nil); res.RowsAffected != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to find game"})
		return
	}

	wishlistItem := database.WishlistItem{UserID: user.ID, GameID: requestBody.GameID}
	if err := database.DB.Create(&wishlistItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to add item to wishlist: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "Added item to wishlist"})
}

// POST /api/secure/remove-from-wishlist
func RemoveFromWishlist(c *gin.Context) {
	var requestBody struct {
		GameID uint `json:"game_id"`
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

	if err := database.DB.Where("user_id = ? AND game_id = ?", user.ID, requestBody.GameID).Delete(&database.WishlistItem{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to remove item from wishlist: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "Removed item from wishlist"})
}

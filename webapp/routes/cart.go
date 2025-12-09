package routes

import (
	"cs3380/database"
	"fmt"
	"net/http"
	"time"

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

		// Check if user already owns this game
		var existingPurchase database.Purchase
		if err := tx.Where("user_id = ? AND game_id = ?", user.ID, requestBody.GameID).First(&existingPurchase).Error; err == nil {
			return fmt.Errorf("you already own this game")
		}

		// Check if already in cart
		var existingCartItem database.CartItem
		if err := tx.Where("user_id = ? AND game_id = ?", user.ID, requestBody.GameID).First(&existingCartItem).Error; err == nil {
			return fmt.Errorf("game already in cart")
		}

		var cartItem database.CartItem
		cartItem.GameID = requestBody.GameID
		cartItem.UserID = user.ID

		return tx.Create(&cartItem).Error
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to add item to cart: %v", err)})
		return
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
			// Check if user already owns this game
			var existingPurchase database.Purchase
			if err := tx.Where("user_id = ? AND game_id = ?", user.ID, item.GameID).First(&existingPurchase).Error; err == nil {
				// User already owns this game, skip it
				continue
			}

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

// POST /api/secure/approve-refund (ADMIN ONLY)
func ApproveRefund(c *gin.Context) {
	var requestBody struct {
		RefundID uint `json:"refund_id" binding:"required"`
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

	if !user.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
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

		// remove library entry if present (hard delete)
		if err := tx.Unscoped().Where("user_id = ? AND game_id = ?", refund.Purchase.UserID, refund.Purchase.GameID).Delete(&database.LibraryItem{}).Error; err != nil {
			return err
		}

		// delete the refund record first (foreign key constraint)
		if err := tx.Unscoped().Delete(&database.Refund{}, requestBody.RefundID).Error; err != nil {
			return err
		}

		// remove purchase (hard delete)
		if err := tx.Unscoped().Delete(&database.Purchase{}, refund.PurchaseID).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to approve refund: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "refund approved"})
}

// GET /api/secure/get-purchases
func GetPurchases(c *gin.Context) {
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

	var purchases []database.Purchase
	if err := database.DB.Preload("Game").Where("user_id = ?", user.ID).Find(&purchases).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get purchases"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"purchases": purchases})
}

// GET /api/secure/get-my-refunds
func GetMyRefunds(c *gin.Context) {
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

	type RefundWithGame struct {
		database.Refund
		GameName string `json:"game_name"`
	}

	var refunds []RefundWithGame
	if err := database.DB.Table("refunds").
		Select("refunds.*, games.name as game_name").
		Joins("JOIN purchases ON purchases.id = refunds.purchase_id").
		Joins("JOIN games ON games.id = purchases.game_id").
		Where("purchases.user_id = ? AND refunds.deleted_at IS NULL", user.ID).
		Scan(&refunds).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get refunds"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"refunds": refunds})
}

// GET /api/secure/get-pending-refunds (ADMIN ONLY)
func GetPendingRefunds(c *gin.Context) {
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

	if !user.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
		return
	}

	type PendingRefund struct {
		database.Refund
		GameName     string    `json:"game_name"`
		Username     string    `json:"username"`
		PurchaseDate time.Time `json:"purchase_date"`
	}

	var pendingRefunds []PendingRefund
	if err := database.DB.Table("refunds").
		Select("refunds.*, games.name as game_name, users.username, purchases.created_at as purchase_date").
		Joins("JOIN purchases ON purchases.id = refunds.purchase_id").
		Joins("JOIN games ON games.id = purchases.game_id").
		Joins("JOIN users ON users.id = purchases.user_id").
		Where("refunds.approved = ? AND refunds.deleted_at IS NULL", false).
		Scan(&pendingRefunds).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get pending refunds"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pending_refunds": pendingRefunds})
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

// GET /api/secure/get-wishlist
func GetWishlist(c *gin.Context) {
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
		var wishlistItems []database.WishlistItem
		if err := tx.Preload("Game").Find(&wishlistItems, "user_id = ?", user.ID).Error; err != nil {
			return err
		}

		for _, item := range wishlistItems {
			result = append(result, item.Game)
		}

		return nil
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to find wishlist items: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"wishlist_items": result})
}

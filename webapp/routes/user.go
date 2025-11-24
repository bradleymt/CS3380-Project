package routes

import (
	"cs3380/database"
	"cs3380/middleware"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func RegisterUser(c *gin.Context) {
	var userRequest database.User
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if hash, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), bcrypt.DefaultCost); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	} else {
		userRequest.Password = string(hash)
	}

	if err := database.DB.Create(&userRequest).Error; err != nil {
		// TODO: remove the error details in production
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user", "error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func LoginUser(c *gin.Context) {
	var loginRequest struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if loginRequest.Email == "" || loginRequest.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email or username is required"})
		return
	}

	var user database.User
	if err := database.DB.Where("email = ? OR username = ?", loginRequest.Email, loginRequest.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username/email or password"})
		return
	}

	if err := user.CheckHash(loginRequest.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username/email or password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
	})

	tokenString, err := token.SignedString(middleware.SecretKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// POST /api/secure/create-family
func CreateFamily(c *gin.Context) {
	var user database.User
	if userID, ok := c.Get("user_id"); !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user ID not found in request context"})
		return
	} else {
		if err := database.DB.First(&user, "id = ?", userID.(uint64)).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find associated user"})
			return
		}
	}

	if user.FamilyID != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "you must leave a family before creating one"})
		return
	}

	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		var family database.Family
		if err := tx.Create(&family).Error; err != nil {
			return fmt.Errorf("failed to create family: %w", err)
		}

		if err := tx.Model(&database.User{}).Where("id = ?", user.ID).Update("family_id", family.ID).Error; err != nil {
			return fmt.Errorf("failed to update user family: %w", err)
		}

		return nil
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to create family: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "created new family"})
}

// POST /api/secure/join-family
func JoinFamily(c *gin.Context) {
	var requestBody struct {
		FamilyID uint `json:"family_id"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to bind request body: %v", err)})
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

	if user.FamilyID != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user already has family"})
		return
	}

	var count int64
	if err := database.DB.Model(&database.Family{}).Where("id = ?", requestBody.FamilyID).Count(&count).Error; err != nil || count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid family id supplied"})
		return
	}

	user.FamilyID = &requestBody.FamilyID
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update familyID for user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "User has been added to family"})
}

// POST /api/secure/leave-family
func LeaveFamily(c *gin.Context) {
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

	if user.FamilyID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot leave a family if you are not a member of any family"})
		return
	}

	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Model(&database.User{}).Where("family_id = ?", *user.FamilyID).Count(&count).Error; err != nil {
			return err
		}

		if count == 1 {
			if err := tx.Model(&database.Family{}).Where("id = ?", *user.FamilyID).Delete(&database.Family{}).Error; err != nil {
				return err
			}
		}
		if err := tx.Model(&database.User{}).Where("id = ?", user.ID).Update("family_id", nil).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to remove user from family: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "User has been removed from family"})
}

// GET /api/secure/get-family
func GetFamily(c *gin.Context) {
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

	if user.FamilyID == nil {
		c.JSON(http.StatusOK, gin.H{"result": "user is not in a family"})
		return
	}

	type FamilyEntry struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	var familyEntries []FamilyEntry
	if err := database.DB.Model(&database.User{}).Where("family_id = ?", *user.FamilyID).Find(&familyEntries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query users for familyID"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"family_members": familyEntries})
}

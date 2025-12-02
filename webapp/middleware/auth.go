package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	// SecretKey used for signing JWT tokens
	SecretKey = []byte("abcdedgfghijklmnopqrstuvwxyz1234567890")
)

func AuthorizeRequest(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" || len(tokenString) < 7 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
		c.Abort()
		return
	}

	if token, err := jwt.Parse(tokenString[7:], func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	}); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token, failed to parse"})
		c.Abort()
		return
	} else if !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	} else {
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
			c.Abort()
			return
		}

		c.Set("user_id", uint64(userID))
	}

	c.Next()
}

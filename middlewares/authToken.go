package middlewares

import (
	"net/http"

	"github.com/Eco-Led/EcoLed-Back_test/services"

	"github.com/gin-gonic/gin"
)

func AuthToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from header
		tokenString := c.GetHeader("Authorization")

		// Check token
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized123",
			})
			return
		}

		// Get userID from token
		userID, err := services.TokenServices{}.ExtractTokenID(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Set userID to context
		c.Set("user_id", userID)

		// Next
		c.Next()
	}
}

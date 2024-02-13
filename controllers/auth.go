package controllers

import (
	"net/http"
	"time"
	"os"

	"github.com/Eco-Led/EcoLed-Back_test/services"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"
)

type AuthControllers struct{}

var tokenService = services.TokenServices{}

func(ctr AuthControllers) RefreshToken(c *gin.Context) {
	// Refresh Token을 요청 헤더나 본문에서 가져옵니다.
	refreshtoken := c.Request.Header.Get("Authorization")
	if refreshtoken == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Refresh token is required"})
		return
	}

	// Parse the token
	claims := &jwt.StandardClaims{}
	_, err := jwt.ParseWithClaims(refreshtoken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})

	// Handle errors from parsing
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Check if token is expired
	if time.Unix(claims.ExpiresAt, 0).Before(time.Now()) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token expired"})
		return
	}

	// Create new Access Token (service)
	userID, err := tokenService.ExtractTokenID(refreshtoken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting user ID from token"})
		return
	}

	tokenDetails, err := tokenService.CreateToken(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating new access token"})
		return
	}

	// Save the details of the new token in Redis (service)
	err = tokenService.SaveToken(userID, tokenDetails)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving the new token details"})
		return
	}

	// Return the new tokens
	c.JSON(http.StatusOK, gin.H{"access_token":  tokenDetails.AccessToken,"refresh_token": tokenDetails.RefreshToken})

}

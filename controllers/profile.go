package controllers

import (
	"net/http"

	"github.com/Eco-Led/EcoLed-Back_test/forms"
	"github.com/Eco-Led/EcoLed-Back_test/services"
	"github.com/Eco-Led/EcoLed-Back_test/utils"

	"github.com/gin-gonic/gin"
)

type ProfileControllers struct{}

var profileService = new(services.ProfileServices)

// Profile is created by register

func (ctr ProfileControllers) UpdateProfile(c *gin.Context) {
	// Bind JSON to profileForm (form)
	var profileForm forms.ProfileForm
	if err := c.ShouldBindJSON(&profileForm); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get userID from token & Chage type to uint (util)
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update profile (service)
	err = profileService.UpdateProfile(userID, profileForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Response
	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})

}


func (ctr ProfileControllers) GetMyProfile(c *gin.Context) {
	// Get userID from token & Chage type to uint (util)
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get profile (service)
	profile, err := profileService.GetProfile(userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Response
	c.JSON(http.StatusOK, gin.H{"profile": profile})

}

func (ctr ProfileControllers) GetUserProfile(c *gin.Context) {
	// Get userID from param (util)
	userID, err := utils.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get profile (service)
	profile, err := profileService.GetProfile(userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Response
	c.JSON(http.StatusOK, gin.H{"profile": profile})

}

func (ctr ProfileControllers) GetProfileByNickname(c *gin.Context) {
	// Get nickname from body (form)
	nickname := c.PostForm("nickname")

	// Get user (service)
	user, err := profileService.GetProfileByNickname(nickname)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Response
	c.JSON(http.StatusOK, gin.H{"user": user})

}

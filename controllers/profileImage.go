package controllers

import (
	"context"
	"net/http"
	"path/filepath"

	"github.com/Eco-Led/EcoLed-Back_test/services"
	"github.com/Eco-Led/EcoLed-Back_test/utils"

	"github.com/gin-gonic/gin"
)

type ProfileImageControllers struct{}

// UploadImage uploads image
func (ctr ProfileImageControllers) UploadProfileImage(c *gin.Context) {
	//By form-data type, file is uploaded
	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	filename := filepath.Base(file.Filename)

	//Open file
	filecontent, _ := file.Open()
	defer filecontent.Close()

	// Get userID from token & Chage type to uint (util)
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Get imageURL
	var imageService services.ImageService
	imageURL, err := imageService.UploadProfileImage(context.Background(), filecontent, userID, filename)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Return imageURL
	c.JSON(http.StatusOK, gin.H{"Success to Upload!": imageURL})

}

func (ctr ProfileImageControllers) DeleteProfileImage(c *gin.Context) {
	// Get userID from token & Chage type to uint (util)
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Delete image(in google cloud storage & DB) (service)
	var imageService services.ImageService
	err = imageService.DeleteProfileImage(context.Background(), userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Return imageURL
	c.JSON(http.StatusOK, gin.H{"message": "Success to Delete!"})

}

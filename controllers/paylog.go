package controllers

import (
	"net/http"

	"github.com/Eco-Led/EcoLed-Back_test/forms"
	"github.com/Eco-Led/EcoLed-Back_test/services"
	"github.com/Eco-Led/EcoLed-Back_test/utils"

	"github.com/gin-gonic/gin"
)

type PaylogControllers struct{}

var paylogService = new(services.PaylogServices)

func (ctr PaylogControllers) CreatePaylog(c *gin.Context) {
	// Bind paylogForm from JSON (form)
	var paylogForm forms.PaylogForm
	if err := c.ShouldBindJSON(&paylogForm); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get userID from token & Chage type to uint (util)
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create paylog (service)
	err = paylogService.CreatePaylog(userID, paylogForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Create Success"})

}


func (ctr PaylogControllers) UpdatePaylog(c *gin.Context) {
	// Get paylogID from param (util)
	paylogID, err := utils.GetPaylogID(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}	

	// Get userID from token & Chage type to uint (util)
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get paylogForm from JSON (form)
	var paylogForm forms.PaylogForm
	if err = c.ShouldBindJSON(&paylogForm); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update paylog (service)
	err = paylogService.UpdatePaylog(userID, paylogID, paylogForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Response paylog
	c.JSON(http.StatusOK, gin.H{"message": "Update Success"})

}


func (ctr PaylogControllers) DeletePaylog(c *gin.Context) {
	// Get paylogID from param (util)
	paylogID, err := utils.GetPaylogID(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}	

	// Get userID from token & Chage type to uint (util)
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Delete paylog (service)
	err = paylogService.DeletePaylog(userID, paylogID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Response paylog
	c.JSON(http.StatusOK, gin.H{"message": "Delete Success"})

}

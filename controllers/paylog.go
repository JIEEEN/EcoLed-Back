package controllers

import (
	"net/http"
	"strconv"

	"github.com/Eco-Led/EcoLed-Back_test/forms"
	"github.com/Eco-Led/EcoLed-Back_test/services"
	"github.com/Eco-Led/EcoLed-Back_test/utils"

	"github.com/gin-gonic/gin"
)

type PaylogControllers struct{}

var paylogService = new(services.PaylogServices)

func (ctr PaylogControllers) CreatePaylog(c *gin.Context) {
	// Bind paylogForm from JSON
	var paylogForm forms.PaylogForm
	if err := c.ShouldBindJSON(&paylogForm); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Get userID from token & Chage type to uint
	userID, err := utils.GetUserIDFromContext(c)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	// Create paylog (service)
	err = paylogService.CreatePaylog(userID, paylogForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Create Success",
	})

}

func (ctr PaylogControllers) UpdatePaylog(c *gin.Context) {
	// Get paylogID from param
	paylogIDstring := c.Param("paylogID")
	paylogIDint64, err := strconv.ParseUint(paylogIDstring, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "failed to get paylogID",
		})
		return
	}
	paylogID := uint(paylogIDint64)

	// Get userID from token & Chage type to uint
	userID, err := utils.GetUserIDFromContext(c)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	// Get paylogForm from JSON
	var paylogForm forms.PaylogForm
	if err = c.ShouldBindJSON(&paylogForm); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Update paylog (service)
	err = paylogService.UpdatePaylog(userID, paylogID, paylogForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Response paylog
	c.JSON(http.StatusOK, gin.H{
		"message": "Update Success",
	})

}

func (ctr PaylogControllers) DeletePaylog(c *gin.Context) {
	// Get paylogID from param
	paylogIDstring := c.Param("paylogID")
	paylogIDint64, err := strconv.ParseUint(paylogIDstring, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "failed to get paylogID",
		})
		return
	}
	paylogID := uint(paylogIDint64)

	// Get userID from token & Chage type to uint
	userID, err := utils.GetUserIDFromContext(c)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	// Delete paylog (service)
	err = paylogService.DeletePaylog(userID, paylogID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Response paylog
	c.JSON(http.StatusOK, gin.H{
		"message": "Delete Success",
	})

}

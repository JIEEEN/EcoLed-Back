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

//TODO: 프론트로부터 시간 정보 어떻게 받아오는지 확인 필요 공부해야함 지금은 2024-01-26, 13:30 형식으로 받아옴
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

func (ctr PaylogControllers) GetPaylog(c *gin.Context) {
	// Get userID from token & Chage type to uint (util)
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Get paylogID from param (util)
	paylogID, err := utils.GetPaylogID(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get paylog (service)
	paylog, err := paylogService.GetPaylog(userID, paylogID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Return paylog
	c.JSON(http.StatusOK, gin.H{"paylog": paylog})

}

func (ctr PaylogControllers) GetPaylogs(c *gin.Context) {
	// Get userID from token & Chage type to uint (util)
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get account (service)
	account, err := accountService.GetAccount(userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get page from query (util)
	page, err := utils.GetPage(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get paylogs (service)
	paylogs, err := paylogService.GetPaylogs(account.ID, page)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Return paylogs
	c.JSON(http.StatusOK, gin.H{"paylogs": paylogs})

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

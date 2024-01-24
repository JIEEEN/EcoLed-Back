package controllers

import (
	"net/http"
	"strconv"

	"github.com/Eco-Led/EcoLed-Back_test/services"
	"github.com/Eco-Led/EcoLed-Back_test/utils"

	"github.com/gin-gonic/gin"
)

type AccountControllers struct{}

var accountService = new(services.AccountServices)

func (ctr AccountControllers) GetAccount(c *gin.Context) {
	// Get userID from token & Chage type to uint
	userID, err := utils.GetUserIDFromContext(c)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	// Get account (service)
	account, err := accountService.GetAccount(userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Return account
	c.JSON(http.StatusOK, gin.H{
		"account": account,
	})

	// Get page from query
	pageString := c.Query("page")
	pageInt, err := strconv.ParseInt(pageString, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "failed to get page",
		})
		return
	}
	page := int(pageInt)

	// Get paylogs (service)
	paylogs, err := paylogService.GetPaylogs(account.ID, page)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Return paylogs
	c.JSON(http.StatusOK, gin.H{
		"paylogs": paylogs,
	})

}

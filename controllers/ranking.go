package controllers

import (
	"net/http"

	"github.com/Eco-Led/EcoLed-Back_test/services"

	"github.com/gin-gonic/gin"
)

type RankingControllers struct{}

var rankingService = new(services.RankingServices)

func (ctr RankingControllers) GetRanking(c *gin.Context) {
	// Get ranking (service)
	ranking, err := rankingService.GetRanking()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Response
	c.JSON(http.StatusOK, gin.H{"ranking": ranking})

}

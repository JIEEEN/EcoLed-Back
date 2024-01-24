package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
)

// GetQuery gets query from context
func GetPage(c *gin.Context) (int, error) {
	pageString := c.Query("page")
	pageInt, err := strconv.ParseInt(pageString, 10, 64)
	if err != nil {
		return 0, errors.New("failed to get page")
	}
	return int(pageInt), nil
}
package utils

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetPaylogID gets paylogID from context
func GetPaylogID(c *gin.Context) (uint, error) {
	paylogIDstring := c.Param("paylogID")
	paylogIDint64, err := strconv.ParseUint(paylogIDstring, 10, 64)
	if err != nil {
		return 0, errors.New("failed to get paylogID")
	}
	return uint(paylogIDint64), nil
}

// GetPostID gets postID from context
func GetPostID(c *gin.Context) (uint, error) {
	postIDstring := c.Param("postID")
	postIDint64, err := strconv.ParseUint(postIDstring, 10, 64)
	if err != nil {
		return 0, errors.New("failed to get postID")
	}
	return uint(postIDint64), nil
}
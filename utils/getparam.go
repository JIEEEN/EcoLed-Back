package utils

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetPaylogID gets paylogID from param	
func GetPaylogID(c *gin.Context) (uint, error) {
	paylogIDstring := c.Param("paylogID")
	paylogIDint64, err := strconv.ParseUint(paylogIDstring, 10, 64)
	if err != nil {
		return 0, errors.New("failed to get paylogID")
	}
	return uint(paylogIDint64), nil
}

// GetPostID gets postID from param
func GetPostID(c *gin.Context) (uint, error) {
	postIDstring := c.Param("postID")
	postIDint64, err := strconv.ParseUint(postIDstring, 10, 64)
	if err != nil {
		return 0, errors.New("failed to get postID")
	}
	return uint(postIDint64), nil
}

// GetUserID gets userID from param
func GetUserID(c *gin.Context) (uint, error) {
	userIDstring := c.Param("userID")
	userIDint64, err := strconv.ParseUint(userIDstring, 10, 64)
	if err != nil {
		return 0, errors.New("failed to get userID")
	}
	return uint(userIDint64), nil
}

// GetNickname gets nickname from param
func GetNickname(c *gin.Context) (string, error) {
	nickname := c.Param("nickname")
	if nickname == "" {
		return "", errors.New("failed to get nickname")
	}
	return nickname, nil
}
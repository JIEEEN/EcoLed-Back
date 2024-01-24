package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
)


func GetUserIDFromContext(c *gin.Context) (uint, error) {
	userIDInterface, ok := c.Get("user_id")
	if !ok {
		return 0, errors.New("failed to get userIDInterface")
	}
	userIDInt64, ok := userIDInterface.(int64)
	if !ok {
		return 0, errors.New("failed to get userID into int64")
	}
	return uint(userIDInt64), nil
}
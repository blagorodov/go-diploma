package auth

import (
	"diploma/internal/logger"
	"fmt"
	"github.com/gin-gonic/gin"
)

func WithToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := GetID(c)
		logger.Log("auth::WithToken")
		logger.Log(fmt.Sprintf("userID: %s", userID))
		c.Set("UserID", userID)
		SetID(c, userID)
		c.Next()
	}
}

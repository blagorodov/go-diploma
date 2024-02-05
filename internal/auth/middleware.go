package auth

import (
	"github.com/gin-gonic/gin"
)

func WithToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := GetID(c)
		c.Set("UserID", userID)
		SetID(c, userID)
		c.Next()
	}
}

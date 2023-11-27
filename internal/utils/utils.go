package utils

import "github.com/gin-gonic/gin"

func ErrorJSON(c *gin.Context, code int, data any) {
	c.JSON(code, gin.H{"error": data})
}

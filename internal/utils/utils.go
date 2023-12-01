package utils

import (
	"diploma/internal/logger"
	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/gin-gonic/gin"
)

func ErrorJSON(c *gin.Context, code int, data any) {
	if data == nil {
		c.Status(code)
	} else {
		c.JSON(code, gin.H{"error": data})
	}
}

func ValidateLuhn(number string) bool {
	err := goluhn.Validate(number)
	return err == nil
}

func GetBody(ctx *gin.Context) string {
	data, err := ctx.GetRawData()
	if err != nil {
		logger.Log("GetBody")
		logger.Log(err.Error())
		return ""
	}
	return string(data)
}

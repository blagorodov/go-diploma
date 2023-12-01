package drivers

import (
	"diploma/internal/logger"
	"github.com/aurowora/compress"
	"github.com/gin-gonic/gin"
)

type GinRouter struct {
	Gin *gin.Engine
}

func NewGinRouter() GinRouter {
	httpRouter := gin.Default()
	httpRouter.Use(compress.Compress(), logger.WithLogging())
	return GinRouter{
		Gin: httpRouter,
	}
}

package drivers

import (
	"diploma/internal/auth"
	"diploma/internal/logger"
	"github.com/aurowora/compress"
	"github.com/gin-gonic/gin"
)

type GinRouter struct {
	Gin *gin.Engine
}

func NewGinRouter() GinRouter {
	httpRouter := gin.Default()
	httpRouter.Use(compress.Compress(), logger.WithLogging(), auth.WithToken())
	return GinRouter{
		Gin: httpRouter,
	}
}

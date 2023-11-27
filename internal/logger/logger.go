package logger

import (
	"diploma/internal/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		gin.ResponseWriter
		responseData *responseData
	}
)

var sugar zap.SugaredLogger

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func Init() {
	logger, err := NewLogger()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	sugar = *logger.Sugar()
}

func NewLogger() (*zap.Logger, error) {
	encCfg := zap.NewProductionEncoderConfig()
	encCfg.EncodeTime = zapcore.RFC3339TimeEncoder

	cfg := zap.Config{
		Level:         zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:   true,
		Encoding:      "console",
		EncoderConfig: encCfg,
		OutputPaths: []string{
			config.Options.LogPath,
			"stdout",
		},
	}

	return cfg.Build()
}

func WithLogging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		c.Writer = &loggingResponseWriter{
			ResponseWriter: c.Writer,
			responseData:   responseData,
		}

		c.Next()

		duration := time.Since(start)
		sugar.Infoln(
			"uri", c.Request.RequestURI,
			"method", c.Request.Method,
			"status", responseData.status,
			"duration", duration,
			"size", responseData.size,
		)
	}
}

func Log(s any) {
	sugar.Infoln(s)
}

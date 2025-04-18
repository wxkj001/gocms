package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Logger struct {
	config *viper.Viper
	log    *zap.Logger
}

func NewLogger(config *viper.Viper, log *zap.Logger) *Logger {
	return &Logger{config: config, log: log}
}
func (l *Logger) Logger(c *gin.Context) {
	// Start timer
	start := time.Now()
	path := c.Request.URL.Path
	raw := c.Request.URL.RawQuery

	// Process request
	c.Next()

	// Log only when it is not being skipped
	l.log.Info("gin",
		zap.Any("TimeStamp", start.Format("2006-01-02 15:04:05")),
		zap.Any("Method", c.Request.Method),
		zap.Any("Path", path+"?"+raw),
		zap.Any("StatusCode", c.Writer.Status()),
		zap.Any("Latency", time.Now().Sub(start)),
		zap.Any("ClientIP", c.ClientIP()))
}

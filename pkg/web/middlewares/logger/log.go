package logger

import (
	"github.com/gopusher/gateway/pkg/log"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		c.Next()
		param := gin.LogFormatterParams{
			Request: c.Request,
			Keys:    c.Keys,
		}
		param.TimeStamp = time.Now()
		param.Latency = param.TimeStamp.Sub(start)
		param.ClientIP = c.ClientIP()
		param.Method = c.Request.Method
		param.StatusCode = c.Writer.Status()
		//param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
		//param.BodySize = c.Writer.Size()
		if raw != "" {
			path = path + "?" + raw
		}
		param.Path = path
		if param.Latency > time.Minute {
			param.Latency = param.Latency - param.Latency%time.Second
		}
		log.Info("[GIN]",
			zap.String("tag", "request"),
			//zap.String("time", param.TimeStamp.Format("2006/01/02 - 15:04:05")),
			zap.Int("statusCode", param.StatusCode),
			zap.Duration("latency", param.Latency),
			zap.String("ip", param.ClientIP),
			zap.String("method", param.Method),
			zap.String("path", param.Path),
			//zap.String("errorMessage", param.ErrorMessage),
		)
	}
}

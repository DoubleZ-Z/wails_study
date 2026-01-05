package interceptor

import (
	"time"
	"wails_study/project/logger"

	"github.com/gin-gonic/gin"
)

func Log(c *gin.Context) {
	start := time.Now()
	path := c.Request.URL.Path
	query := c.Request.URL.RawQuery
	c.Next()
	cost := time.Since(start)
	status := c.Writer.Status()

	fields := []interface{}{
		"method", c.Request.Method,
		"path", path,
		"query", query,
		"status", status,
		"cost_ms", cost.Milliseconds(),
		"client_ip", c.ClientIP(),
	}
	if status >= 500 {
		logger.Errors("HTTP 5xx 错误", fields...)
	} else if status >= 400 {
		logger.Warns("HTTP 4xx 客户端错误", fields...)
	} else {
		logger.Infos("HTTP 请求", fields...)
	}
}

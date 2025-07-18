package middleware

import (
	"bytes"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/algorithm9/flash-deal/pkg/logger"
)

// TraceAndLogger 注入 traceID 并记录请求日志
func TraceAndLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := uuid.New().String()
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, "trace_id", traceID)
		c.Request = c.Request.WithContext(ctx)

		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		status := c.Writer.Status()
		cost := time.Since(start)

		logger.L().
			Info().
			Str("trace_id", traceID).
			Msgf("HTTP Request method %s path:%s status:%d cost:%s", method, path, status, cost)
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// RecoveryWithLog 替代 gin.Recovery，统一捕获 panic 并记录日志
func RecoveryWithLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.L().Panic().
					Str("trace_id", c.GetString("trace_id")).
					Msgf("panic recover:%v", err)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

package middleware

import (
	"runtime/debug"
	"time"

	"github.com/img21326/andrew_framework/helper"

	"github.com/gin-gonic/gin"
	uuid "github.com/google/uuid"
)

func WithLoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()
		uuid := uuid.New().String()
		logger := helper.NewLogger(uuid)

		defer func() {
			if err := recover(); err != nil {
				logger.Error(ctx, "error: %v", string(debug.Stack()))
			}
		}()

		ctx.Set("logger", logger)
		ctx.Set("uuid", uuid)
		body, err := ctx.GetRawData()
		if err != nil {
			logger.Error(ctx, "get raw data error: %v", err)
		}
		logger.Info(ctx, "request params: %+v, request data: %+v", ctx.Request.URL.Query(), string(body))

		ctx.Next()
		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)
		reqMethod := ctx.Request.Method
		reqUri := ctx.Request.RequestURI
		statusCode := ctx.Writer.Status()
		clientIP := ctx.ClientIP()
		logger.Info(ctx, "| %3d | %13v | %15s | %s | %s |",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
		)
	}
}

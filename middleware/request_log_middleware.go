package middleware

import (
	"time"

	"github.com/img21326/andrew_framework/helper"

	"github.com/gin-gonic/gin"
)

func WithRequestLogMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logger := ctx.MustGet("logger").(*helper.Logger)
		startTime := time.Now()

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

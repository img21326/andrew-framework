package middleware

import (
	"bytes"
	"io"
	"strings"
	"time"

	"github.com/img21326/andrew_framework/helper"

	"github.com/gin-gonic/gin"
)

func WithRequestLogMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logger := ctx.MustGet("logger").(*helper.Logger)
		startTime := time.Now()

		contentType := ctx.GetHeader("Content-Type")
		// 確認不是上傳檔案的資料
		if strings.Contains(contentType, "multipart/form-data") {
			body, err := io.ReadAll(ctx.Request.Body)
			if err != nil {
				logger.Error(ctx, "get raw data error: %v", err)
			}
			logger.Info(ctx, "request params: %+v, request data: %+v", ctx.Request.URL.Query(), string(body))
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}

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

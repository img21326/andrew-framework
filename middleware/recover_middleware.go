package middleware

import (
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/img21326/andrew_framework/helper"
)

func WithRecoverMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger := ctx.MustGet("logger").(*helper.Logger)
				logger.Error(ctx, "error: %v", err)
				logger.Error(ctx, "stack: %s", debug.Stack())

				ctx.JSON(500, gin.H{
					"code":    500,
					"message": "Internal Server Error",
				})
				ctx.Abort()
			}
		}()

		ctx.Next()
	}
}

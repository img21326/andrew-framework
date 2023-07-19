package middleware

import (
	"github.com/img21326/andrew_framework/helper"

	"github.com/gin-gonic/gin"
	uuid "github.com/google/uuid"
)

func WithLoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uuid := uuid.New().String()
		logger := helper.NewLogger(uuid)

		ctx.Set("logger", logger)
		ctx.Set("uuid", uuid)

		ctx.Next()
	}
}

package middleware

import (
	"github.com/img21326/andrew-framework/helper"

	"github.com/gin-gonic/gin"
	uuid "github.com/google/uuid"
)

func WithLoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var uid string
		uid, ok := ctx.Value("aws_request_id").(string)
		if !ok {
			uid = uuid.New().String()
		}
		logger := helper.NewLogger(uid)

		ctx.Set("logger", logger)
		ctx.Set("uuid", uid)

		ctx.Next()

		defer logger.Close()
	}
}

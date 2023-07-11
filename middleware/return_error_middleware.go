package middleware

import (
	"errors"

	"github.com/img21326/andrew_framework/helper"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ReturnErrorMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		logger := ctx.MustGet("logger").(helper.Logger)
		if len(ctx.Errors) > 0 {
			for _, err := range ctx.Errors {
				logger.Error(ctx, "error: %v", err)
			}
			err := ctx.Errors.Last().Err
			errorKey := "internal_server_error"
			if errors.Is(err, gorm.ErrRecordNotFound) {
				errorKey = "not_found"
			}
			ctx.JSON(helper.ErrorMap[errorKey].Code, helper.ErrorMap[errorKey])
			ctx.Abort()
		}
	}
}

package middleware

import (
	"github.com/img21326/andrew_framework/helper"

	"github.com/gin-gonic/gin"
)

func ReturnErrorMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		logger := ctx.MustGet("logger").(*helper.Logger)
		if len(ctx.Errors) > 0 {
			for _, err := range ctx.Errors {
				logger.Error(ctx, "error: %v", err)
			}
			err, ok := ctx.Errors.Last().Err.(helper.ErrorInterface)
			if ok && err != nil {
				err = helper.ErrorMap[err.ErrorCode()]
				ctx.JSON(err.HttpStatus(), gin.H{
					"code":    err.ErrorCode(),
					"message": err.Message(),
				})
				ctx.Abort()
			} else {
				panic("undefined error")
			}
		}
	}
}

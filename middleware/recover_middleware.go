package middleware

import (
	"fmt"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/img21326/andrew-framework/helper"
	"github.com/spf13/viper"
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

				if mailHelper := helper.GetEmailHelper(); mailHelper != nil && gin.Mode() == gin.ReleaseMode {
					viper := viper.GetViper()
					adminEmail := viper.GetStringSlice("ADMIN_EMAIL")
					if len(adminEmail) == 0 {
						return
					}
					requestRawBody := ctx.Value("raw_body").(*[]byte)
					body := fmt.Sprintf("request: %s %s\ndata: %+v\nerror: %v\nstack: %s", ctx.Request.Method, ctx.Request.URL.String(), string(*requestRawBody), err, debug.Stack())
					mailHelper.SendEmail(helper.EmailSendOption{
						To:      adminEmail,
						Subject: "Internal Server Error",
						Body:    body,
					})
				}
			}
		}()

		ctx.Next()
	}
}

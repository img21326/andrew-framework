package middleware

import (
	"github.com/img21326/andrew_framework/helper"

	"github.com/gin-gonic/gin"
)

func WithGormMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_ = ctx.MustGet("logger").(*helper.Logger)
		gorm := helper.NewGorm(helper.GormOption{
			Ctx: ctx,
		})
		ctx.Set("gormDB", gorm)
		ctx.Next()
	}
}

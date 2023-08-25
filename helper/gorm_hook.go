package helper

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var gormHook []func(*gorm.DB, *gin.Context)

func RegisterGormHook(f func(*gorm.DB, *gin.Context)) {
	gormHook = append(gormHook, f)
}

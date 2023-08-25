package helper

import "gorm.io/gorm"

var gormHook []func(*gorm.DB)

func RegisterGormHook(f func(*gorm.DB)) {
	gormHook = append(gormHook, f)
}

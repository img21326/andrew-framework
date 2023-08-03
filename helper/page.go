package helper

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func Pagination(ctx *gin.Context, perPageCount int) (limit int, offset int) {
	pageString := ctx.Query("page")
	page, err := strconv.Atoi(pageString)
	if err != nil {
		page = 1
	}
	limit = perPageCount
	offset = (page - 1) * perPageCount
	return
}

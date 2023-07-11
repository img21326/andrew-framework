package andrewframework

import "github.com/gin-gonic/gin"

type RouterInterface interface {
	AddRoute(r *gin.Engine)
}

var RouterList []RouterInterface

func AddRouter(r RouterInterface) {
	RouterList = append(RouterList, r)
}

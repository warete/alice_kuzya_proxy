package serviceproxy

import "github.com/gin-gonic/gin"

type IServiceProxy interface {
	AddRoutes(r *gin.Engine)
}

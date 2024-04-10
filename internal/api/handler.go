package api

import "github.com/gin-gonic/gin"

type IHandler interface {
	Register(rg *gin.RouterGroup)
}

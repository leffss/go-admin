package router

import (
	"github.com/gin-gonic/gin"
	"github.com/leffss/go-admin/app/admin/api"
)

func InitRouter(router *gin.RouterGroup) *gin.RouterGroup {
	router.GET("/test", api.Test)
	router.GET("/test/:id", api.TestTask)
	return router
}

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/leffss/go-admin/app/auth/api"
)

func InitRouter(router *gin.RouterGroup) *gin.RouterGroup {
	router.POST("/login", api.Login)
	router.POST("/logout", api.LogOut)
	router.GET("/info", api.Info)
	return router
}

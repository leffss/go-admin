package router

import (
	"github.com/gin-gonic/gin"
	"github.com/leffss/go-admin/app/profile/api"
)

func InitRouter(router *gin.RouterGroup) *gin.RouterGroup {
	router.PUT("/information", api.ChangeInformation)
	router.PUT("/password", api.ChangePassword)
	return router
}

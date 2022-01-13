package api

import (
	"github.com/gin-gonic/gin"
	"github.com/leffss/go-admin/app/auth/service"
	"github.com/leffss/go-admin/pkg/app"
	"github.com/leffss/go-admin/pkg/codemsg"
	"github.com/leffss/go-admin/pkg/jwt"
	"net/http"
)

// Info
// @Summary 当前用户信息
// @Description 当前用户信息，包含权限信息
// @Tags 认证
// @Produce json
// @Param token header string true "token"
// @Success 200 {object} app.Response
// @Failure 400 {object} app.Response
// @Router /api/oauth/info [get]
func Info(c *gin.Context) {
	appG := app.Gin{C: c}
	token := c.GetString("token")
	claims, _ := jwt.ParseToken(token)
	user := service.GetUserByUsername(claims.UserName)
	user.GetPermissions()
	user.Avatar = "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif"
	appG.Response(http.StatusOK, codemsg.Success, user)
}

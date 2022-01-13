// jwt 的设计原则是无状态的，只保存在客户端，客户端发过来后服务器解密得到相关信息。
// 即获取 token 后，不能撤销 token （删除、变更），只能等其自动失效。所以也就不需
// 要传统的 /logout 了，客户端只需等 token 过期或者删除有效的 token 即可。
// 如果想做到失效 token，需要将 token 存入 DB（如 Redis）中，失效则删除；但增加了一个每次校验
// 时候都要先从 DB 中查询 token 是否存在的步骤，而且违背了 JWT 的无状态原则（这不就和 session
// 一样了么？）。

// 如果要所有生成的 token 直接失效，还有一种方案就是动态修改 app.ini 中 JwtSecret 的值，修改后
// 前面所有以此加密的 token 都会失效，但是动态修改需要用到读写锁，会降低一些性能，需要权衡利弊。

package api

import (
	"github.com/gin-gonic/gin"
	"github.com/leffss/go-admin/tasks"
	"net/http"

	"github.com/leffss/go-admin/pkg/app"
	"github.com/leffss/go-admin/pkg/codemsg"
	"github.com/leffss/go-admin/pkg/jwt"
)

// LogOut
// @Summary 退出
// @Description 由于 JWT 机制的原因，退出操作已无意义，这个接口只是配合前端而已
// @Tags 认证
// @Produce json
// @Param token header string true "token"
// @Success 200 {object} app.Response
// @Failure 400 {object} app.Response
// @Router /api/oauth/logout [post]
func LogOut(c *gin.Context) {
	appG := app.Gin{C: c}

	//token, exist := c.Get("token")
	//
	//if exist {
	//	claims, err := jwt.ParseToken(common.StrVal(token))
	//	if err == nil {
	//		_, err = tasks.SendTask("authTaskInsertLoginLog", claims.UserName, c.ClientIP(), c.Request.UserAgent(), claims.UserName + " - 注销成功", true, "注销")
	//		if err != nil {
	//			logger.Error(err.Error())
	//		}
	//	}
	//}

	// token 在中间件中已经验证过正确性了
	token := c.GetString("token")
	claims, _ := jwt.ParseToken(token)
	_, err := tasks.SendTask("authTaskInsertLoginLog", claims.UserName, c.ClientIP(), c.Request.UserAgent(), claims.UserName + " - 注销成功", true, "注销")
	if err != nil {
		logger.Error(err.Error())
	}
	appG.Response(http.StatusOK, codemsg.Success, nil)
}

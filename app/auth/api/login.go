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
	"github.com/leffss/go-admin/app/auth/service"
	"github.com/leffss/go-admin/pkg/common"
	"github.com/leffss/go-admin/tasks"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/leffss/go-admin/pkg/app"
	"github.com/leffss/go-admin/pkg/codemsg"
	"github.com/leffss/go-admin/pkg/jwt"
)

type LoginFrom struct {
	Username string `json:"username" valid:"Required;MinSize(3);MaxSize(32)"`
	Password string `json:"password" valid:"Required;MinSize(6);MaxSize(32)"`
}

// Login
// @Summary 登陆
// @Description 登陆获取 TOKEN
// @Tags 认证
// @Accept json
// @Produce json
// @Param username body string true "用户名"
// @Param password body string true "密码"
// @Success 200 {object} app.Response
// @Failure 400 {object} app.Response
// @Router /api/oauth/login [post]
func Login(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form LoginFrom
	)

	_, errCode, err := app.BindJsonAndValid(c, &form)
	if err != nil {
		logger.Error(err.Error())
		appG.Response(http.StatusBadRequest, codemsg.Error, nil)
		return
	}

	if errCode != codemsg.Success {
		appG.Response(http.StatusBadRequest, errCode, nil)
		return
	}

	user := service.GetUserByUsername(form.Username)
	if user.Username == "" {
		logger.Warn(form.Username + " - " + codemsg.GetMsg(codemsg.ErrorUserNotExist))
		appG.Response(http.StatusBadRequest, codemsg.ErrorUserNotExist, nil)
		return
	}

	if common.EncodeSHA256(form.Password) != user.Password {
		logger.Warn(user.Username + " - " + codemsg.GetMsg(codemsg.ErrorUserWrongPass))
		_, err := tasks.SendTask("authTaskInsertLoginLog", user.Username, c.ClientIP(), c.Request.UserAgent(), user.Username + " - " + codemsg.GetMsg(codemsg.ErrorUserWrongPass), false, "登陆")
		if err != nil {
			logger.Error(err.Error())
		}
		appG.Response(http.StatusBadRequest, codemsg.ErrorUserWrongPass, nil)
		return
	}

	if !user.Status {
		logger.Warn(user.Username + " - " + codemsg.GetMsg(codemsg.ErrorUserStatus))
		_, err := tasks.SendTask("authTaskInsertLoginLog", user.Username, c.ClientIP(), c.Request.UserAgent(), user.Username + " - " + codemsg.GetMsg(codemsg.ErrorUserStatus), false, "登陆")
		if err != nil {
			logger.Error(err.Error())
		}
		appG.Response(http.StatusBadRequest, codemsg.ErrorUserStatus, nil)
		return
	}

	token, expireTime, err := jwt.GenerateToken(user.Id, user.Username)
	if err != nil {
		logger.Error(err.Error())
		_, err := tasks.SendTask("authTaskInsertLoginLog", user.Username, c.ClientIP(), c.Request.UserAgent(), user.Username + " - " + codemsg.GetMsg(codemsg.ErrorTokenCreate), false, "登陆")
		if err != nil {
			logger.Error(err.Error())
		}
		appG.Response(http.StatusBadRequest, codemsg.ErrorTokenCreate, nil)
		return
	}

	_, err = tasks.SendTask("authTaskInsertLoginLog", user.Username, c.ClientIP(), c.Request.UserAgent(), user.Username + " - 登陆成功", true, "登陆")
	if err != nil {
		logger.Error(err.Error())
	}

	appG.Response(http.StatusOK, codemsg.Success, gin.H{
		"token": token,
		"expire_time": expireTime.Format(time.RFC3339),
	})
}

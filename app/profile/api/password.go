package api

import (
	"github.com/gin-gonic/gin"
	"github.com/leffss/go-admin/app/auth/service"
	profileService "github.com/leffss/go-admin/app/profile/service"
	"github.com/leffss/go-admin/pkg/app"
	"github.com/leffss/go-admin/pkg/codemsg"
	"github.com/leffss/go-admin/pkg/common"
	"github.com/leffss/go-admin/pkg/jwt"
	"net/http"
)

type PasswordFrom struct {
	OldPassword string `json:"old_password" validate:"required,min=6,max=32"`
	Password string `json:"password" validate:"required,min=6,max=32,nefield=OldPassword"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=6,max=32,eqfield=Password"`
}

// ChangePassword
// @Summary 修改个人密码
// @Description 修改个人密码
// @Tags 个人中心
// @Produce json
// @Param old_password body string true "旧密码"
// @Param password body string true "新密码"
// @Param confirm_password body string true "确认密码"
// @Success 200 {object} app.Response
// @Failure 400 {object} app.Response
// @Router /api/profile/password [put]
func ChangePassword(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form PasswordFrom
	)
	_, errCode, err := app.BindJsonAndValid2(c, &form)
	if err != nil {
		logger.Error(err.Error())
		appG.Response(http.StatusBadRequest, codemsg.ErrorInvalidParams, nil)
		return
	}

	if errCode != codemsg.Success {
		appG.Response(http.StatusBadRequest, errCode, nil)
		return
	}

	// token 在中间件中已经验证过正确性了
	token := c.GetString("token")
	claims, _ := jwt.ParseToken(token)

	user := service.GetUserByUsername(claims.UserName)
	if user.Username == "" {
		logger.Warn(claims.UserName + " - " + codemsg.GetMsg(codemsg.ErrorUserNotExist))
		appG.Response(http.StatusBadRequest, codemsg.ErrorUserNotExist, nil)
		return
	}

	if common.EncodeSHA256(form.OldPassword) != user.Password {
		appG.Response(http.StatusBadRequest, codemsg.ErrorUserWrongPass, nil)
		return
	}

	updateValues := map[string]interface{}{
		"password": common.EncodeSHA256(form.Password),
	}
	profileService.UpdatesUser(user, updateValues)

	appG.Response(http.StatusOK, codemsg.Success, nil)
}

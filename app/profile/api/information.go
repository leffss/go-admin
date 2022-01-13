package api

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	authService "github.com/leffss/go-admin/app/auth/service"
	profileService "github.com/leffss/go-admin/app/profile/service"
	"github.com/leffss/go-admin/pkg/app"
	"github.com/leffss/go-admin/pkg/codemsg"
	"github.com/leffss/go-admin/pkg/jwt"
	"net/http"
	"regexp"
)

type InformationFrom struct {
	Name string `json:"name" validate:"omitempty,min=2,max=32"`
	//Mobile string `json:"mobile" validate:"omitempty,e164"`
	Mobile string `json:"mobile" validate:"omitempty,MobileValidation"`
	Email string `json:"email" validate:"omitempty,email"`
}

// MobileValidation 自定义验证函数
func MobileValidation(fl validator.FieldLevel) bool {
	rule := "^[1][3-9][0-9]{9}$"
	ok, err := regexp.MatchString(rule, fl.Field().String())
	if err != nil {
		logger.Error(err.Error())
		return false
	}
	return ok
}

// NameValidation 自定义验证函数，只是示例，未使用
func NameValidation(fl validator.FieldLevel) bool {
	return fl.Field().String() == "admin"
}

// ChangeInformation
// @Summary 修改个人信息
// @Description 修改个人信息
// @Tags 个人中心
// @Produce json
// @Param name body string false "姓名"
// @Param mobile body int false "手机"
// @Param email body string false "邮箱"
// @Success 200 {object} app.Response
// @Failure 400 {object} app.Response
// @Router /api/profile/information [put]
func ChangeInformation(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form InformationFrom
	)
	httpCode, errCode, err := app.BindJsonAndValid2(c, &form, MobileValidation, NameValidation)
	if err != nil {
		logger.Error(err.Error())
		appG.Response(httpCode, codemsg.ErrorInvalidParams, nil)
		return
	}

	if errCode != codemsg.Success {
		appG.Response(httpCode, errCode, nil)
		return
	}

	// token 在中间件中已经验证过正确性了
	token := c.GetString("token")
	claims, _ := jwt.ParseToken(token)

	user := authService.GetUserByUsername(claims.UserName)
	if user.Username == "" {
		logger.Warn(claims.UserName + " - " + codemsg.GetMsg(codemsg.ErrorUserNotExist))
		appG.Response(http.StatusBadRequest, codemsg.ErrorUserNotExist, nil)
		return
	}

	updateValues := map[string]interface{}{
		"name": form.Name,
		"mobile": form.Mobile,
		"email": form.Email,
	}
	profileService.UpdatesUser(user, updateValues)

	appG.Response(http.StatusOK, codemsg.Success, nil)
}

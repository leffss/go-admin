package task

import (
	"fmt"
	"github.com/gojuukaze/YTask/v2/message"
	"github.com/json-iterator/go"
	"github.com/leffss/go-admin/app/auth/service"
	"github.com/leffss/go-admin/models"
	"github.com/leffss/go-admin/pkg/logging"
	"github.com/mssola/user_agent"
)

var logger = logging.GetLogger()

func InsertLoginLog(username, ip, userAgent, detail string, status bool, kind string) string {
	user := service.GetUserByUsername(username)
	if user.Username == "" {
		msg := fmt.Sprintf("保存登陆日志错误：用户 [%s] 不存在", username)
		logger.Error(msg)
		return msg
	}
	UA := user_agent.New(userAgent)
	clientName, clientVersion := UA.Browser()
	loginLog := &models.LoginLog{
		Ip:           ip,
		Client:       clientName + " " + clientVersion,
		Os:           UA.OS(),
		Status:       status,
		Detail:       detail,
		Kind:         kind,
		DeptId:       user.DeptId,
		ControlBy:    models.ControlBy{CreateBy: user.Id, UpdateBy: user.Id},
	}
	if len(loginLog.Client) > 255 {
		loginLog.Client = loginLog.Client[0:254]
	}
	if len(loginLog.Os) > 255 {
		loginLog.Os = loginLog.Os[0:254]
	}
	if len(loginLog.Detail) > 255 {
		loginLog.Detail = loginLog.Detail[0:254]
	}
	models.DB.Create(loginLog)
	return ""
}

func InsertLoginLogCallback(username, ip, userAgent, detail string, status bool, kind string, result *message.Result) {
	if result.IsSuccess() {
		res, _ := jsoniter.MarshalToString(result)
		logger.Info(result.Id + " is success, result: " + res)
	} else {
		logger.Error(result.Id + " is failed")
	}
}

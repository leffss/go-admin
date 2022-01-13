package app

import (
	"errors"
	"github.com/leffss/go-admin/pkg/logging"
	"go.uber.org/zap"
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"

	"github.com/leffss/go-admin/pkg/codemsg"
)

var logger = logging.GetLogger()

// BindJsonAndValid BindAndValid binds and validates data
func BindJsonAndValid(c *gin.Context, form interface{}) (int, int, error) {
	//return http.StatusBadRequest, codemsg.InvalidParams, errors.New("valid check not passed")

	// BindJSON 如果绑定错误，会返回 400，会导致 gin 打印的 request 日志，只有日志等级和时间戳，例如
	// {"L":"error","T":"2021-12-24T12:26:35.480+0800"}
	err := c.BindJSON(form)
	if err != nil {
		return http.StatusBadRequest, codemsg.ErrorInvalidParams, err
	}

	valid := validation.Validation{}
	check, err := valid.Valid(form)
	if err != nil {
		return http.StatusInternalServerError, codemsg.Error, err
	}
	if !check {
		MarkErrors(valid.Errors)
		return http.StatusBadRequest, codemsg.ErrorInvalidParams, errors.New("valid check not passed")
	}
	return http.StatusOK, codemsg.Success, nil
}

// MarkErrors logs error logs
func MarkErrors(errors []*validation.Error) {
	for _, err := range errors {
		//logging.Logger.Error(fmt.Sprintf(err.Key, "(", err.Value.(string), ")", err.Message))
		logger.Error(err.Message, zap.String("key", err.Key), zap.String("value", err.Value.(string)))
	}
	return
}

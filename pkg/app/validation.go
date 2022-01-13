package app

import (
	"errors"
	"go.uber.org/zap"
	"net/http"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/leffss/go-admin/pkg/codemsg"
	"github.com/leffss/go-admin/pkg/common"
)

// BindJsonAndValid2 BindAndValid binds and validates data
func BindJsonAndValid2(c *gin.Context, form interface{}, validatorFunc ...validator.Func) (int, int, error) {
	//return http.StatusBadRequest, codemsg.InvalidParams, errors.New("valid check not passed")

	// BindJSON 如果绑定错误，会返回 400，会导致 gin 打印的 request 日志，只有日志等级和时间戳，例如
	// {"L":"error","T":"2021-12-24T12:26:35.480+0800"}
	err := c.BindJSON(form)
	if err != nil {
		return http.StatusBadRequest, codemsg.ErrorInvalidParams, err
	}

	valid := validator.New()
	if len(validatorFunc) != 0 {
		for _, v := range validatorFunc {
			funcName := strings.TrimPrefix(filepath.Ext(runtime.FuncForPC(reflect.ValueOf(v).Pointer()).Name()), ".")
			if err := valid.RegisterValidation(funcName, v); err != nil {
				logger.Error(err.Error())
			}
		}
	}

	validErrs := valid.Struct(form)
	if validErrs != nil {
		if _, ok := validErrs.(*validator.InvalidValidationError); ok {
			return http.StatusBadRequest, codemsg.ErrorInvalidParams, errors.New("valid rule is error")
		}
		var errs []validator.FieldError
		for _, err := range validErrs.(validator.ValidationErrors) {
			errs = append(errs, err)
		}
		MarkErrors2(errs)
		return http.StatusBadRequest, codemsg.ErrorInvalidParams, errors.New("valid check not passed")
	}
	return http.StatusOK, codemsg.Success, nil
}

func MarkErrors2(errors []validator.FieldError) {
	for _, err := range errors {
		//logging.Logger.Error(fmt.Sprintf(err.Key, "(", err.Value.(string), ")", err.Message))
		logger.Error(err.Error(), zap.String("key", err.Field()), zap.String("value", common.StrVal(err.Value())))
	}
	return
}

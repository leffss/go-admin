package app

import (
	"github.com/gin-gonic/gin"
	"github.com/leffss/go-admin/pkg/codemsg"
)

type Gin struct {
	C *gin.Context
}

type Response struct {
	Code int         `json:"code,omitempty"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

// Response setting gin.JSON
func (g *Gin) Response(httpCode, errCode int, data interface{}) {
	g.C.JSON(httpCode, Response{
		Code: errCode,
		Msg:  codemsg.GetMsg(errCode),
		Data: data,
	})
	return
}

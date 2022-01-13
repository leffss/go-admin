package jwt

import (
	"fmt"
	"github.com/leffss/go-admin/pkg/app"
	"go.uber.org/zap"
	"net/http"
	"strings"

	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/leffss/go-admin/pkg/codemsg"
	"github.com/leffss/go-admin/pkg/jwt"
	"github.com/leffss/go-admin/pkg/logging"
)

var logger = logging.GetLogger()

func JWT(key, prefix string, whiteList ...string) gin.HandlerFunc {
	whiteLists := map[string]bool{}
	if len(whiteList) > 0 {
		for _, v := range whiteList {
			whiteLists[v] = true
		}
	}
	return func(c *gin.Context) {
		if _, ok := whiteLists[c.FullPath()]; ok {
			// 白名单不验证 token
			logger.Debug("token white list: " + c.FullPath())
			c.Next()
		} else {
			var code int
			var data interface{}
			code = codemsg.Success
			token := strings.TrimSpace(c.Request.Header.Get(key))

			if token == "" {
				code = codemsg.ErrorNoToken
			} else {
				tmp := strings.Split(token, prefix)
				if len(tmp) == 1 {
					code = codemsg.ErrorTokenFormat
					logger.Warn(token + " " + codemsg.GetMsg(code), zap.String("type", "token"))
				} else {
					token = strings.TrimSpace(tmp[1])
					_, err := jwt.ParseToken(token)
					if err != nil {
						//logging.Logger.Error(token + " " + err.Error(), zap.String("type", "token"))
						logger.Error(token + " " + err.Error(), zap.String("type", "token"))
						switch err.(*jwtGo.ValidationError).Errors {
						case jwtGo.ValidationErrorExpired:
							code = codemsg.ErrorTokenExpire
						default:
							code = codemsg.ErrorToken
						}
					}
				}
			}

			if code != codemsg.Success {
				c.JSON(http.StatusUnauthorized, app.Response{
					Code: code,
					Msg:  codemsg.GetMsg(code),
					Data: data,
				})
				c.Abort()
				return
			}

			c.Set("token", token)

			c.Next()
		}
	}
}

func JwtWs(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var token string
		code = codemsg.Success
		originProtocols := strings.TrimSpace(c.Request.Header.Get(key))
		if originProtocols == "" {
			code = codemsg.ErrorNoToken
		} else {
			protocols := strings.Split(originProtocols, ",")
			for i := range protocols {
				protocols[i] = strings.TrimSpace(protocols[i])
			}
			token = protocols[0]
			_, err := jwt.ParseToken(token)
			if err != nil {
				//logging.Logger.Error(token + " " + err.Error(), zap.String("type", "ws_token"))
				logger.Error(token + " " + err.Error(), zap.String("type", "ws_token"))
				switch err.(*jwtGo.ValidationError).Errors {
				case jwtGo.ValidationErrorExpired:
					code = codemsg.ErrorTokenExpire
				default:
					code = codemsg.ErrorToken
				}
			}
		}

		if code != codemsg.Success {
			upGrader := websocket.Upgrader{
				// cross origin domain
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
				// 处理 Sec-WebSocket-Protocol Header
				//Subprotocols: []string{c.Request.Header.Get("Sec-WebSocket-Protocol")},
				Subprotocols: websocket.Subprotocols(c.Request),
				ReadBufferSize: 1024,
				WriteBufferSize: 1024,
			}
			ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
			if err == nil {
				_ = ws.WriteMessage(websocket.BinaryMessage,
					[]byte(fmt.Sprintf("code: %d msg: %s\r\n", code, codemsg.GetMsg(code))))
				_ = ws.Close()
			}
			c.Abort()
			return
		}

		c.Set("token", token)
		c.Next()
	}
}

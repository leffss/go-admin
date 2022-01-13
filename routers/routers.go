package routers

import (
	"github.com/leffss/go-admin/middleware/jwt"
	zapMiddleware "github.com/leffss/go-admin/middleware/zap"
	"github.com/leffss/go-admin/pkg/logging"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"go.uber.org/zap"
	"net/http"
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	adminRouter "github.com/leffss/go-admin/app/admin/router"
	authRouter "github.com/leffss/go-admin/app/auth/router"
	profileRouter "github.com/leffss/go-admin/app/profile/router"
	_ "github.com/leffss/go-admin/docs"
	"github.com/leffss/go-admin/middleware/cors"
	"github.com/leffss/go-admin/pkg/setting"
)

var serverSetting = setting.GetServerSetting()
var appSetting = setting.GetAppSetting()

func InitRouter() *gin.Engine {
	gin.DisableConsoleColor()
	router := gin.New()
	router.Use(requestid.New())
	setLoggerRecovery(router)
	setCors(router)

	if serverSetting.Swagger {
		// Swagger api 文档
		router.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// apiGroup 需要验证 token，注意白名单
	apiGroup := router.Group("/api", jwt.JWT(appSetting.JwtHeaderKey, appSetting.JwtHeaderPrefixKey, appSetting.JwtWhiteLists...))

	// 初始化各 app 的 router
	authApiGroup := apiGroup.Group("/oauth")
	authRouter.InitRouter(authApiGroup)

	adminApiGroup := apiGroup.Group("/admin")
	adminRouter.InitRouter(adminApiGroup)

	profileApiGroup := apiGroup.Group("/profile")
	profileRouter.InitRouter(profileApiGroup)

	// 捕获 405 错误
	router.HandleMethodNotAllowed = true
	router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"code": 405,
			"msg": "Method Not Allowed",
		})
	})

	return router
}

func setLoggerRecovery(router *gin.Engine)  {
	//if serverSetting.GinLog {
	//	router.Use(gin.Logger())
	//}
	//router.Use(gin.Recovery())

	//logger, _ := zap.NewProduction()
	logger := logging.RequestLogger
	//if serverSetting.GinLog {
	//	router.Use(ginzap.Ginzap(logger, time.RFC3339Nano, false))
	//}
	//router.Use(ginzap.RecoveryWithZap(logger, true))

	// 自定义 zap logger 中间件
	// 基于 https://github.com/gin-contrib/zap 的 #16 pr 修改 https://github.com/gin-contrib/zap/pull/16
	if serverSetting.GinLog {
		router.Use(zapMiddleware.Logger(logger, zapMiddleware.WithTimeFormat(time.RFC3339Nano), zapMiddleware.WithUTC(true),
			zapMiddleware.WithBody(true),
			zapMiddleware.WithCustomFields(
				// 需要使用 github.com/gin-contrib/requestid 中间件后，才能获取到
				//func(c *gin.Context) zap.Field { return zap.String("x-request-id", c.Writer.Header().Get("X-Request-ID")) },
				func(c *gin.Context) zap.Field { return zap.String("x-request-id", requestid.Get(c)) },
			),
			))
	}

	router.Use(zapMiddleware.RecoveryWithZap(logger, true))
}

func setCors(router *gin.Engine)  {
	router.Use(cors.Cors())
}

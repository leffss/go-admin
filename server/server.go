package server

import (
	"context"
	"fmt"
	"github.com/leffss/go-admin/pkg/logging"
	"net/http"

	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
	"github.com/leffss/go-admin/pkg/setting"
	"github.com/leffss/go-admin/routers"
)

var server *http.Server
var serverSetting = setting.GetServerSetting()
var logger = logging.GetLogger()

func StartHttpServer() {
	gin.SetMode(serverSetting.RunMode)
	routersInit := routers.InitRouter()
	if serverSetting.Ginpprof {
		ginpprof.Wrap(routersInit)
	}
	endPoint := fmt.Sprintf(":%d", serverSetting.Port)
	server = &http.Server{
		Addr:           endPoint,
		Handler:        routersInit,
		ReadTimeout:    serverSetting.ReadTimeout,
		WriteTimeout:   serverSetting.WriteTimeout,
		//WriteTimeout:   3000 * time.Second,
		MaxHeaderBytes: serverSetting.MaxHeaderBytes,
	}
	logger.Info(fmt.Sprintf("start http server listening %s", endPoint))
	go func() {
		// 服务连接
		if serverSetting.SSL {
			// https
			if err := server.ListenAndServeTLS(serverSetting.SSLCa, serverSetting.SSLKey); err != nil && err != http.ErrServerClosed {
				logger.Fatal(fmt.Sprintf("Server Start Error: %s", err.Error()))
			}
		} else {
			// http
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logger.Fatal(fmt.Sprintf("Server Start Error: %s", err))
			}
		}
	}()
}

func StopHttpServer()  {
	ctx, cancel := context.WithTimeout(context.Background(), serverSetting.CloseTimeout)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal(fmt.Sprintf("Server Shutdown Error: %s", err))
	}
	logger.Info("Server Shutdown")
}

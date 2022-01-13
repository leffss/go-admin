package cmd

import (
	"fmt"
	"github.com/leffss/go-admin/models"
	"github.com/leffss/go-admin/pkg/encryption"
	"github.com/leffss/go-admin/pkg/gredis"
	"github.com/leffss/go-admin/pkg/jwt"
	"github.com/leffss/go-admin/pkg/logging"
	_ "github.com/leffss/go-admin/pkg/logging"
	_ "github.com/leffss/go-admin/pkg/setting"
	"github.com/leffss/go-admin/server"
	"github.com/leffss/go-admin/tasks"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
)

var logger = logging.GetLogger()

func initPkg() {
	// 加密
	encryption.Setup()
	// 加密解密测试
	//enc, _ := encryption.AesCtrEncryptToString([]byte("123456"))
	//fmt.Println(enc)
	//dec, _ := encryption.AesCtrStringToDecrypt(enc)
	//fmt.Println(string(dec))

	// jwt token
	jwt.Setup()

	// gredis
	if err := gredis.Setup(); err != nil {
		logger.Fatal(err.Error())
	}
	// 测试 redis
	//_ = gredis.Set("leffss", "123456", 30)
	//tmp, _ := gredis.Get("leffss")
	//fmt.Println(string(tmp))

	// database
	if err := models.Setup(); err != nil {
		logger.Fatal(err.Error())
	}

	// 插入测试数据
	//logger.Info("插入测试数据")
	//models.InitDatabase("123456")

	tasks.Setup()
	// 后台任务测试
	//tasks.SendTask("adminTaskTest", 1)
	//tasks.SendTask("hostTaskTest", 2)
}

func start()  {
	if mode == "all" {
		logger.Info("mode [all]")
		tasks.StartTasksServer()
		server.StartHttpServer()
	} else if mode == "http" {
		logger.Info("mode [http]")
		server.StartHttpServer()
	} else if mode == "task" {
		logger.Info("mode [task]")
		tasks.StartTasksServer()
	} else {
		logger.Fatal("-m or --mode just support [all, http, task]")
	}
}

func stop()  {
	if mode == "all" {
		server.StopHttpServer()

		if err := tasks.StopTasksServer(); err != nil {
			logger.Error(err.Error())
		}

		if err := models.CloseDB(); err != nil {
			logger.Error(err.Error())
		}

		if err := gredis.CloseRedis(); err != nil {
			logger.Error(err.Error())
		}
	} else if mode == "http" {
		server.StopHttpServer()

		if err := models.CloseDB(); err != nil {
			logger.Error(err.Error())
		}

		if err := gredis.CloseRedis(); err != nil {
			logger.Error(err.Error())
		}
	} else if mode == "task" {
		if err := tasks.StopTasksServer(); err != nil {
			logger.Error(err.Error())
		}

		if err := models.CloseDB(); err != nil {
			logger.Error(err.Error())
		}

		if err := gredis.CloseRedis(); err != nil {
			logger.Error(err.Error())
		}
	}

	logger.Info("exit")
	_ = logger.Sync()
}

var (
	rootCmd = &cobra.Command{
		Use:   "go-admin",
		Short: "Background management program written by gin web framework",
		Long: `Background management program written by gin web framework`,
		PreRun: func(cmd *cobra.Command, args []string) {
			// 初始化依赖包
			initPkg()
		},
		Run: func(cmd *cobra.Command, args []string) {
			start()

			// 等待中断信号以优雅地关闭服务器（可以设置超时时间）
			quit := make(chan os.Signal, 1)
			//signal.Notify(quit, os.Interrupt)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
			<-quit

			stop()
		},
	}
)

func init()  {
	rootCmd.Flags().StringVarP(&mode, "mode", "m", "all", "run mode [all, http, task]")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// 错误信息打印到控制台
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

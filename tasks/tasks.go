package tasks

import (
	"context"
	"github.com/leffss/go-admin/pkg/logging"
	"time"

	ytask "github.com/gojuukaze/YTask/v2"
	"github.com/gojuukaze/YTask/v2/server"
	adminTask "github.com/leffss/go-admin/app/admin/task"
	authTask "github.com/leffss/go-admin/app/auth/task"
	"github.com/leffss/go-admin/pkg/setting"
)

var (
	Server server.Server
	Client server.Client
	Tasks map[string]interface{}
	TasksCallback map[string]interface{}
	//redisSetting = setting.GetRedisSetting()
	tasksSetting = setting.GetTasksSetting()
	logger = logging.GetLogger()
	yLogger = NewMyLogger(logger)
)

func Setup() {
	Tasks = map[string]interface{}{
		"adminTaskTest": adminTask.Test,
		"authTaskInsertLoginLog": authTask.InsertLoginLog,
	}
	// 任务回调函数，没有则设置 nil，需要和 Tasks 对应
	TasksCallback = map[string]interface{}{
		"adminTaskTest": nil,
		"authTaskInsertLoginLog": authTask.InsertLoginLogCallback,
	}
	//broker := ytask.Broker.NewRedisBroker(redisSetting.Host, redisSetting.Port, redisSetting.Password, 0, redisSetting.MaxActive)
	//backend := ytask.Backend.NewRedisBackend(redisSetting.Host, redisSetting.Port, redisSetting.Password, 0, redisSetting.MaxActive)
	//yLogger := ytask.Logger.NewYTaskLogger()

	broker := NewMyRedisBroker(5 * time.Second)	// 设置时间不能太长，否则正常关闭程序时会等待很久时间
	backend := NewMyRedisBackend()

	Server = ytask.Server.NewServer(
		ytask.Config.Broker(broker),
		ytask.Config.Backend(backend),		// 可以不设置 backend
		ytask.Config.Logger(yLogger),		// 可以不设置 logger
		ytask.Config.Debug(tasksSetting.Debug),
		ytask.Config.StatusExpires(tasksSetting.StatusExpires),
		ytask.Config.ResultExpires(tasksSetting.ResultExpires),
	)
	Client = Server.GetClient()
}

func StartTasksServer()  {
	for taskName, task := range Tasks {
		// v2.4 支持设置 callback 回调函数
		if TasksCallback[taskName] != nil {
			Server.Add(tasksSetting.Group, taskName, task, TasksCallback[taskName])
		} else {
			Server.Add(tasksSetting.Group, taskName, task)
		}
	}

	// v2.2 开始支持运行多个 group
	//Server.Run("group1", 3)
	//Server.Run("group2", 3)
	// 如果要使用延时任务，把 enableDelayServer 设为 true
	//Server.Run("group", 3,true)

	Server.Run(tasksSetting.Group, tasksSetting.NumWorkers, true)
}

func SendTask(taskName string, args ...interface{}) (string, error) {
	// 设置普通任务
	return Client.Send(tasksSetting.Group, taskName, args...)
}

func SendRetryTask(taskName string, retryCount int, args ...interface{}) (string, error) {
	// 设置重试任务
	return Client.SetTaskCtl(Client.RetryCount, retryCount).Send(tasksSetting.Group, taskName, args...)
}

func SendDelayTask(taskName string, delay time.Duration, args ...interface{}) (string, error) {
	// 设置延迟任务
	return Client.SetTaskCtl(Client.RunAfter, delay).Send(tasksSetting.Group, taskName, args...)
}

func SendRunAtTask(taskName string, runAt time.Time, args ...interface{}) (string, error) {
	// 设置时间任务
	return Client.SetTaskCtl(Client.RunAt, runAt).Send(tasksSetting.Group, taskName, args...)
}

func StopTasksServer() error {
	return Server.Shutdown(context.Background())
}

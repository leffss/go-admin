package logging

import (
	"fmt"
	"os"

	"github.com/leffss/go-admin/pkg/setting"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger
var RequestLogger *zap.Logger

func GetLogger() *zap.Logger {
	return Logger
}

//func GetRequestLogger() *zap.Logger {
//	return RequestLogger
//}

// Setup initialize the log instance
func Setup() {
	appSetting := setting.GetAppSetting()
	serverHook := lumberjack.Logger{
		Filename:   fmt.Sprintf(
			"%s%s%s.%s",appSetting.RuntimeRootPath, appSetting.LogSavePath,
			appSetting.LogServerName, appSetting.LogFileExt), // 日志文件路径
		MaxSize:    appSetting.LogMaxSize,                      // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: appSetting.LogMaxBackups,                       // 日志文件最多保存多少个备份
		MaxAge:     appSetting.LogMaxAge,                        // 文件最多保存多少天
		Compress:   appSetting.LogCompress,                     // 是否压缩
		LocalTime:  true,
	}

	requsetHook := lumberjack.Logger{
		Filename:   fmt.Sprintf(
			"%s%s%s.%s",appSetting.RuntimeRootPath, appSetting.LogSavePath,
			appSetting.LogRequestName, appSetting.LogFileExt),
		MaxSize:    appSetting.LogMaxSize,
		MaxBackups: appSetting.LogMaxBackups,
		MaxAge:     appSetting.LogMaxAge,
		Compress:   appSetting.LogCompress,
		LocalTime:  true,
	}

	serverEncoderConfig := zapcore.EncoderConfig{
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		//EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	requsetEncoderConfig := zapcore.EncoderConfig{
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 设置日志级别
	serverAtomicLevel := zap.NewAtomicLevel()
	if appSetting.LogLevel == "debug" {
		serverAtomicLevel.SetLevel(zap.DebugLevel)
	} else if appSetting.LogLevel == "info" {
		serverAtomicLevel.SetLevel(zap.InfoLevel)
	} else if appSetting.LogLevel == "warn" {
		serverAtomicLevel.SetLevel(zap.WarnLevel)
	} else if appSetting.LogLevel == "error" {
		serverAtomicLevel.SetLevel(zap.ErrorLevel)
	} else if appSetting.LogLevel == "panic" {
		serverAtomicLevel.SetLevel(zap.PanicLevel)
	} else if appSetting.LogLevel == "fatal" {
		serverAtomicLevel.SetLevel(zap.FatalLevel)
	} else {
		serverAtomicLevel.SetLevel(zap.ErrorLevel)
	}

	requestAtomicLevel := zap.NewAtomicLevel()
	requestAtomicLevel.SetLevel(zap.InfoLevel)

	var serverCore zapcore.Core
	if appSetting.LogConsoleOutput {
		serverCore = zapcore.NewCore(
			zapcore.NewJSONEncoder(serverEncoderConfig),                                           // 编码器配置
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stderr), zapcore.AddSync(&serverHook)), // 打印到控制台和文件
			serverAtomicLevel,                                                                     // 日志级别
		)
	} else if appSetting.LogFileOutput {
		serverCore = zapcore.NewCore(
			zapcore.NewJSONEncoder(serverEncoderConfig),                                           // 编码器配置
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(&serverHook)),							   // 仅输出到文件
			serverAtomicLevel,                                                                     // 日志级别
		)
	}

	var requestCore zapcore.Core
	if appSetting.LogConsoleOutput {
		requestCore = zapcore.NewCore(
			zapcore.NewJSONEncoder(requsetEncoderConfig),                                           // 编码器配置
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stderr), zapcore.AddSync(&requsetHook)), // 打印到控制台和文件
			requestAtomicLevel,                                                                     // 日志级别
		)
	} else if appSetting.LogFileOutput {
		requestCore = zapcore.NewCore(
			zapcore.NewJSONEncoder(requsetEncoderConfig),                                           // 编码器配置
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(&requsetHook)), 							// 仅输出到文件
			requestAtomicLevel,                                                                     // 日志级别
		)
	}

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()

	// 设置初始化字段
	//filed := zap.Fields(zap.String("service", "go-admin"))

	// 构造日志
	//Logger = zap.New(serverCore, caller, development, filed)
	Logger = zap.New(serverCore, caller, development)

	//RequestLogger = zap.New(requestCore, caller, development, filed)
	RequestLogger = zap.New(requestCore, caller, development)

	//logger.Info("log 初始化成功")
	//logger.Info("无法获取网址",
	//	zap.String("url", "http://www.baidu.com"),
	//	zap.Int("attempt", 3),
	//	zap.Duration("backoff", time.Second))
}

func init()  {
	Setup()
}

package tasks

import (
	"fmt"
	"github.com/gojuukaze/YTask/v2/log"
	"go.uber.org/zap"
)

type MyLogger struct {
	logger *zap.Logger
}

func NewMyLogger(logger *zap.Logger) *MyLogger {
	return &MyLogger{
		logger: logger,
	}
}

func (ml *MyLogger) Debug(msg string) {
	ml.logger.Debug(msg)
}

func (ml *MyLogger) DebugWithField(msg string, key string, val interface{}) {
	ml.logger.Debug(msg, zap.String(key, fmt.Sprint(val)))
}

func (ml *MyLogger) Info(msg string) {
	ml.logger.Info(msg)
}

func (ml *MyLogger) InfoWithField(msg string, key string, val interface{}) {
	ml.logger.Info(msg, zap.String(key, fmt.Sprint(val)))
}

func (ml *MyLogger) Warn(msg string) {
	ml.logger.Warn(msg)
}

func (ml *MyLogger) WarnWithField(msg string, key string, val interface{}) {
	ml.logger.Warn(msg, zap.String(key, fmt.Sprint(val)))
}

func (ml *MyLogger) Error(msg string) {
	ml.logger.Error(msg)
}

func (ml *MyLogger) ErrorWithField(msg string, key string, val interface{}) {
	ml.logger.Error(msg, zap.String(key, fmt.Sprint(val)))
}

func (ml *MyLogger) Fatal(msg string) {
	ml.logger.Fatal(msg)
}

func (ml *MyLogger) FatalWithField(msg string, key string, val interface{}) {
	ml.logger.Fatal(msg, zap.String(key, fmt.Sprint(val)))
}

func (ml *MyLogger) Panic(msg string) {
	ml.logger.Panic(msg)
}

func (ml *MyLogger) PanicWithField(msg string, key string, val interface{}) {
	ml.logger.Panic(msg, zap.String(key, fmt.Sprint(val)))
}

func (ml *MyLogger) SetLevel(level string)  {
	switch level {
	case "debug":

	case "info":

	case "warn":

	case "error":

	case "fatal":

	case "panic":

	default:

	}
}

func (ml *MyLogger) Clone() log.LoggerInterface {
	return ml
}

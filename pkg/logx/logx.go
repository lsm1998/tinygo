package logx

import (
	"github.com/lsm1998/tinygo"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func SetLog(config Config) {
	// 设置日志级别
	logrus.SetLevel(config.GetLevel())
	// 设置日志格式
	logrus.SetFormatter(config.GetFormatter())
	// 设置日志输出
	if tinygo.AppMod == "local" {
		logrus.SetOutput(logrus.StandardLogger().Out) // 默认为标准输出
	} else {
		logrus.SetOutput(&lumberjack.Logger{
			Filename:   config.GetFilename(),
			MaxSize:    config.GetMaxSize(),    // 最大XMB
			MaxBackups: config.GetMaxBackups(), // 最多保留X个旧日志文件
			MaxAge:     config.GetMaxAge(),     // 保留X天的日志文件
			Compress:   config.GetCompress(),   // 是否压缩旧日志文件
		})
	}
}

func Errorf(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}

func Error(args ...interface{}) {
	logrus.Error(args...)
}

func Warnf(format string, args ...interface{}) {
	logrus.Warnf(format, args...)
}

func Warn(args ...interface{}) {
	logrus.Warn(args...)
}

func Infof(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

func Info(args ...interface{}) {
	logrus.Info(args...)
}

func Debugf(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}

func Debug(args ...interface{}) {
	logrus.Debug(args...)
}

func Tracef(format string, args ...interface{}) {
	logrus.Tracef(format, args...)
}

func Trace(args ...interface{}) {
	logrus.Trace(args...)
}

func Fatalf(format string, args ...interface{}) {
	logrus.Fatalf(format, args...)
}

func Fatal(args ...interface{}) {
	logrus.Fatal(args...)
}

func Panicf(format string, args ...interface{}) {
	logrus.Panicf(format, args...)
}

func Panic(args ...interface{}) {
	logrus.Panic(args...)
}

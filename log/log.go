package log

import (
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func LogInit(logLevel string) {
	logger = logrus.New()

	level := logrus.DebugLevel
	logger.SetLevel(logrus.ErrorLevel)
	switch {
	case logLevel == "debug":
		level = logrus.DebugLevel
	case logLevel == "info":
		level = logrus.InfoLevel
	case logLevel == "error":
		level = logrus.ErrorLevel
	default:
		level = logrus.DebugLevel
	}
	logger.Formatter = &logrus.JSONFormatter{}
	logger.SetLevel(level)

}

func LogInfo(v ...interface{}) {
	logger.Info(v...)
}

func LogError(v ...interface{}) {
	logger.Error(v...)
}

func Debug(v ...interface{}) {
	logger.Debug(v...)
}

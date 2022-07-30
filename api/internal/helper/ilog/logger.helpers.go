package ilog

import "github.com/sirupsen/logrus"

func ErrorlnIf(err error, log StdLogger) {
	if err != nil {
		log.Errorln(err)
	}
}

func MakeLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		DisableTimestamp: true,
		PrettyPrint:      true,
	})

	return logger
}

func MakeLoggerWithTag(tag string) logrus.FieldLogger {
	logger := MakeLogger()
	return logger.WithFields(logrus.Fields{
		"tag": tag,
	})
}

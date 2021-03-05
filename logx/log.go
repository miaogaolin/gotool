package logx

import (
	"github.com/sirupsen/logrus"
	"io"
)

var log = logrus.New()


func init() {
	log.SetFormatter(&logrus.JSONFormatter{})
}

func Out() io.Writer {
	return log.Out
}

func SetOutput(output io.Writer)  {
	log.SetOutput(output)
}

func Label(label string) *logrus.Entry {
	return log.WithField("label", label)
}

func WithFields(fields map[string]interface{}) *logrus.Entry {
	return log.WithFields(fields)
}

func WithField(key string, value interface{}) *logrus.Entry {
	return log.WithField(key, value)
}

func Error(args ...interface{}) {
	log.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func Info(args ...interface{}) {
	log.Info(args...)
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func Warn(args ...interface{}) {
	log.Warn(args)
}

func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

func Fatal(args ...interface{}) {
	log.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

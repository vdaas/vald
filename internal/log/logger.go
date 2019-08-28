package log

import "sync"

type Logger interface {
	Info(vals ...interface{})
	Infof(format string, vals ...interface{})
	Debug(vals ...interface{})
	Debugf(format string, vals ...interface{})
	Warn(vals ...interface{})
	Warnf(format string, vals ...interface{})
	Error(vals ...interface{})
	Errorf(format string, vals ...interface{})
	Fatal(vals ...interface{})
	Fatalf(format string, vals ...interface{})
}

var (
	logger Logger
	once   sync.Once
)

func Init(l Logger) {
	once.Do(func() {
		logger = l
	})
}

func Info(vals ...interface{}) {
	logger.Info(vals...)
}
func Infof(format string, vals ...interface{}) {
	logger.Infof(format, vals...)
}

func Debug(vals ...interface{}) {
	logger.Debug(vals...)
}

func Debugf(format string, vals ...interface{}) {
	logger.Debugf(format, vals...)
}

func Warn(vals ...interface{}) {
	logger.Warn(vals...)
}

func Warnf(format string, vals ...interface{}) {
	logger.Warnf(format, vals...)
}

func Error(vals ...interface{}) {
	logger.Error(vals...)
}

func Errorf(format string, vals ...interface{}) {
	logger.Errorf(format, vals...)
}

func Fatal(vals ...interface{}) {
	logger.Fatal(vals...)
}

func Fatalf(format string, vals ...interface{}) {
	logger.Fatalf(format, vals...)
}

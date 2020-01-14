package info

import (
	"github.com/vdaas/vald/internal/log"
)

type loggerMock struct {
	log.Logger
	InfoFunc func(vals ...interface{})
}

func (l *loggerMock) Info(vals ...interface{}) {
	l.InfoFunc(vals...)
}

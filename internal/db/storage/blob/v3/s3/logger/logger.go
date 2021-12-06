package logger

import (
	"github.com/aws/smithy-go/logging"

	"github.com/vdaas/vald/internal/log"
)

type logger struct{}

func New() logging.Logger {
	return new(logger)
}

func (l logger) Logf(classification logging.Classification, format string, v ...interface{}) {
	switch classification {
	case logging.Warn:
		log.Warnf(format, v...)
	case logging.Debug:
		log.Debugf(format, v...)
	default:
		log.Infof(format, v...)
	}
}

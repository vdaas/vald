package log

import (
	"github.com/kpango/glg"
)

type glglogger struct {
	log *glg.Glg
}

// New returns a new glglogger instance.
func NewGlg(g *glg.Glg) Logger {
	return &glglogger{
		log: g,
	}
}

func DefaultGlg() Logger {
	return &glglogger{
		log: glg.Get(),
	}
}

func (l *glglogger) Info(vals ...interface{}) {
	l.log.Info(vals...)
}

func (l *glglogger) Infof(format string, vals ...interface{}) {
	l.log.Infof(format, vals...)
}

func (l *glglogger) Debug(vals ...interface{}) {
	l.log.Debug(vals...)
}

func (l *glglogger) Debugf(format string, vals ...interface{}) {
	l.log.Debugf(format, vals...)
}

func (l *glglogger) Warn(vals ...interface{}) {
	l.log.Warn(vals...)
}

func (l *glglogger) Warnf(format string, vals ...interface{}) {
	l.log.Warnf(format, vals...)
}

func (l *glglogger) Error(vals ...interface{}) {
	l.log.Error(vals...)
}

func (l *glglogger) Errorf(format string, vals ...interface{}) {
	l.log.Errorf(format, vals...)
}

func (l *glglogger) Fatal(vals ...interface{}) {
	l.log.Fatal(vals...)
}

func (l *glglogger) Fatalf(format string, vals ...interface{}) {
	l.log.Fatalf(format, vals...)
}

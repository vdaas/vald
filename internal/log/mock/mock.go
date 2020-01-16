package mock

type Logger struct {
	InfoFunc   func(vals ...interface{})
	InfofFunc  func(format string, vals ...interface{})
	DebugFunc  func(vals ...interface{})
	DebugfFunc func(format string, vals ...interface{})
	WarnFunc   func(vals ...interface{})
	WarnfFunc  func(format string, vals ...interface{})
	ErrorFunc  func(vals ...interface{})
	ErrorfFunc func(format string, vals ...interface{})
	FatalFunc  func(vals ...interface{})
	FatalfFunc func(format string, vals ...interface{})
}

func (l *Logger) Info(vals ...interface{}) {
	l.InfoFunc(vals...)
}

func (l *Logger) Infof(format string, vals ...interface{}) {
	l.InfofFunc(format, vals...)
}

func (l *Logger) Debug(vals ...interface{}) {
	l.DebugFunc(vals...)
}

func (l *Logger) Debugf(format string, vals ...interface{}) {
	l.DebugfFunc(format, vals...)
}

func (l *Logger) Warn(vals ...interface{}) {
	l.WarnFunc(vals...)
}

func (l *Logger) Warnf(format string, vals ...interface{}) {
	l.WarnfFunc(format, vals...)
}

func (l *Logger) Error(vals ...interface{}) {
	l.ErrorFunc(vals...)
}

func (l *Logger) Errorf(format string, vals ...interface{}) {
	l.ErrorfFunc(format, vals...)
}

func (l *Logger) Fatal(vals ...interface{}) {
	l.FatalFunc(vals...)
}

func (l *Logger) Fatalf(format string, vals ...interface{}) {
	l.FatalfFunc(format, vals...)
}

package log

type loggerMock struct {
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

func (lm *loggerMock) Info(vals ...interface{}) {
	lm.InfoFunc(vals...)
}

func (lm *loggerMock) Infof(format string, vals ...interface{}) {
	lm.InfofFunc(format, vals...)
}

func (lm *loggerMock) Debug(vals ...interface{}) {
	lm.DebugFunc(vals...)
}

func (lm *loggerMock) Debugf(format string, vals ...interface{}) {
	lm.DebugfFunc(format, vals...)
}

func (lm *loggerMock) Warn(vals ...interface{}) {
	lm.WarnFunc(vals...)
}

func (lm *loggerMock) Warnf(format string, vals ...interface{}) {
	lm.WarnfFunc(format, vals...)
}

func (lm *loggerMock) Error(vals ...interface{}) {
	lm.ErrorFunc(vals...)
}

func (lm *loggerMock) Errorf(format string, vals ...interface{}) {
	lm.ErrorfFunc(format, vals...)
}

func (lm *loggerMock) Fatal(vals ...interface{}) {
	lm.FatalFunc(vals...)
}

func (lm *loggerMock) Fatalf(format string, vals ...interface{}) {
	lm.FatalfFunc(format, vals...)
}

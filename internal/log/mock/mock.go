package mock

type MockLogger struct {
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

func (ml *MockLogger) Info(vals ...interface{}) {
	ml.InfoFunc(vals...)
}

func (ml *MockLogger) Infof(format string, vals ...interface{}) {
	ml.InfofFunc(format, vals...)
}

func (ml *MockLogger) Debug(vals ...interface{}) {
	ml.DebugFunc(vals...)
}

func (ml *MockLogger) Debugf(format string, vals ...interface{}) {
	ml.DebugfFunc(format, vals...)
}

func (ml *MockLogger) Warn(vals ...interface{}) {
	ml.WarnFunc(vals...)
}

func (ml *MockLogger) Warnf(format string, vals ...interface{}) {
	ml.WarnfFunc(format, vals...)
}

func (ml *MockLogger) Error(vals ...interface{}) {
	ml.ErrorFunc(vals...)
}

func (ml *MockLogger) Errorf(format string, vals ...interface{}) {
	ml.ErrorfFunc(format, vals...)
}

func (ml *MockLogger) Fatal(vals ...interface{}) {
	ml.FatalFunc(vals...)
}

func (ml *MockLogger) Fatalf(format string, vals ...interface{}) {
	ml.FatalfFunc(format, vals...)
}

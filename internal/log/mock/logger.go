//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package mock

// Logger represents struct of each log level function.
type Logger struct {
	DebugFunc  func(vals ...interface{})
	DebugfFunc func(format string, vals ...interface{})
	InfoFunc   func(vals ...interface{})
	InfofFunc  func(format string, vals ...interface{})
	WarnFunc   func(vals ...interface{})
	WarnfFunc  func(format string, vals ...interface{})
	ErrorFunc  func(vals ...interface{})
	ErrorfFunc func(format string, vals ...interface{})
	FatalFunc  func(vals ...interface{})
	FatalfFunc func(format string, vals ...interface{})
}

// Debug calls DebugFunc of Logger.
func (l *Logger) Debug(vals ...interface{}) {
	l.DebugFunc(vals...)
}

// Debugf calls DebugfFunc of Logger.
func (l *Logger) Debugf(format string, vals ...interface{}) {
	l.DebugfFunc(format, vals...)
}

// Info calls InfoFunc of Logger.
func (l *Logger) Info(vals ...interface{}) {
	l.InfoFunc(vals...)
}

// Infof calls InfofFunc of Logger.
func (l *Logger) Infof(format string, vals ...interface{}) {
	l.InfofFunc(format, vals...)
}

// Warn calls WarnFunc of Logger.
func (l *Logger) Warn(vals ...interface{}) {
	l.WarnFunc(vals...)
}

// Warnf calls WarnfFunc of Logger.
func (l *Logger) Warnf(format string, vals ...interface{}) {
	l.WarnfFunc(format, vals...)
}

// Error calls ErrorFunc of Logger.
func (l *Logger) Error(vals ...interface{}) {
	l.ErrorFunc(vals...)
}

// Errorf calls ErrorfFunc of Logger.
func (l *Logger) Errorf(format string, vals ...interface{}) {
	l.ErrorfFunc(format, vals...)
}

// Fatal calls FatalFunc of Logger.
func (l *Logger) Fatal(vals ...interface{}) {
	l.FatalFunc(vals...)
}

// Fatalf calls FatalfFunc of Logger.
func (l *Logger) Fatalf(format string, vals ...interface{}) {
	l.FatalfFunc(format, vals...)
}

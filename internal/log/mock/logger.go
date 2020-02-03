//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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

func (l *Logger) Debug(vals ...interface{}) {
	l.DebugFunc(vals...)
}

func (l *Logger) Debugf(format string, vals ...interface{}) {
	l.DebugfFunc(format, vals...)
}

func (l *Logger) Info(vals ...interface{}) {
	l.InfoFunc(vals...)
}

func (l *Logger) Infof(format string, vals ...interface{}) {
	l.InfofFunc(format, vals...)
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

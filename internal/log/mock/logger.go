//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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
	DebugFunc  func(vals ...any)
	DebugfFunc func(format string, vals ...any)
	InfoFunc   func(vals ...any)
	InfofFunc  func(format string, vals ...any)
	WarnFunc   func(vals ...any)
	WarnfFunc  func(format string, vals ...any)
	ErrorFunc  func(vals ...any)
	ErrorfFunc func(format string, vals ...any)
	FatalFunc  func(vals ...any)
	FatalfFunc func(format string, vals ...any)
}

// Debug calls DebugFunc of Logger.
func (l *Logger) Debug(vals ...any) {
	l.DebugFunc(vals...)
}

// Debugf calls DebugfFunc of Logger.
func (l *Logger) Debugf(format string, vals ...any) {
	l.DebugfFunc(format, vals...)
}

// Debugd calls DebugfFunc of Logger.
func (l *Logger) Debugd(msg string, details ...any) {
	l.DebugfFunc(msg, details...)
}

// Info calls InfoFunc of Logger.
func (l *Logger) Info(vals ...any) {
	l.InfoFunc(vals...)
}

// Infof calls InfofFunc of Logger.
func (l *Logger) Infof(format string, vals ...any) {
	l.InfofFunc(format, vals...)
}

// Infod calls InfofFunc of Logger.
func (l *Logger) Infod(msg string, details ...any) {
	l.InfofFunc(msg, details...)
}

// Warn calls WarnFunc of Logger.
func (l *Logger) Warn(vals ...any) {
	l.WarnFunc(vals...)
}

// Warnf calls WarnfFunc of Logger.
func (l *Logger) Warnf(format string, vals ...any) {
	l.WarnfFunc(format, vals...)
}

// Warnd calls WarnfFunc of Logger.
func (l *Logger) Warnd(msg string, details ...any) {
	l.WarnfFunc(msg, details...)
}

// Error calls ErrorFunc of Logger.
func (l *Logger) Error(vals ...any) {
	l.ErrorFunc(vals...)
}

// Errorf calls ErrorfFunc of Logger.
func (l *Logger) Errorf(format string, vals ...any) {
	l.ErrorfFunc(format, vals...)
}

// Errord calls ErrorfFunc of Logger.
func (l *Logger) Errord(msg string, details ...any) {
	l.ErrorfFunc(msg, details...)
}

// Fatal calls FatalFunc of Logger.
func (l *Logger) Fatal(vals ...any) {
	l.FatalFunc(vals...)
}

// Fatalf calls FatalfFunc of Logger.
func (l *Logger) Fatalf(format string, vals ...any) {
	l.FatalfFunc(format, vals...)
}

// Fatald calls FatalfFunc of Logger.
func (l *Logger) Fatald(msg string, details ...any) {
	l.FatalfFunc(msg, details...)
}

func (*Logger) Close() error {
	return nil
}

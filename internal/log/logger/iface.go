//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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

package logger

type Logger interface {
	// Debug logs the vals at Debug level.
	Debug(vals ...interface{})

	// Debugf logs the formatted message at Debug level.
	Debugf(format string, vals ...interface{})

	// Debugd logs the message with details at Debug level.
	Debugd(msg string, details ...interface{})

	// Info logs the vals at Info level.
	Info(vals ...interface{})

	// Infof logs the formatted message at Info level.
	Infof(format string, vals ...interface{})

	// Infod logs the message with details at Info level.
	Infod(msg string, details ...interface{})

	// Warn logs the vals at Warn level.
	Warn(vals ...interface{})

	// Warnf logs the formatted message at Warn level.
	Warnf(format string, vals ...interface{})

	// Warnd logs the message with details at Warn level.
	Warnd(msg string, details ...interface{})

	// Error logs the vals at Error level.
	Error(vals ...interface{})

	// Errorf logs the formatted message at Error level.
	Errorf(format string, vals ...interface{})

	// Errord logs the message with details at Error level.
	Errord(msg string, details ...interface{})

	// Fatal logs the vals at Fatal level, then calls os.Exit(1).
	Fatal(vals ...interface{})

	// Fatalf logs the formatted message at Fatal level, then calls os.Exit(1).
	Fatalf(format string, vals ...interface{})

	// Fatald logs the message with details at Fatal level, then calls os.Exit(1).
	Fatald(msg string, details ...interface{})

	// Close calls finalizer of logger implementations.
	Close() error
}

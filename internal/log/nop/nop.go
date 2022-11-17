//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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
package nop

import (
	"github.com/vdaas/vald/internal/log/logger"
)

type nopLogger struct{}

// New returns a new logger instance.
func New() logger.Logger {
	return new(nopLogger)
}

// Debug logs the vals at Debug level.
func (*nopLogger) Debug(vals ...interface{}) {}

// Debugf logs the formatted message at Debug level.
func (*nopLogger) Debugf(format string, vals ...interface{}) {}

// Debugd logs the message with details at Debug level.
func (*nopLogger) Debugd(msg string, details ...interface{}) {}

// Info logs the vals at Info level.
func (*nopLogger) Info(vals ...interface{}) {}

// Infof logs the formatted message at Info level.
func (*nopLogger) Infof(format string, vals ...interface{}) {}

// Infod logs the message with details at Info level.
func (*nopLogger) Infod(msg string, details ...interface{}) {}

// Warn logs the vals at Warn level.
func (*nopLogger) Warn(vals ...interface{}) {}

// Warnf logs the formatted message at Warn level.
func (*nopLogger) Warnf(format string, vals ...interface{}) {}

// Warnd logs the message with details at Warn level.
func (*nopLogger) Warnd(msg string, details ...interface{}) {}

// Error logs the vals at Error level.
func (*nopLogger) Error(vals ...interface{}) {}

// Errorf logs the formatted message at Error level.
func (*nopLogger) Errorf(format string, vals ...interface{}) {}

// Errord logs the message with details at Error level.
func (*nopLogger) Errord(msg string, details ...interface{}) {}

// Fatal logs the vals at Fatal level, then calls os.Exit(1).
func (*nopLogger) Fatal(vals ...interface{}) {}

// Fatalf logs the formatted message at Fatal level, then calls os.Exit(1).
func (*nopLogger) Fatalf(format string, vals ...interface{}) {}

// Fatald logs the message with details at Fatal level, then calls os.Exit(1).
func (*nopLogger) Fatald(msg string, details ...interface{}) {}

// Close calls finalizer of logger implementations.
func (*nopLogger) Close() error {
	return nil
}

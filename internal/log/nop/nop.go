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
package nop

import "github.com/vdaas/vald/internal/log/logger"

type nopLogger struct{}

// New returns a new logger instance.
func New() logger.Logger {
	return new(nopLogger)
}

// Debug logs the vals at Debug level.
func (*nopLogger) Debug(...any) {}

// Debugf logs the formatted message at Debug level.
func (*nopLogger) Debugf(string, ...any) {}

// Debugd logs the message with details at Debug level.
func (*nopLogger) Debugd(string, ...any) {}

// Info logs the vals at Info level.
func (*nopLogger) Info(...any) {}

// Infof logs the formatted message at Info level.
func (*nopLogger) Infof(string, ...any) {}

// Infod logs the message with details at Info level.
func (*nopLogger) Infod(string, ...any) {}

// Warn logs the vals at Warn level.
func (*nopLogger) Warn(...any) {}

// Warnf logs the formatted message at Warn level.
func (*nopLogger) Warnf(string, ...any) {}

// Warnd logs the message with details at Warn level.
func (*nopLogger) Warnd(string, ...any) {}

// Error logs the vals at Error level.
func (*nopLogger) Error(...any) {}

// Errorf logs the formatted message at Error level.
func (*nopLogger) Errorf(string, ...any) {}

// Errord logs the message with details at Error level.
func (*nopLogger) Errord(string, ...any) {}

// Fatal logs the vals at Fatal level, then calls os.Exit(1).
func (*nopLogger) Fatal(...any) {}

// Fatalf logs the formatted message at Fatal level, then calls os.Exit(1).
func (*nopLogger) Fatalf(string, ...any) {}

// Fatald logs the message with details at Fatal level, then calls os.Exit(1).
func (*nopLogger) Fatald(string, ...any) {}

// Close calls finalizer of logger implementations.
func (*nopLogger) Close() error {
	return nil
}

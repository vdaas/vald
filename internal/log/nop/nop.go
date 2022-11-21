// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
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
func (*nopLogger) Debug(_ ...interface{}) {}

// Debugf logs the formatted message at Debug level.
func (*nopLogger) Debugf(_ string, _ ...interface{}) {}

// Debugd logs the message with details at Debug level.
func (*nopLogger) Debugd(_ string, _ ...interface{}) {}

// Info logs the vals at Info level.
func (*nopLogger) Info(_ ...interface{}) {}

// Infof logs the formatted message at Info level.
func (*nopLogger) Infof(_ string, _ ...interface{}) {}

// Infod logs the message with details at Info level.
func (*nopLogger) Infod(_ string, _ ...interface{}) {}

// Warn logs the vals at Warn level.
func (*nopLogger) Warn(_ ...interface{}) {}

// Warnf logs the formatted message at Warn level.
func (*nopLogger) Warnf(_ string, _ ...interface{}) {}

// Warnd logs the message with details at Warn level.
func (*nopLogger) Warnd(_ string, _ ...interface{}) {}

// Error logs the vals at Error level.
func (*nopLogger) Error(_ ...interface{}) {}

// Errorf logs the formatted message at Error level.
func (*nopLogger) Errorf(_ string, _ ...interface{}) {}

// Errord logs the message with details at Error level.
func (*nopLogger) Errord(_ string, _ ...interface{}) {}

// Fatal logs the vals at Fatal level, then calls os.Exit(1).
func (*nopLogger) Fatal(_ ...interface{}) {}

// Fatalf logs the formatted message at Fatal level, then calls os.Exit(1).
func (*nopLogger) Fatalf(_ string, _ ...interface{}) {}

// Fatald logs the message with details at Fatal level, then calls os.Exit(1).
func (*nopLogger) Fatald(_ string, _ ...interface{}) {}

// Close calls finalizer of logger implementations.
func (*nopLogger) Close() error {
	return nil
}

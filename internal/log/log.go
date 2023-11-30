//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

package log

import (
	"sync"

	"github.com/vdaas/vald/internal/log/format"
	"github.com/vdaas/vald/internal/log/glg"
	"github.com/vdaas/vald/internal/log/level"
	logger "github.com/vdaas/vald/internal/log/logger"
	"github.com/vdaas/vald/internal/log/nop"
	"github.com/vdaas/vald/internal/log/retry"
	"github.com/vdaas/vald/internal/log/zap"
)

var (
	l    logger.Logger
	once sync.Once
)

func init() {
	l = glg.New(
		glg.WithLevel(level.DEBUG.String()),
		glg.WithFormat(format.RAW.String()),
		glg.WithRetry(
			retry.New(
				retry.WithError(Error),
				retry.WithWarn(Warn),
			),
		),
	)
}

func Init(opts ...Option) {
	once.Do(func() {
		o := new(option)
		for _, opt := range append(defaultOptions, opts...) {
			opt(o)
		}
		l = getLogger(o)
	})
}

func Close() error {
	return l.Close()
}

func getLogger(o *option) logger.Logger {
	switch o.logType {
	case logger.NOP:
		return nop.New()
	case logger.ZAP:
		z, err := zap.New(
			zap.WithLevel(o.level.String()),
			zap.WithFormat(o.format.String()),
		)
		if err == nil {
			return z
		}

		// fallback
		fallthrough
	case logger.GLG:
		fallthrough
	default:
		return glg.New(
			glg.WithLevel(o.level.String()),
			glg.WithFormat(o.format.String()),
			glg.WithRetry(
				retry.New(
					retry.WithError(Error),
					retry.WithWarn(Warn),
				),
			),
		)

	}
}

func Bold(str string) string {
	return "\033[1m" + str + "\033[22m"
}

func Debug(vals ...interface{}) {
	l.Debug(vals...)
}

func Debugf(format string, vals ...interface{}) {
	l.Debugf(format, vals...)
}

func Debugd(msg string, details ...interface{}) {
	l.Debugd(msg, details...)
}

func Info(vals ...interface{}) {
	l.Info(vals...)
}

func Infof(format string, vals ...interface{}) {
	l.Infof(format, vals...)
}

func Infod(msg string, details ...interface{}) {
	l.Infod(msg, details...)
}

func Warn(vals ...interface{}) {
	l.Warn(vals...)
}

func Warnf(format string, vals ...interface{}) {
	l.Warnf(format, vals...)
}

func Warnd(msg string, details ...interface{}) {
	l.Warnd(msg, details...)
}

func Error(vals ...interface{}) {
	l.Error(vals...)
}

func Errorf(format string, vals ...interface{}) {
	l.Errorf(format, vals...)
}

func Errord(msg string, details ...interface{}) {
	l.Errord(msg, details...)
}

func Fatal(vals ...interface{}) {
	l.Fatal(vals...)
}

func Fatalf(format string, vals ...interface{}) {
	l.Fatalf(format, vals...)
}

func Fatald(msg string, details ...interface{}) {
	l.Fatald(msg, details...)
}

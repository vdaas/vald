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

package log

import (
	"sync"

	"github.com/vdaas/vald/internal/log/glg"
	logger "github.com/vdaas/vald/internal/log/logger"
)

type Logger interface {
	Debug(vals ...interface{})
	Debugf(format string, vals ...interface{})
	Info(vals ...interface{})
	Infof(format string, vals ...interface{})
	Warn(vals ...interface{})
	Warnf(format string, vals ...interface{})
	Error(vals ...interface{})
	Errorf(format string, vals ...interface{})
	Fatal(vals ...interface{})
	Fatalf(format string, vals ...interface{})
}

var (
	l    Logger
	once sync.Once
)

func Init(opts ...Option) {
	once.Do(func() {
		o := new(option)
		for _, opt := range append(defaultOptions, opts...) {
			opt(o)
		}
		l = getLogger(o)
	})
}

func getLogger(o *option) Logger {
	switch o.logType {
	case logger.GLG:
		fallthrough
	default:
		gopts := []glg.Option{
			glg.WithLevel(o.level.String()),
			glg.WithFormat(o.format.String()),
		}
		return glg.New(gopts...)
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

func Info(vals ...interface{}) {
	l.Info(vals...)
}

func Infof(format string, vals ...interface{}) {
	l.Infof(format, vals...)
}

func Warn(vals ...interface{}) {
	l.Warn(vals...)
}

func Warnf(format string, vals ...interface{}) {
	l.Warnf(format, vals...)
}

func Error(vals ...interface{}) {
	l.Error(vals...)
}

func Errorf(format string, vals ...interface{}) {
	l.Errorf(format, vals...)
}

func Fatal(vals ...interface{}) {
	l.Fatal(vals...)
}

func Fatalf(format string, vals ...interface{}) {
	l.Fatalf(format, vals...)
}

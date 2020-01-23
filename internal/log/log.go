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
	"reflect"
	"sync"

	"github.com/vdaas/vald/internal/errors"
)

type Logger interface {
	Info(vals ...interface{})
	Infof(format string, vals ...interface{})
	Debug(vals ...interface{})
	Debugf(format string, vals ...interface{})
	Warn(vals ...interface{})
	Warnf(format string, vals ...interface{})
	Error(vals ...interface{})
	Errorf(format string, vals ...interface{})
	Fatal(vals ...interface{})
	Fatalf(format string, vals ...interface{})
}

var (
	logger Logger
	once   sync.Once
)

func Init(opts ...Option) {
	once.Do(func() {
		o := (&option{}).apply(append(defaultOptions, opts...)...)
		logger = o.logger
	})
}

func Bold(str string) string {
	return "\033[1m" + str + "\033[0m"
}

func Info(vals ...interface{}) {
	logger.Info(vals...)
}
func Infof(format string, vals ...interface{}) {
	logger.Infof(format, vals...)
}

func Debug(vals ...interface{}) {
	logger.Debug(vals...)
}

func Debugf(format string, vals ...interface{}) {
	logger.Debugf(format, vals...)
}

func Warn(vals ...interface{}) {
	logger.Warn(vals...)
}

func Warnf(format string, vals ...interface{}) {
	logger.Warnf(format, vals...)
}

func Error(vals ...interface{}) {
	logger.Error(vals...)
}

func Errorf(format string, vals ...interface{}) {
	logger.Errorf(format, vals...)
}

func Fatal(vals ...interface{}) {
	logger.Fatal(vals...)
}

func Fatalf(format string, vals ...interface{}) {
	logger.Fatalf(format, vals...)
}

var (
	retryOut = func(fn func(vals ...interface{}) error, vals ...interface{}) {
		retryOutf(func(format string, vals ...interface{}) error {
			return fn(vals...)
		}, "", vals...)
	}

	retryOutf = func(fn func(format string, vals ...interface{}) error, format string, vals ...interface{}) {
		if err := fn(format, vals...); err != nil {
			rv := reflect.ValueOf(fn)

			Warn(errors.ErrLoggingRetry(err, rv))

			err = fn(format, vals...)
			if err != nil {
				Error(errors.ErrLoggingFaild(err, rv))

				err = fn(format, vals...)
				if err != nil {
					panic(err)
				}
			}
		}
	}
)

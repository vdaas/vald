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
package log

import (
	"github.com/vdaas/vald/internal/log/format"
	"github.com/vdaas/vald/internal/log/glg"
	"github.com/vdaas/vald/internal/log/level"
	logger "github.com/vdaas/vald/internal/log/logger"
	"github.com/vdaas/vald/internal/log/retry"
)

type Option func(*option)

var defaultOptions = []Option{
	WithLevel(level.DEBUG.String()),
	WithLogger(
		glg.New(
			glg.WithRetry(
				retry.New(
					retry.WithError(Error),
					retry.WithWarn(Warn),
				),
			),
		),
	),
}

type option struct {
	logType logger.Type
	level   level.Level
	format  format.Format
	logger  logger.Logger
}

func WithLogger(l logger.Logger) Option {
	return func(o *option) {
		if l == nil {
			return
		}
		o.logger = l
	}
}

func WithLoggerType(str string) Option {
	return func(o *option) {
		if str == "" {
			return
		}
		o.logType = logger.Atot(str)
	}
}

func WithLevel(str string) Option {
	return func(o *option) {
		if str == "" {
			return
		}
		o.level = level.Atol(str)
	}
}

func WithFormat(str string) Option {
	return func(o *option) {
		if str == "" {
			return
		}
		o.format = format.Atof(str)
	}
}

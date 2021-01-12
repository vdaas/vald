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

package glg

import (
	"github.com/kpango/glg"
	"github.com/vdaas/vald/internal/log/format"
	"github.com/vdaas/vald/internal/log/level"
	log "github.com/vdaas/vald/internal/log/logger"
	"github.com/vdaas/vald/internal/log/retry"
)

type logger struct {
	format format.Format
	level  level.Level
	retry  retry.Retry
	glg    *glg.Glg
}

// New returns a new logger instance.
func New(opts ...Option) log.Logger {
	l := new(logger)
	for _, opt := range append(defaultOptions, opts...) {
		opt(l)
	}

	return l.
		setLevelMode(l.level).
		setLogFormat(l.format)
}

func (l *logger) setLevelMode(lv level.Level) *logger {
	l.glg.SetMode(glg.NONE)

	switch lv {
	case level.DEBUG:
		l.glg.SetLevelMode(glg.DEBG, glg.STD)
		fallthrough
	case level.INFO:
		l.glg.SetLevelMode(glg.INFO, glg.STD)
		fallthrough
	case level.WARN:
		l.glg.SetLevelMode(glg.WARN, glg.STD)
		fallthrough
	case level.ERROR:
		l.glg.SetLevelMode(glg.ERR, glg.STD)
		fallthrough
	case level.FATAL:
		l.glg.SetLevelMode(glg.FATAL, glg.STD)
	}

	return l
}

func (l *logger) setLogFormat(fmt format.Format) *logger {
	if fmt == format.JSON {
		l.glg.EnableJSON()
	}
	return l
}

func (l *logger) Info(vals ...interface{}) {
	l.retry.Out(l.glg.Info, vals...)
}

func (l *logger) Infof(format string, vals ...interface{}) {
	l.retry.Outf(l.glg.Infof, format, vals...)
}

func (l *logger) Debug(vals ...interface{}) {
	l.retry.Out(l.glg.Debug, vals...)
}

func (l *logger) Debugf(format string, vals ...interface{}) {
	l.retry.Outf(l.glg.Debugf, format, vals...)
}

func (l *logger) Warn(vals ...interface{}) {
	l.retry.Out(l.glg.Warn, vals...)
}

func (l *logger) Warnf(format string, vals ...interface{}) {
	l.retry.Outf(l.glg.Warnf, format, vals...)
}

func (l *logger) Error(vals ...interface{}) {
	l.retry.Out(l.glg.Error, vals...)
}

func (l *logger) Errorf(format string, vals ...interface{}) {
	l.retry.Outf(l.glg.Errorf, format, vals...)
}

func (l *logger) Fatal(vals ...interface{}) {
	l.glg.Fatal(vals...)
}

func (l *logger) Fatalf(format string, vals ...interface{}) {
	l.glg.Fatalf(format, vals...)
}

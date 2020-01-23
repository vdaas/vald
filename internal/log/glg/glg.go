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

package glg

import (
	logger "github.com/kpango/glg"
	"github.com/vdaas/vald/internal/log/retry"
)

type GlgLogger struct {
	lv  level
	log *logger.Glg
}

// New returns a new GlgLogger instance.
func New(g *logger.Glg, opts ...Option) *GlgLogger {
	gl := (&GlgLogger{
		log: g,
	}).apply(append(defaultOpts, opts...)...)

	gl.setLevelMode(gl.lv)
	return gl
}

func Default() *GlgLogger {
	gl := (&GlgLogger{
		log: logger.Get(),
	}).apply(defaultOpts...)

	gl.setLevelMode(gl.lv)
	return gl
}

func (l *GlgLogger) apply(opts ...Option) *GlgLogger {
	for _, opt := range opts {
		opt(l)
	}
	return l
}

func (l *GlgLogger) setLevelMode(lv level) {
	l.log.SetMode(logger.NONE)

	switch lv {
	case DEBUG:
		l.log.SetLevelMode(logger.DEBG, logger.STD)
		fallthrough
	case INFO:
		l.log.SetLevelMode(logger.INFO, logger.STD)
		fallthrough
	case WARN:
		l.log.SetLevelMode(logger.WARN, logger.STD)
		fallthrough
	case ERROR:
		l.log.SetLevelMode(logger.ERR, logger.STD)
		fallthrough
	case FATAL:
		l.log.SetLevelMode(logger.FAIL, logger.STD)
	}
}

func (l *GlgLogger) Info(vals ...interface{}) {
	retry.Out(l.log.Info, vals...)
}

func (l *GlgLogger) Infof(format string, vals ...interface{}) {
	retry.Outf(l.log.Infof, format, vals...)
}

func (l *GlgLogger) Debug(vals ...interface{}) {
	retry.Out(l.log.Debug, vals...)
}

func (l *GlgLogger) Debugf(format string, vals ...interface{}) {
	retry.Outf(l.log.Debugf, format, vals...)
}

func (l *GlgLogger) Warn(vals ...interface{}) {
	retry.Out(l.log.Warn, vals...)
}

func (l *GlgLogger) Warnf(format string, vals ...interface{}) {
	retry.Outf(l.log.Warnf, format, vals...)
}

func (l *GlgLogger) Error(vals ...interface{}) {
	retry.Out(l.log.Error, vals...)
}

func (l *GlgLogger) Errorf(format string, vals ...interface{}) {
	retry.Outf(l.log.Errorf, format, vals...)
}

func (l *GlgLogger) Fatal(vals ...interface{}) {
	l.log.Fatal(vals...)
}

func (l *GlgLogger) Fatalf(format string, vals ...interface{}) {
	l.log.Fatalf(format, vals...)
}

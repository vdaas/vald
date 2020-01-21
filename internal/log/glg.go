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
	"github.com/kpango/glg"
)

type glglogger struct {
	lv  logLevel
	log *glg.Glg
}

// New returns a new glglogger instance.
func NewGlg(g *glg.Glg, opts ...GlgOption) Logger {
	gl := (&glglogger{
		log: g,
	}).apply(append(defaultGlgOpts, opts...)...)

	gl.setMode(gl.lv)
	return gl
}

func DefaultGlg() Logger {
	gl := (&glglogger{
		log: glg.Get(),
	}).apply(defaultGlgOpts...)

	gl.setMode(gl.lv)
	return gl
}

func (l *glglogger) apply(opts ...GlgOption) *glglogger {
	for _, opt := range opts {
		opt(l)
	}
	return l
}

func (l *glglogger) setMode(lv logLevel) {
	l.log.SetMode(glg.NONE)

	switch lv {
	case DEBUG:
		l.log.SetLevelMode(glg.DEBG, glg.STD)
		fallthrough
	case INFO:
		l.log.SetLevelMode(glg.INFO, glg.STD)
		fallthrough
	case WARN:
		l.log.SetLevelMode(glg.WARN, glg.STD)
		fallthrough
	case ERROR:
		l.log.SetLevelMode(glg.ERR, glg.STD)
		fallthrough
	case FATAL:
		l.log.SetLevelMode(glg.FAIL, glg.STD)
	}
}

func (l *glglogger) Info(vals ...interface{}) {
	l.log.Info(vals...)
}

func (l *glglogger) Infof(format string, vals ...interface{}) {
	l.log.Infof(format, vals...)
}

func (l *glglogger) Debug(vals ...interface{}) {
	l.log.Debug(vals...)
}

func (l *glglogger) Debugf(format string, vals ...interface{}) {
	l.log.Debugf(format, vals...)
}

func (l *glglogger) Warn(vals ...interface{}) {
	l.log.Warn(vals...)
}

func (l *glglogger) Warnf(format string, vals ...interface{}) {
	l.log.Warnf(format, vals...)
}

func (l *glglogger) Error(vals ...interface{}) {
	l.log.Error(vals...)
}

func (l *glglogger) Errorf(format string, vals ...interface{}) {
	l.log.Errorf(format, vals...)
}

func (l *glglogger) Fatal(vals ...interface{}) {
	l.log.Fatal(vals...)
}

func (l *glglogger) Fatalf(format string, vals ...interface{}) {
	l.log.Fatalf(format, vals...)
}

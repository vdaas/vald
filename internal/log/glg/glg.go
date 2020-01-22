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
	kglg "github.com/kpango/glg"
	"github.com/vdaas/vald/internal/log"
)

type glglogger struct {
	lv  log.Level
	log *kglg.Glg
}

// New returns a new glglogger instance.
func New(g *kglg.Glg, opts ...Option) log.Logger {
	gl := (&glglogger{
		log: g,
	}).apply(append(defaultOpts, opts...)...)

	gl.setLevelMode(gl.lv)
	return gl
}

func Default() log.Logger {
	gl := (&glglogger{
		log: kglg.Get(),
	}).apply(defaultOpts...)

	gl.setLevelMode(gl.lv)
	return gl
}

func (l *glglogger) apply(opts ...Option) *glglogger {
	for _, opt := range opts {
		opt(l)
	}
	return l
}

func (l *glglogger) setLevelMode(lv log.Level) {
	l.log.SetMode(kglg.NONE)

	switch lv {
	case log.DEBUG:
		l.log.SetLevelMode(kglg.DEBG, kglg.STD)
		fallthrough
	case log.INFO:
		l.log.SetLevelMode(kglg.INFO, kglg.STD)
		fallthrough
	case log.WARN:
		l.log.SetLevelMode(kglg.WARN, kglg.STD)
		fallthrough
	case log.ERROR:
		l.log.SetLevelMode(kglg.ERR, kglg.STD)
		fallthrough
	case log.FATAL:
		l.log.SetLevelMode(kglg.FAIL, kglg.STD)
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

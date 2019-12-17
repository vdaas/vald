//
// Copyright (C) 2019 Vdaas.org Vald team ( kpango, kmrmt, rinx )
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
	log *glg.Glg
}

// New returns a new glglogger instance.
func NewGlg(g *glg.Glg) Logger {
	return &glglogger{
		log: g,
	}
}

func DefaultGlg() Logger {
	return &glglogger{
		log: glg.Get(),
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

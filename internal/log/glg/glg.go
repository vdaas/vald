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
	"github.com/kpango/glg"
)

type logger struct {

	*glg.Glg
}

// New returns a new logger instance.
func New(opts ...Option) *logger {
	return &logger{g}
}

func (l *logger) Info(vals ...interface{}) {
	l.Info(vals...)
}

func (l *logger) Infof(format string, vals ...interface{}) {
	l.Infof(format, vals...)
}

func (l *logger) Debug(vals ...interface{}) {
	l.Debug(vals...)
}

func (l *logger) Debugf(format string, vals ...interface{}) {
	l.Debugf(format, vals...)
}

func (l *logger) Warn(vals ...interface{}) {
	l.Warn(vals...)
}

func (l *logger) Warnf(format string, vals ...interface{}) {
	l.Warnf(format, vals...)
}

func (l *logger) Error(vals ...interface{}) {
	l.Error(vals...)
}

func (l *logger) Errorf(format string, vals ...interface{}) {
	l.Errorf(format, vals...)
}

func (l *logger) Fatal(vals ...interface{}) {
	l.Fatal(vals...)
}

func (l *logger) Fatalf(format string, vals ...interface{}) {
	l.Fatalf(format, vals...)
}

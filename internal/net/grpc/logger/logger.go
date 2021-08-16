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
package logger

import (
	"os"
	"strconv"
	"sync"

	"github.com/vdaas/vald/internal/log"
	glog "google.golang.org/grpc/grpclog"
)

type logger struct {
	v int
}

var once sync.Once

func Init() {
	once.Do(func() {
		var v int
		if vl, err := strconv.Atoi(os.Getenv("GRPC_GO_LOG_VERBOSITY_LEVEL")); err == nil {
			v = vl
		}
		glog.SetLoggerV2(&logger{v: v})
	})
}

func (l *logger) Info(args ...interface{}) {
	log.Info(args...)
}

func (l *logger) Infoln(args ...interface{}) {
	log.Info(args...)
}

func (l *logger) Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func (l *logger) Warning(args ...interface{}) {
	log.Warn(args...)
}

func (l *logger) Warningln(args ...interface{}) {
	log.Warn(args...)
}

func (l *logger) Warningf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

func (l *logger) Error(args ...interface{}) {
	log.Error(args...)
}

func (l *logger) Errorln(args ...interface{}) {
	log.Error(args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func (l *logger) Fatal(args ...interface{}) {
	log.Fatal(args...)
}

func (l *logger) Fatalln(args ...interface{}) {
	log.Fatal(args...)
}

func (l *logger) Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

func (l *logger) V(v int) bool {
	return v <= l.v
}

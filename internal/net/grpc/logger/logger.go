// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
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

const (
	tag = "[gRPC Log]"
)

// Init initialize the logging level and the logger.
func Init() {
	once.Do(func() {
		var v int
		if vl, err := strconv.Atoi(os.Getenv("GRPC_GO_LOG_VERBOSITY_LEVEL")); err == nil {
			v = vl
		}
		glog.SetLoggerV2(&logger{v: v})
	})
}

// Info prints the debug log to the logger.
func (*logger) Info(args ...interface{}) {
	log.Debug(append([]interface{}{tag}, args...)...)
}

// Infoln prints the debug log to the logger.
func (*logger) Infoln(args ...interface{}) {
	log.Debug(append([]interface{}{tag}, args...)...)
}

// Infof prints the debug log to the logger.
func (*logger) Infof(format string, args ...interface{}) {
	log.Debugf(tag+"\t"+format, args...)
}

// Warning prints the warning log to the logger.
func (*logger) Warning(args ...interface{}) {
	log.Warn(append([]interface{}{tag}, args...)...)
}

// Warningln prints the warning log to the logger.
func (*logger) Warningln(args ...interface{}) {
	log.Warn(append([]interface{}{tag}, args...)...)
}

// Warningf prints the warning log to the logger.
func (*logger) Warningf(format string, args ...interface{}) {
	log.Warnf(tag+"\t"+format, args...)
}

// Error prints the error log to the logger.
func (*logger) Error(args ...interface{}) {
	log.Error(append([]interface{}{tag}, args...)...)
}

// Errorln prints the error log to the logger.
func (*logger) Errorln(args ...interface{}) {
	log.Error(append([]interface{}{tag}, args...)...)
}

// Errorf prints the error log to the logger.
func (*logger) Errorf(format string, args ...interface{}) {
	log.Errorf(tag+"\t"+format, args...)
}

// Fatal prints the fatal log to the logger and exit the program.
func (*logger) Fatal(args ...interface{}) {
	// skipcq: RVV-A0003
	log.Fatal(append([]interface{}{tag}, args...)...)
}

// Fatalln prints the fatal log to the logger and exit the program.
func (*logger) Fatalln(args ...interface{}) {
	// skipcq: RVV-A0003
	log.Fatal(append([]interface{}{tag}, args...)...)
}

// Fatalf prints the fatal log to the logger and exit the program.
func (*logger) Fatalf(format string, args ...interface{}) {
	// skipcq: RVV-A0003
	log.Fatalf(tag+"\t"+format, args...)
}

// V returns if the v is less than the verbosity level.
func (l *logger) V(v int) bool {
	return v <= l.v
}

//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

// Package k8s provides kubernetes control functionality
package k8s

import (
	"fmt"

	"github.com/go-logr/logr"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
)

var setLoggerOnce sync.Once

// setControllerRuntimeLogger bridges controller-runtime's logr output into
// the Vald logger so that manager/controller runtime information (leader
// election details, reconcile worker count, cache sync) appears in the
// component logs instead of being silently discarded.
func setControllerRuntimeLogger() {
	setLoggerOnce.Do(func() {
		ctrllog.SetLogger(logr.New(new(logSink)))
	})
}

// logSink adapts the logr.LogSink interface onto internal/log. V-levels above
// zero are treated as debug output; everything else is informational.
type logSink struct {
	name string
	kvs  []any
}

func (*logSink) Init(logr.RuntimeInfo) {}

func (*logSink) Enabled(int) bool { return true }

func (s *logSink) Info(level int, msg string, kvs ...any) {
	if level > 0 {
		log.Debugf("[controller-runtime] %s%s", s.prefix()+msg, render(s.kvs, kvs))
		return
	}
	log.Infof("[controller-runtime] %s%s", s.prefix()+msg, render(s.kvs, kvs))
}

func (s *logSink) Error(err error, msg string, kvs ...any) {
	log.Errorf("[controller-runtime] %s: %v%s", s.prefix()+msg, err, render(s.kvs, kvs))
}

func (s *logSink) WithValues(kvs ...any) logr.LogSink {
	ns := *s
	ns.kvs = append(append(make([]any, 0, len(s.kvs)+len(kvs)), s.kvs...), kvs...)
	return &ns
}

func (s *logSink) WithName(name string) logr.LogSink {
	ns := *s
	if ns.name != "" {
		ns.name += "."
	}
	ns.name += name
	return &ns
}

func (s *logSink) prefix() string {
	if s.name == "" {
		return ""
	}
	return s.name + ": "
}

func render(base, kvs []any) string {
	if len(base) == 0 && len(kvs) == 0 {
		return ""
	}
	var b strings.Builder
	for i := 0; i+1 < len(base); i += 2 {
		fmt.Fprintf(&b, "\t%v=%v", base[i], base[i+1])
	}
	for i := 0; i+1 < len(kvs); i += 2 {
		fmt.Fprintf(&b, "\t%v=%v", kvs[i], kvs[i+1])
	}
	return b.String()
}

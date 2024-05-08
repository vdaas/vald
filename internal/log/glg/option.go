//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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

package glg

import (
	"github.com/kpango/glg"
	"github.com/vdaas/vald/internal/log/format"
	"github.com/vdaas/vald/internal/log/level"
	"github.com/vdaas/vald/internal/log/retry"
)

type Option func(l *logger)

var defaultOptions = []Option{
	WithGlg(glg.Get()),
	WithLevel(level.DEBUG.String()),
	WithRetry(retry.New()),
}

func WithGlg(g *glg.Glg) Option {
	return func(l *logger) {
		if g == nil {
			return
		}
		l.glg = g
	}
}

func WithFormat(str string) Option {
	return func(l *logger) {
		if str == "" {
			return
		}
		l.format = format.Atof(str)
	}
}

func WithLevel(str string) Option {
	return func(l *logger) {
		if str == "" {
			return
		}
		l.level = level.Atol(str)
	}
}

func WithRetry(rt retry.Retry) Option {
	return func(l *logger) {
		if rt == nil {
			return
		}
		l.retry = rt
	}
}

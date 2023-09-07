// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package info

import (
	"runtime"

	"github.com/vdaas/vald/internal/errors"
)

// Option represent the functional option for info.
type Option func(i *info) error

var defaultOpts = []Option{
	WithRuntimeCaller(runtime.Caller),
	WithRuntimeFuncForPC(runtime.FuncForPC),
}

// WithServerName returns the option to set the server name.
func WithServerName(s string) Option {
	return func(i *info) error {
		if len(s) == 0 {
			return errors.NewErrInvalidOption("ServerName", s)
		}
		i.detail.ServerName = s
		return nil
	}
}

// WithRuntimeCaller returns the option to set the runtime Caller function.
func WithRuntimeCaller(f func(skip int) (pc uintptr, file string, line int, ok bool)) Option {
	return func(i *info) error {
		if f == nil {
			return errors.NewErrInvalidOption("RuntimeCaller", f)
		}
		i.rtCaller = f
		return nil
	}
}

// WithRuntimeFuncForPC returns the option to set the runtime FuncForPC function.
func WithRuntimeFuncForPC(f func(pc uintptr) *runtime.Func) Option {
	return func(i *info) error {
		if f == nil {
			return errors.NewErrInvalidOption("RuntimeFuncForPC", f)
		}
		i.rtFuncForPC = f
		return nil
	}
}

//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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

// Package safety provides safety functionality like revcover
package safety

import (
	"runtime"
	"runtime/debug"

	"github.com/vdaas/vald/internal/conv"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/log"
)

func RecoverFunc(fn func() error) func() error {
	return recoverFn(fn, true)
}

func RecoverWithoutPanicFunc(fn func() error) func() error {
	return recoverFn(fn, false)
}

func recoverFn(fn func() error, withPanic bool) func() error {
	return func() (err error) {
		if fn == nil {
			return nil
		}
		defer func() {
			if r := recover(); r != nil {
				ds := debug.Stack()
				infoStr := info.Get().String() + "\r\n" + conv.Btoa(ds)
				log.Warnf("function %#v panic recovered: %#v\ninfo:\n%s", fn, r, infoStr)
				switch x := r.(type) {
				case runtime.Error:
					err = errors.ErrRuntimeError(err, x)
					if withPanic {
						log.Errorf("recovered but this thread is going to panic: the reason is runtimer.Error\nerror:\t%v\ninfo:\n%s\nrecovered:\t%#v", err, infoStr, r)
						panic(err)
					}
				case *string:
					err = errors.ErrPanicString(err, *x)
				case string:
					err = errors.ErrPanicString(err, x)
				case error:
					err = errors.Join(err, x)
				default:
					err = errors.ErrPanicRecovered(err, x)
				}
				if err != nil {
					log.Errorf("function %#v panic recovered error:\t%v\ninfo:\n%s\nrecovered:\t%#v", fn, err, infoStr, r)
				}
			}
		}()
		err = fn()
		return err
	}
}

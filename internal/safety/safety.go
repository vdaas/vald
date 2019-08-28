// MIT License
//
// Copyright (c) 2019 kpango (Yusuke Kato)
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Package safety provides safety functionallity like revcover
package safety

import (
	"runtime"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
)

func RecoverFunc(fn func() error) func() error {
	return func() (err error) {
		defer func() {
			err = RecoverWithError(err)
		}()
		err = fn()
		return err
	}
}

func RecoverWithError(err error) error {
	if r := recover(); r != nil {
		log.Warnf("recovered in f %v", r)
		switch x := r.(type) {
		case runtime.Error:
			err = errors.ErrRuntimeError(err, x)
			panic(err)
		case string:
			err = errors.ErrPanicString(err, x)
		case error:
			err = errors.Wrap(err, x.Error())
		default:
			err = errors.ErrPanicRecovered(err, x)
		}
		log.Error(err)
	}
	return err
}

//
// Copyright (C) 2019 Vdaas.org Vald team ( kpango, kou-m, rinx )
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Package errors provides error types and function
package errors

import (
	"reflect"
	"runtime"

	// "github.com/pkg/errors"
	"github.com/cockroachdb/errors"
	"github.com/cockroachdb/errors/errbase"
)

var (
	ErrTimeoutParseFailed = func(timeout string) error {
		return Errorf("invalid timeout value: %s\t:timeout parse error out put failed", timeout)
	}

	ErrServerNotFound = func(name string) error {
		return Errorf("server %s not found", name)
	}

	ErrOptionFailed = func(err error, ref reflect.Value) error {
		return Wrapf(err, "failed to setup option :\t%s",
			runtime.FuncForPC(ref.Pointer()).Name())
	}

	ErrArgumentParseFailed = func(err error) error {
		return Wrap(err, "argument parse failed")
	}

	ErrBackoffTimeout = func(err error) error {
		return Wrap(err, "backoff timeout by limitation")
	}

	ErrInvalidTypeConversion = func(i interface{}, tgt interface{}) error {
		return Errorf("invalid type conversion %v to %v", reflect.TypeOf(i), reflect.TypeOf(tgt))
	}

	New = func(msg string) error {
		if msg == "" {
			return nil
		}
		return errors.New(msg)
	}

	Wrap = func(err error, msg string) error {
		if err != nil {
			if msg != "" {
				return errors.Wrap(err, msg)
			}
			return err
		}
		return New(msg)
	}

	Wrapf = func(err error, format string, args ...interface{}) error {
		if err != nil {
			if format != "" && len(args) > 0 {
				return errors.Wrapf(err, format, args...)
			}
			return err
		}
		return Errorf(format, args...)
	}

	Cause = func(err error) error {
		if err != nil {
			return errors.Cause(err)
		}
		return nil
	}

	Errorf = func(format string, args ...interface{}) error {
		if format != "" && args != nil && len(args) > 0 {
			return errors.Errorf(format, args...)
		}
		return nil
	}

	As         = errors.As
	Is         = errors.Is
	UnWrapOnce = errbase.UnwrapOnce
	UnWrapAll  = errbase.UnwrapAll
)

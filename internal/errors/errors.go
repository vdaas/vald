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

// Package errors provides error types and function
package errors

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
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

	ErrLoggingRetry = func(err error, ref reflect.Value) error {
		return Wrapf(err, "failed to output %s logs, retrying...",
			runtime.FuncForPC(ref.Pointer()).Name())
	}

	ErrLoggingFailed = func(err error, ref reflect.Value) error {
		return Wrapf(err, "failed to output %s logs",
			runtime.FuncForPC(ref.Pointer()).Name())
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
				return fmt.Errorf("%s: %w", msg, err)
			}
			return err
		}
		return New(msg)
	}

	Wrapf = func(err error, format string, args ...interface{}) error {
		if err != nil {
			if format != "" && len(args) != 0 {
				return Wrap(err, fmt.Sprintf(format, args...))
			}
			return err
		}
		return Errorf(format, args...)
	}

	Cause = func(err error) error {
		if err != nil {
			return errors.Unwrap(err)
		}
		return nil
	}

	Unwrap = errors.Unwrap

	Errorf = func(format string, args ...interface{}) error {
		const delim = " "
		if format == "" && len(args) == 0 {
			return nil
		}
		if len(args) != 0 {
			if format == "" {
				for range args {
					format += "%v" + delim
				}
				format = strings.TrimSuffix(format, delim)
			}
			return fmt.Errorf(format, args...)
		}
		return New(format)
	}

	Is = func(err, target error) bool {
		if target == nil {
			return err == target
		}

		isComparable := reflect.TypeOf(target).Comparable()
		for {
			if isComparable && (err == target ||
				err.Error() == target.Error()) {
				return true
			}
			if x, ok := err.(interface {
				Is(error) bool
			}); ok && x.Is(target) {
				return true
			}
			if uerr := Unwrap(err); uerr == nil {
				return err.Error() == target.Error()
			} else {
				err = uerr
			}
		}
	}

	As = errors.As
)

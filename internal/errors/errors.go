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
	// ErrTimeoutParseFailed represents a function to generate an error that the timeout value is invalid.
	ErrTimeoutParseFailed = func(timeout string) error {
		return Errorf("invalid timeout value: %s\t:timeout parse error out put failed", timeout)
	}

	// ErrServerNotFound represents a funtion to generate an error that the server is not found.
	ErrServerNotFound = func(name string) error {
		return Errorf("server %s not found", name)
	}

	// ErrOptionFailed represents a function to generate an error that the option setup is failed.
	// When ref is zero Value, it will return error with ref is invalid.
	ErrOptionFailed = func(err error, ref reflect.Value) error {
		var str string
		if ref.IsValid() {
			str = runtime.FuncForPC(ref.Pointer()).Name()
		}
		return Wrapf(err, "failed to setup option :\t%s", str)
	}

	// ErrArgumentParseFailed represents a function to generate an error that argument parse is failed.
	ErrArgumentParseFailed = func(err error) error {
		return Wrap(err, "argument parse failed")
	}

	// ErrBackoffTimeout represents a function to generate an error that backoff is timeout by limitation.
	ErrBackoffTimeout = func(err error) error {
		return Wrap(err, "backoff timeout by limitation")
	}

	// Err represents a function to generate an error that type conversion fails due to an invalid input type.
	ErrInvalidTypeConversion = func(i interface{}, tgt interface{}) error {
		return Errorf("invalid type conversion %v to %v", reflect.TypeOf(i), reflect.TypeOf(tgt))
	}

	// ErrLoggingRetry represents a function to generate an error that failing to output logs and retrying to output.
	ErrLoggingRetry = func(err error, ref reflect.Value) error {
		var str string
		if ref.IsValid() {
			str = runtime.FuncForPC(ref.Pointer()).Name()
		}
		return Wrapf(err, "failed to output %s logs, retrying...", str)
	}

	// ErrLoggingFailed represents a function to generate an error that failing to output logs.
	ErrLoggingFailed = func(err error, ref reflect.Value) error {
		var str string
		if ref.IsValid() {
			str = runtime.FuncForPC(ref.Pointer()).Name()
		}
		return Wrapf(err, "failed to output %s logs", str)
	}

	// New represents a function to generate the new error with message.
	// When message is nil, it will return nil instead of error.
	New = func(msg string) error {
		if msg == "" {
			return nil
		}
		return errors.New(msg)
	}

	// Wrap represents a function to generate an error which is used by input error and message.
	// When both of input is nil, it will return new error with message even message is nil.
	// When input error is not nil, it will return error based on input error.
	Wrap = func(err error, msg string) error {
		if err != nil {
			if msg != "" {
				return fmt.Errorf("%s: %w", msg, err)
			}
			return err
		}
		return New(msg)
	}

	// Wrapf represents a function to generate an error which is used by input error, format and args.
	// When all of input is nil, it will return new error based on format and args even these are nil.
	// When input error is not nil, it will return error based on input error.
	Wrapf = func(err error, format string, args ...interface{}) error {
		if err != nil {
			if format != "" && len(args) != 0 {
				return Wrap(err, fmt.Sprintf(format, args...))
			}
			return err
		}
		return Errorf(format, args...)
	}

	// Cause represents a function to generate error when inpurt error is not nil.
	// When input is nil, it will return nil.
	Cause = func(err error) error {
		if err != nil {
			return errors.Unwrap(err)
		}
		return nil
	}

	// Unwrap represents errors.Unwrap
	Unwrap = errors.Unwrap

	// Errorf represents a function to generate an error that based on format and args.
	// When format and args do not satisfy the condition, it will return nil.
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

	// Is represents a function to check whether err and target is same or not.
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

	// As represents errors.As
	As = errors.As
)

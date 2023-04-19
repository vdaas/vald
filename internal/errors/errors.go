//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
	"unsafe"
)

var (
	// ErrTimeoutParseFailed represents a function to generate an error that the timeout value is invalid.
	ErrTimeoutParseFailed = func(timeout string) error {
		return Errorf("invalid timeout value: %s\t:timeout parse error out put failed", timeout)
	}

	// ErrServerNotFound represents a function to generate an error that the server not found.
	ErrServerNotFound = func(name string) error {
		return Errorf("server %s not found", name)
	}

	// ErrOptionFailed represents a function to generate an error that the option setup failed.
	// When ref is zero value, it will return an error with ref is invalid.
	ErrOptionFailed = func(err error, ref reflect.Value) error {
		var str string
		if ref.IsValid() {
			str = runtime.FuncForPC(ref.Pointer()).Name()
		}
		return Wrapf(err, "failed to setup option :\t%s", str)
	}

	// ErrArgumentParseFailed represents a function to generate an error that argument parse failed.
	ErrArgumentParseFailed = func(err error) error {
		return Wrap(err, "argument parse failed")
	}

	// ErrBackoffTimeout represents a function to generate an error that backoff is timeout by limitation.
	ErrBackoffTimeout = func(err error) error {
		return Wrap(err, "backoff timeout by limitation")
	}

	// ErrInvalidTypeConversion represents a function to generate an error that type conversion fails due to an invalid input type.
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

	// New represents a function to generate the new error with a message.
	// When the message is nil, it will return nil instead of an error.
	New = func(msg string) error {
		if msg == "" {
			return nil
		}
		return errors.New(msg)
	}

	// Wrap represents a function to generate an error that is used by input error and message.
	// When both of the input are nil, it will return an error when the error message is not empty.
	// When the input error is not nil, it will return the error based on the input error.
	Wrap = func(err error, msg string) error {
		if err != nil {
			if msg != "" && err.Error() != msg {
				return fmt.Errorf("%s: %w", msg, err)
			}
			return err
		}
		return New(msg)
	}

	// Wrapf represents a function to generate an error that is used by input error, format, and args.
	// When all of the input is nil, it will return a new error based on format and args even these are nil.
	// When the input error is not nil, it will return an error based on the input error.
	Wrapf = func(err error, format string, args ...interface{}) error {
		if err != nil {
			if format != "" && len(args) != 0 {
				return Wrap(err, fmt.Sprintf(format, args...))
			}
			return err
		}
		return Errorf(format, args...)
	}

	// Cause represents a function to generate an error when the input error is not nil.
	// When the input is nil, it will return nil.
	Cause = func(err error) error {
		if err != nil {
			return Unwrap(err)
		}
		return nil
	}

	// Errorf represents a function to generate an error based on format and args.
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

	// As represents errors.As.
	As = errors.As

	// errExpectedErrIsNil represents a function to generate an error that given name's error object is nil.
	errExpectedErrIsNil = func(n string) error {
		return Errorf("expected err is nil: %s", n)
	}
)

// Is represents a function to check whether err and the target is the same or not.
func Is(err, target error) bool {
	if target == nil || err == nil {
		return err == target
	}
	isComparable := reflect.TypeOf(target).Comparable()
	for {
		if isComparable && (err == target ||
			err.Error() == target.Error() ||
			strings.EqualFold(err.Error(), target.Error())) {
			return true
		}

		if x, ok := err.(interface {
			Is(error) bool
		}); ok && x.Is(target) {
			return true
		}
		switch x := err.(type) {
		case interface{ Unwrap() error }:
			err = x.Unwrap()
			if err == nil {
				return isComparable && err == target ||
					err.Error() == target.Error() ||
					strings.EqualFold(err.Error(), target.Error())
			}
		case interface{ Unwrap() []error }:
			for _, err = range x.Unwrap() {
				if Is(err, target) {
					return true
				}
			}
			return isComparable && err == target ||
				err.Error() == target.Error() ||
				strings.EqualFold(err.Error(), target.Error())
		default:
			return isComparable && err == target ||
				err.Error() == target.Error() ||
				strings.EqualFold(err.Error(), target.Error())
		}
	}
}

// Unwrap represents errors.Unwrap.
func Unwrap(err error) error {
	if err == nil {
		return nil
	}
	switch x := err.(type) {
	case interface{ Unwrap() error }:
		return x.Unwrap()
	case interface{ Unwrap() []error }:
		errs := x.Unwrap()
		switch l := len(errs); l {
		case 0:
			return nil
		case 1, 2:
			return errs[0]
		default:
			return Join(errs[:l-1]...)
		}
	default:
		return nil
	}
}

// Join represents a function to generate multiple error when the input errors is not nil.
func Join(errs ...error) error {
	l := len(errs)
	switch l {
	case 0:
		return nil
	case 1:
		if errs[0] != nil {
			return errs[0]
		}
	case 2:
		switch {
		case errs[0] != nil && errs[1] != nil && !Is(errs[0], errs[1]):
			var es []error
			switch x := errs[1].(type) {
			case *joinError:
				es = x.errs
			case interface{ Unwrap() []error }:
				es = x.Unwrap()
			default:
				es = []error{errs[1]}
			}
			switch x := errs[0].(type) {
			case *joinError:
				x.errs = append(x.errs, es...)
				return x
			case interface{ Unwrap() []error }:
				return &joinError{errs: append(x.Unwrap(), es...)}
			default:
				return &joinError{errs: []error{errs[0], errs[1]}}
			}
		case errs[0] != nil:
			return errs[0]
		case errs[1] != nil:
			return errs[1]
		}
	}
	n := 0
	for _, err := range errs {
		if err != nil {
			n++
		}
	}
	if n == 0 {
		return nil
	}
	e := &joinError{
		errs: make([]error, 0, n),
	}
	for _, err := range errs {
		if err != nil {
			e.errs = append(e.errs, err)
		}
	}
	return e
}

type joinError struct {
	errs []error
}

func (e *joinError) Error() string {
	switch len(e.errs) {
	case 0:
		return ""
	case 1:
		return e.errs[0].Error()
	}
	b := make([]byte, 0, len(e.errs)*16)
	for i, err := range e.errs {
		if i > 0 {
			b = append(b, '\n')
		}
		b = append(b, err.Error()...)
	}
	if len(b) == 0 {
		return ""
	}
	// skipcq: GSC-G103
	return unsafe.String(&b[0], len(b))
}

func (e *joinError) Unwrap() []error {
	return e.errs
}

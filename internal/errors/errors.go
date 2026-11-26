// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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

package errors

import (
	"cmp"
	"errors"
	"fmt"
	"maps"
	"reflect"
	"runtime"
	"slices"

	"github.com/vdaas/vald/internal/strings"
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

	// ErrUnimplemented represents a function to generate an error that feature is unimplemented.
	ErrUnimplemented = func(name string) error {
		return Errorf("%s is unimplemented", name)
	}
)

// errExpectedErrIsNil generates an error indicating the given name's error object is nil.
func errExpectedErrIsNil(n string) error {
	return Errorf("expected err is nil: %s", n)
}

// New generates a new error with a message, or nil if msg is empty.
func New(msg string) error {
	if msg == "" {
		return nil
	}
	return errors.New(msg)
}

type wrapError struct {
	err error
	msg string
}

func (e *wrapError) Error() string {
	return e.msg + ": " + e.err.Error()
}

func (e *wrapError) Unwrap() error {
	return e.err
}

// Wrap wraps err with msg unless msg is empty or err already contains msg.
// When err is nil it returns New(msg).
func Wrap(err error, msg string) error {
	if err != nil {
		if msg != "" && !hasMessage(err, msg) {
			return &wrapError{msg: msg, err: err}
		}
		return err
	}
	return New(msg)
}

func hasMessage(err error, msg string) bool {
	for {
		emsg := err.Error()
		if emsg == msg || strings.EqualFold(emsg, msg) {
			return true
		}
		if x, ok := err.(interface{ Is(error) bool }); ok {
			if x.Is(errors.New(msg)) {
				return true
			}
		}
		{
			var x interface{ Unwrap() error }
			var x1 interface{ Unwrap() []error }
			switch {
			case errors.As(err, &x):
				err = x.Unwrap()
				if err == nil {
					return false
				}
			case errors.As(err, &x1):
				for _, e := range x1.Unwrap() {
					if e != nil && hasMessage(e, msg) {
						return true
					}
				}
				return false
			default:
				return false
			}
		}
	}
}

// Wrapf wraps err with the formatted message, returning New/Errorf when err is nil.
func Wrapf(err error, format string, args ...any) error {
	if err != nil {
		if format == "" {
			return err
		}
		if len(args) == 0 {
			return Wrap(err, format)
		}
		return Wrap(err, fmt.Sprintf(format, args...))
	}
	return Errorf(format, args...)
}

// Cause unwraps err once when non-nil.
func Cause(err error) error {
	if err != nil {
		return Unwrap(err)
	}
	return nil
}

// Errorf formats an error, returning nil when both format and args are empty.
func Errorf(format string, args ...any) error {
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

// As delegates to stdlib errors.As.
func As(err error, target any) bool {
	return errors.As(err, target)
}

// Is represents a function to check whether err and the target is the same or not.
func Is(err, target error) (same bool) {
	if target == nil || err == nil {
		return err == target
	}
	return is(err, target)
}

// IsAny reports whether err is non-nil and matches at least one non-nil target.
func IsAny(err error, targets ...error) (same bool) {
	if err == nil || len(targets) == 0 {
		return false
	}
	for _, target := range targets {
		if target != nil && is(err, target) {
			return true
		}
	}
	return false
}

// IsNot reports whether err is an actionable error that matches none of the
// targets: it returns true only when err is non-nil and Is(err, target) is
// false for every non-nil target. It is NOT the negation of Is — for a nil
// err it always returns false (nil is not an error to act on), so callers can
// use it directly as a "should this error be reported" filter. Conversely, a
// non-nil err with an empty target list returns true.
func IsNot(err error, targets ...error) (same bool) {
	if err == nil || len(targets) == 0 {
		return err != nil
	}
	for _, target := range targets {
		if target != nil && is(err, target) {
			return false
		}
	}
	return true
}

func is(err, target error) (same bool) {
	targetComparable := reflect.TypeOf(target).Comparable()
	// target.Error() is computed lazily (only once identity and custom Is checks
	// have failed) and reused across layers and joined children.
	var (
		tmsg   string
		hasMsg bool
	)
	return isRec(err, target, targetComparable, &tmsg, &hasMsg)
}

func isRec(err, target error, targetComparable bool, tmsg *string, hasMsg *bool) bool {
	for {
		//nolint:errorlint // isRec is the Is primitive itself; direct == identity check on comparable target is intentional.
		if targetComparable && err == target {
			return true
		}
		if x, ok := err.(interface {
			Is(error) bool
		}); ok && x.Is(target) {
			return true
		}
		emsg := err.Error()
		if !*hasMsg {
			*tmsg, *hasMsg = target.Error(), true
		}
		if emsg == *tmsg || strings.EqualFold(emsg, *tmsg) {
			return true
		}
		switch x := err.(type) {
		case interface{ Unwrap() error }:
			err = x.Unwrap()
			if err == nil {
				return false
			}
		case interface{ Unwrap() []error }:
			for _, e := range x.Unwrap() {
				if e != nil && isRec(e, target, targetComparable, tmsg, hasMsg) {
					return true
				}
			}
			return false
		default:
			return false
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
		return errs[0]
	case 2: //nolint:mnd // fast path for 2 errors without allocating a map
		switch {
		case errs[0] != nil && errs[1] != nil:
			// is() already short-circuits on identity, but only when the type
			// is comparable; a raw == here panics for uncomparable error types.
			if is(errs[0], errs[1]) {
				return errs[0]
			}
			var es []error
			switch x := errs[1].(type) {
			case *joinError:
				es = x.errs
			case interface{ Unwrap() []error }:
				es = x.Unwrap()
			}
			switch x := errs[0].(type) {
			case *joinError:
				// Copy x.errs into a fresh slice instead of mutating and
				// returning the existing *joinError: error values are expected
				// to be immutable once returned, and a retained reference to
				// errs[0] must not gain the newly joined errors.
				if es != nil {
					return &joinError{errs: RemoveDuplicates(append(append(make([]error, 0, len(x.errs)+len(es)), x.errs...), es...))}
				}
				return &joinError{errs: RemoveDuplicates(append(append(make([]error, 0, len(x.errs)+1), x.errs...), errs[1]))}
			case interface{ Unwrap() []error }:
				u := x.Unwrap()
				if es != nil {
					return &joinError{errs: RemoveDuplicates(append(append(make([]error, 0, len(u)+len(es)), u...), es...))}
				}
				return &joinError{errs: RemoveDuplicates(append(append(make([]error, 0, len(u)+1), u...), errs[1]))}
			default:
				if es != nil {
					list := make([]error, 0, 1+len(es))
					list = append(list, errs[0])
					return &joinError{errs: RemoveDuplicates(append(list, es...))}
				}
				return &joinError{errs: []error{errs[0], errs[1]}}
			}
		case errs[0] != nil:
			return errs[0]
		case errs[1] != nil:
			return errs[1]
		}
		return nil
	}
	var e *joinError
	switch x := errs[0].(type) {
	case *joinError:
		// Copy x.errs into a fresh slice rather than aliasing errs[0]: the
		// append + in-place RemoveDuplicates below would otherwise mutate the
		// caller's *joinError.
		if x != nil && len(x.errs) != 0 {
			e = &joinError{errs: append(make([]error, 0, len(x.errs)+l), x.errs...)}
		}
		errs = errs[1:]
	case interface{ Unwrap() []error }:
		if x != nil {
			if u := x.Unwrap(); len(u) != 0 {
				e = &joinError{errs: append(make([]error, 0, len(u)+l), u...)}
			}
		}
		errs = errs[1:]
	}
	if e == nil {
		e = &joinError{
			errs: make([]error, 0, l),
		}
	}
	e.errs = RemoveDuplicates(append(e.errs, errs...))
	switch len(e.errs) {
	case 0:
		return nil
	case 1:
		// A single remaining error is returned as is, like the 1- and 2-argument cases.
		return e.errs[0]
	}
	return e
}

// JoinMap joins all non-nil errors in errs ordered by their sorted map keys so
// the leading error (whose status code a joined error reports) is deterministic.
func JoinMap[K cmp.Ordered](errs map[K]error) error {
	switch len(errs) {
	case 0:
		return nil
	case 1:
		for _, err := range errs {
			return err
		}
	}
	list := make([]error, 0, len(errs))
	for _, k := range slices.Sorted(maps.Keys(errs)) {
		list = append(list, errs[k])
	}
	return Join(list...)
}

func RemoveDuplicates(errs []error) []error {
	switch len(errs) {
	case 0:
		return nil
	case 1:
		return errs
	case 2: //nolint:mnd // fast path for 2 errors without allocating a map
		if errs[0] == nil {
			if errs[1] == nil {
				clear(errs)
				return errs[:0]
			}
			errs[0] = errs[1]
			errs[1] = nil
			return errs[:1]
		}
		if errs[1] == nil || errs[0].Error() == errs[1].Error() {
			errs[1] = nil
			return errs[:1]
		}
		return errs
	}
	// Two and more elements are deduplicated uniformly by their message: the
	// former length==2 special case used semantic Is() (unwrap + case-insensitive
	// fold) while length>=3 used an exact-string map, so whether two errors were
	// treated as duplicates depended on the argument count rather than the errors
	// themselves. The exact-string map below is the consistent rule.
	seen := make(map[string]struct{}, len(errs))
	var idx uint64
	for _, err := range errs {
		if err != nil {
			key := err.Error()
			if _, ok := seen[key]; !ok {
				seen[key] = struct{}{}
				errs[idx] = err
				idx++
			}
		}
	}
	clear(errs[idx:])
	return errs[:idx]
}

type joinError struct {
	errs []error
}

func (e *joinError) Error() (str string) {
	switch len(e.errs) {
	case 0:
		return ""
	case 1:
		return e.errs[0].Error()
	}
	// A pooled Builder cannot keep its buffer (Reset drops it and String aliases
	// it), so a local Builder is as cheap and simpler.
	var sb strings.Builder
	sb.Grow(len(e.errs) * 16)
	for i, err := range e.errs {
		if i > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(err.Error())
	}
	if sb.Len() == 0 {
		return ""
	}
	return sb.String()
}

func (e *joinError) String() (str string) {
	return e.Error()
}

func (e *joinError) Unwrap() []error {
	return e.errs
}

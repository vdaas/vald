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

// combineOptionErrors merges the non-nil errors in errs into one error.
// The message is identical to chaining Wrap(err, prev.Error()) over errs, but
// every merged error stays reachable through Is/As via the multi-Unwrap chain
// instead of only the last one. It returns nil when errs has no non-nil error.
func combineOptionErrors(errs []error) (e error) {
	for _, err := range errs {
		switch {
		case err == nil:
		case e == nil:
			e = err
		case e.Error() != "" && e.Error() != err.Error() && !Is(err, New(e.Error())):
			e = Errorf("%w: %w", e, err)
		default:
			// e is already represented in err (same behavior as Wrap).
			e = err
		}
	}
	return e
}

// ErrInvalidOption represents the invalid option error.
type ErrInvalidOption struct {
	err    error
	origin error
}

// NewErrInvalidOption represents a function to generate a new error of ErrInvalidOption that invalid option.
func NewErrInvalidOption(name string, val any, errs ...error) error {
	if len(errs) == 0 {
		return &ErrInvalidOption{
			err: Errorf("invalid option, name: %s, val: %v", name, val),
		}
	}
	e := combineOptionErrors(errs)

	return &ErrInvalidOption{
		err:    Wrapf(e, "invalid option, name: %s, val: %v", name, val),
		origin: e,
	}
}

// Error returns a string of ErrInvalidOption.err.
func (e *ErrInvalidOption) Error() string {
	if e.err == nil {
		e.err = errExpectedErrIsNil("ErrInvalidOption")
	}
	return e.err.Error()
}

// Unwrap returns an origin error of ErrInvalidOption.
func (e *ErrInvalidOption) Unwrap() error {
	return e.origin
}

// ErrCriticalOption represents the critical option error.
type ErrCriticalOption struct {
	err    error
	origin error
}

// NewErrCriticalOption represents a function to generate a new error of ErrCriticalOption that invalid option.
func NewErrCriticalOption(name string, val any, errs ...error) error {
	if len(errs) == 0 {
		return &ErrCriticalOption{
			err: Errorf("invalid critical option, name: %s, val: %v", name, val),
		}
	}

	e := combineOptionErrors(errs)

	return &ErrCriticalOption{
		err:    Wrapf(e, "invalid critical option, name: %s, val: %v", name, val),
		origin: e,
	}
}

// Error returns a string of ErrCriticalOption.err.
func (e *ErrCriticalOption) Error() string {
	if e.err == nil {
		e.err = errExpectedErrIsNil("ErrCriticalOption")
	}
	return e.err.Error()
}

// Unwrap returns an origin error of ErrCriticalOption.
func (e *ErrCriticalOption) Unwrap() error {
	return e.origin
}

// ErrIgnoredOption represents the ignored option error.
type ErrIgnoredOption struct {
	err    error
	origin error
}

// NewErrIgnoredOption represents a function to generate a new error of ErrIgnoredOption that option is ignored.
func NewErrIgnoredOption(name string, errs ...error) error {
	if len(errs) == 0 {
		return &ErrIgnoredOption{
			err: Errorf("ignored option, name: %s", name),
		}
	}
	e := combineOptionErrors(errs)

	return &ErrIgnoredOption{
		err:    Wrapf(e, "ignored option, name: %s", name),
		origin: e,
	}
}

// Error returns a string of ErrIgnoredOption.err.
func (e *ErrIgnoredOption) Error() string {
	if e.err == nil {
		e.err = errExpectedErrIsNil("ErrIgnoredOption")
	}
	return e.err.Error()
}

// Unwrap returns an origin error of ErrIgnoredOption.
func (e *ErrIgnoredOption) Unwrap() error {
	return e.origin
}

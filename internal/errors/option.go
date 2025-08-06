// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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

// InvalidOptionError represents the invalid option error.
type InvalidOptionError struct {
	err    error
	origin error
}

// NewErrInvalidOption represents a function to generate a new error of InvalidOptionError that invalid option.
func NewErrInvalidOption(name string, val any, errs ...error) error {
	if len(errs) == 0 {
		return &InvalidOptionError{
			err: Errorf("invalid option, name: %s, val: %v", name, val),
		}
	}
	var e error
	for _, err := range errs {
		if err == nil {
			continue
		}

		if e != nil {
			e = Wrap(err, e.Error())
		} else {
			e = err
		}
	}

	return &InvalidOptionError{
		err:    Wrapf(e, "invalid option, name: %s, val: %v", name, val),
		origin: e,
	}
}

// Error returns a string of InvalidOptionError.err.
func (e *InvalidOptionError) Error() string {
	if e.err == nil {
		e.err = errExpectedErrIsNil("InvalidOptionError")
	}
	return e.err.Error()
}

// Unwrap returns an origin error of InvalidOptionError.
func (e *InvalidOptionError) Unwrap() error {
	return e.origin
}

// CriticalOptionError represents the critical option error.
type CriticalOptionError struct {
	err    error
	origin error
}

// NewErrCriticalOption represents a function to generate a new error of CriticalOptionError that invalid option.
func NewErrCriticalOption(name string, val any, errs ...error) error {
	if len(errs) == 0 {
		return &CriticalOptionError{
			err: Errorf("invalid critical option, name: %s, val: %v", name, val),
		}
	}

	var e error
	for _, err := range errs {
		if err == nil {
			continue
		}

		if e != nil {
			e = Wrap(err, e.Error())
		} else {
			e = err
		}
	}

	return &CriticalOptionError{
		err:    Wrapf(e, "invalid critical option, name: %s, val: %v", name, val),
		origin: e,
	}
}

// Error returns a string of CriticalOptionError.err.
func (e *CriticalOptionError) Error() string {
	if e.err == nil {
		e.err = errExpectedErrIsNil("CriticalOptionError")
	}
	return e.err.Error()
}

// Unwrap returns an origin error of CriticalOptionError.
func (e *CriticalOptionError) Unwrap() error {
	return e.origin
}

// IgnoredOptionError represents the ignored option error.
type IgnoredOptionError struct {
	err    error
	origin error
}

// NewErrIgnoredOption represents a function to generate a new error of IgnoredOptionError that option is ignored.
func NewErrIgnoredOption(name string, errs ...error) error {
	if len(errs) == 0 {
		return &IgnoredOptionError{
			err: Errorf("ignored option, name: %s", name),
		}
	}
	var e error
	for _, err := range errs {
		if err == nil {
			continue
		}

		if e != nil {
			e = Wrap(err, e.Error())
		} else {
			e = err
		}
	}

	return &IgnoredOptionError{
		err:    Wrapf(e, "ignored option, name: %s", name),
		origin: e,
	}
}

// Error returns a string of IgnoredOptionError.err.
func (e *IgnoredOptionError) Error() string {
	if e.err == nil {
		e.err = errExpectedErrIsNil("IgnoredOptionError")
	}
	return e.err.Error()
}

// Unwrap returns an origin error of IgnoredOptionError.
func (e *IgnoredOptionError) Unwrap() error {
	return e.origin
}

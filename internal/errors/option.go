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
package errors

// ErrInvalidOption represents the invalid option error.
type ErrInvalidOption struct {
	err    error
	origin error
}

// NewErrInvalidOption represents a function to generate a new error of ErrInvalidOption that invalid option.
func NewErrInvalidOption(name string, val interface{}, errs ...error) error {
	if len(errs) == 0 {
		return &ErrInvalidOption{
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

	return &ErrInvalidOption{
		err:    Wrapf(e, "invalid option, name: %s, val: %v", name, val),
		origin: e,
	}
}

// Error returns a string of ErrInvalidOption.err.
func (e *ErrInvalidOption) Error() string {
	if e.err == nil {
		return ""
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
func NewErrCriticalOption(name string, val interface{}, errs ...error) error {
	if len(errs) == 0 {
		return &ErrCriticalOption{
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

	return &ErrCriticalOption{
		err:    Wrapf(e, "invalid critical option, name: %s, val: %v", name, val),
		origin: e,
	}
}

// Error returns a string of ErrCriticalOption.err.
func (e *ErrCriticalOption) Error() string {
	if e.err == nil {
		return ""
	}
	return e.err.Error()
}

// Unwrap returns an origin error of ErrCriticalOption.
func (e *ErrCriticalOption) Unwrap() error {
	return e.origin
}

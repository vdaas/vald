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

// ErrInvalidOption represent the invalid option error.
type ErrInvalidOption struct {
	err    error
	origin error
}

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

func (e *ErrInvalidOption) Error() string {
	return e.err.Error()
}

func (e *ErrInvalidOption) Unwrap() error {
	return e.origin
}

/*
   ErrCriticalOption
*/

// ErrCriticalOption represent the critical option error.
type ErrCriticalOption struct {
	err    error
	origin error
}

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

func (e *ErrCriticalOption) Error() string {
	return e.err.Error()
}

func (e *ErrCriticalOption) Unwrap() error {
	return e.origin
}

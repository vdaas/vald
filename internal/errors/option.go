//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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

import (
	"fmt"
	"strings"
)

// ErrInvalidOption represent the invalid option error
type ErrInvalidOption struct {
	name string
	val  interface{}
	errs []error
}

func NewErrInvalidOption(name string, val interface{}, errs ...error) error {
	return &ErrInvalidOption{
		name: name,
		val:  val,
		errs: errs,
	}
}

func (e *ErrInvalidOption) Error() string {
	if len(e.errs) == 0 {
		return fmt.Sprintf("invalid option, name: %s, val: %#v", e.name, e.val)
	}

	errStrs := make([]string, 0, len(e.errs))
	for i := 0; i < len(e.errs); i++ {
		errStrs[i] = e.errs[i].Error()
	}

	return fmt.Sprintf("invalid option, name: %s, val: %#v, error: %v", e.name, e.val, strings.Join(errStrs, ", "))
}

/*
   ErrCriticalOption
*/

// ErrCriticalOption represent the critical option error
type ErrCriticalOption struct {
	err error
}

func NewErrCriticalOption(name string, val interface{}, errs ...error) error {
	return &ErrCriticalOption{
		err: NewErrInvalidOption(name, val, errs...),
	}
}

func (e *ErrCriticalOption) Error() string {
	return Wrap(e.err, "invalid critical option").Error()
}

func IsCriticalOptionError(err error) bool {
	switch err.(type) {
	case *ErrCriticalOption:
		return true
	default:
		return false
	}
}

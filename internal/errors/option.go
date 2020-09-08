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

var (
	// ErrInvalidOption returns invalid option error.
	ErrInvalidOption = func(name string, val interface{}) error {
		if val == nil {
			return Errorf("invalid option. name: %s, val: nil", name)
		}
		return Errorf("invalid option. name: %s, val: %#v", name, val)
	}
)

// ErrCriticalOption represent the critical option error
type ErrCriticalOption struct {
	Err error
}

func (e *ErrCriticalOption) Error() string {
	return e.Err.Error()
}

func IsCriticalOptionError(err error) bool {
	switch err.(type) {
	case *ErrCriticalOption:
		return true
	default:
		return false
	}
}

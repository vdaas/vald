//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

import "runtime"

var (
	// ErrPanicRecovered represents a function to generate an error that the panic recovered.
	ErrPanicRecovered = func(err error, rec interface{}) error {
		return Wrap(err, Errorf("panic recovered: %v", rec).Error())
	}

	// ErrPanicString represents a function to generate an error that the panic recovered with a string message.
	ErrPanicString = func(err error, msg string) error {
		return Wrap(err, Errorf("panic recovered: %v", msg).Error())
	}

	// ErrRuntimeError represents a function to generate an error that the panic caused by runtime error.
	ErrRuntimeError = func(err error, r runtime.Error) error {
		return Wrap(err, Errorf("system panicked caused by runtime error: %v", r).Error())
	}
)

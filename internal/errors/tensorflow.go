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

var (
	ErrFailedToCastTF = func(v interface{}) error {
		return Errorf("failed to cast tensorflow result %+v", v)
	}
	ErrInputLength = func(i int, f int) error {
		return Errorf("inputs length %d does not match feeds length %d", i, f)
	}
	ErrNilTensorTF = func(v interface{}) error {
		return Errorf("nil tensorflow tensor %+v", v)
	}
	ErrNilTensorValueTF = func(v interface{}) error {
		return Errorf("nil tensorflow tensor value %+v", v)
	}
)

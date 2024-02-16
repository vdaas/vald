// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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
package mock

// Retry represents struct of mock retry structure.
type Retry struct {
	OutFunc func(
		fn func(vals ...interface{}) error,
		vals ...interface{},
	)

	OutfFunc func(
		fn func(format string, vals ...interface{}) error,
		format string,
		vals ...interface{},
	)
}

// Out calls OutFunc.
func (r *Retry) Out(
	fn func(vals ...interface{}) error,
	vals ...interface{},
) {
	r.OutFunc(fn, vals...)
}

// Outf calls OutfFunc.
func (r *Retry) Outf(
	fn func(format string, vals ...interface{}) error,
	format string, vals ...interface{},
) {
	r.OutfFunc(fn, format, vals...)
}

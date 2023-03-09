//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

// Package errors provides benchmark error
package errors

var (
	ErrInvalidCoreMode = New("invalid core mode")

	// ErrFailedToCreateBenchmarkJob represents a function to generate an error that failed to create benchmark job crd.
	ErrFailedToCreateBenchmarkJob = func(err error, jn string) error {
		return Wrapf(err, "could not create benchmark job resource: %s ", jn)
	}

	// ErrFailedToCreateJob represents a function to generate an error that failed to create job resource.
	ErrFailedToCreateJob = func(err error, jn string) error {
		return Wrapf(err, "could not create job: %s ", jn)
	}

	// ErrMismatchAtomics represents a function to generate an error that mismatch each atomic.Pointer stored corresponding to benchmark tasks.
	ErrMismatchAtomics = func(job, benchjob, benchscenario interface{}) error {
		return Errorf("mismatch atomics: job=%v\tbenchjob=%v\tbenchscenario=%v", job, benchjob, benchscenario)
	}
)

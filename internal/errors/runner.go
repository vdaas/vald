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

// Package errors provides error types and function
package errors

var (
	// ErrDaemonStartFailed represents a function to generate an error that failed to start daemon.
	ErrDaemonStartFailed = func(err error) error {
		return Wrap(err, "failed to start daemon")
	}

	// ErrDaemonStopFailed represents a function to generate an error that failed to stop daemon.
	ErrDaemonStopFailed = func(err error) error {
		return Wrap(err, "failed to stop daemon")
	}

	// ErrStartFunc represents a function to generate an error that occurred in the start function.
	ErrStartFunc = func(name string, err error) error {
		return Wrapf(err, "error occurred in runner.Start at %s", name)
	}

	// ErrPreStopFunc represents a function to generate an error that occurred in the pre-stop function.
	ErrPreStopFunc = func(name string, err error) error {
		return Wrapf(err, "error occurred in runner.PreStop at %s", name)
	}

	// ErrStopFunc represents a function to generate an error that occurred in the stop function.
	ErrStopFunc = func(name string, err error) error {
		return Wrapf(err, "error occurred in runner.Stop at %s", name)
	}

	// ErrPostStopFunc represents a function to generate an error that occurred in the post-stop function.
	ErrPostStopFunc = func(name string, err error) error {
		return Wrapf(err, "error occurred in runner.PostStop at %s", name)
	}

	// ErrRunnerWait represents a function to generate an error during runner.Wait.
	ErrRunnerWait = func(name string, err error) error {
		return Wrapf(err, "error occurred in runner.Wait at %s", name)
	}
)

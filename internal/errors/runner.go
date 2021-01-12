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
	ErrDaemonStartFailed = func(err error) error {
		return Wrap(err, "failed to start daemon")
	}

	ErrDaemonStopFailed = func(err error) error {
		return Wrap(err, "failed to stop daemon")
	}

	ErrStartFunc = func(name string, err error) error {
		return Wrapf(err, "error occurred in runner.Start at %s", name)
	}

	ErrPreStopFunc = func(name string, err error) error {
		return Wrapf(err, "error occurred in runner.PreStop at %s", name)
	}

	ErrStopFunc = func(name string, err error) error {
		return Wrapf(err, "error occurred in runner.Stop at %s", name)
	}

	ErrPostStopFunc = func(name string, err error) error {
		return Wrapf(err, "error occurred in runner.PostStop at %s", name)
	}

	ErrRunnerWait = func(name string, err error) error {
		return Wrapf(err, "error occurred in runner.Wait at %s", name)
	}
)

//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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

	// ErrWorkerIsNotRunning represents a function to generate worker is not running error.
	ErrWorkerIsNotRunning = func(name string) error {
		return Errorf("worker %s is not running", name)
	}

	// ErrWorkerIsAlreadyRunning represents a function to generate worker is already running error.
	ErrWorkerIsAlreadyRunning = func(name string) error {
		return Errorf("worker %s is already running", name)
	}

	// ErrQueueIsNotRunning represents a function to generate the queue is not running error.
	ErrQueueIsNotRunning = func() error {
		return New("queue is not running")
	}

	// ErrQueueIsAlreadyRunning represents a function to generate the queue is already running error.
	ErrQueueIsAlreadyRunning = func() error {
		return New("queue is already running")
	}

	// ErrJobFuncIsNil represents a function to generate job function is nil error.
	ErrJobFuncIsNil = func() error {
		return New("JobFunc is nil")
	}

	// ErrJobFuncNotFound represents a function to generate job function is not found.
	ErrJobFuncNotFound = func() error {
		return New("JobFunc is not found")
	}
)

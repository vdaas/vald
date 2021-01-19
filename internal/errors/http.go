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

import "time"

var (
	// ErrInvalidAPIConfig represents an error that the API configuration is invalid.
	ErrInvalidAPIConfig = New("invalid api config")

	// ErrInvalidRequest represents an error that the API request is invalid.
	ErrInvalidRequest = New("invalid request")

	// ErrHandler represents a function to generate an error that the handler returned an error.
	ErrHandler = func(err error) error {
		return Wrap(err, "handler returned error")
	}

	// ErrHandlerTimeout represents a function to generate an error that the handler was time out.
	ErrHandlerTimeout = func(err error, dur time.Duration) error {
		return Wrapf(err, "handler timeout %s", dur.String())
	}

	// ErrRequestBodyCloseAndFlush represents a function to generate an error that the flush of the request body and the close failed.
	ErrRequestBodyCloseAndFlush = func(err error) error {
		return Wrap(err, "request body flush & close failed")
	}

	// ErrRequestBodyClose represents a function to generate an error that the close of the request body failed.
	ErrRequestBodyClose = func(err error) error {
		return Wrap(err, "request body close failed")
	}

	// ErrRequestBodyFlush represents a function to generate an error that the flush of the request body failed.
	ErrRequestBodyFlush = func(err error) error {
		return Wrap(err, "request body flush failed")
	}

	// ErrTransportRetryable represents an error that the transport is retryable.
	ErrTransportRetryable = New("transport is retryable")
)

// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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
package errors

var (
	// ErrCircuitBreakerHalfOpenFlowLimitation is returned in case of flow limitation in half-open state.
	ErrCircuitBreakerHalfOpenFlowLimitation = New("circuitbreaker breaker half-open flow limitation")
	// ErrCircuitBreakerOpenState is returned when the CB state is open.
	ErrCircuitBreakerOpenState = New("circuit breaker is open")
)

type CircuitBreakerIgnorableError struct {
	err error
}

func NewCircuitBreakerIgnorableError(err error) error {
	if err == nil {
		return nil
	}
	return &CircuitBreakerIgnorableError{
		err: err,
	}
}

func (e *CircuitBreakerIgnorableError) Error() string {
	var errstr string
	if e.err != nil {
		errstr = ": " + e.err.Error()
	}
	return "circuit breaker ignorable error" + errstr
}

func (e *CircuitBreakerIgnorableError) Unwrap() error {
	return e.err
}

type CircuitBreakerMarkWithSuccessError struct {
	err error
}

func NewCircuitBreakerMarkWithSuccessError(err error) error {
	if err == nil {
		return nil
	}
	return &CircuitBreakerMarkWithSuccessError{
		err: err,
	}
}

func (e *CircuitBreakerMarkWithSuccessError) Error() string {
	var errstr string
	if e.err != nil {
		errstr = ": " + e.err.Error()
	}
	return "circuit breaker mark with success" + errstr
}

func (e *CircuitBreakerMarkWithSuccessError) Unwrap() error {
	return e.err
}

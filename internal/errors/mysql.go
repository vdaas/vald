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

var (
	// ErrMySQLConnectionPingFailed represents an error that the ping failed.
	ErrMySQLConnectionPingFailed = New("error MySQL connection ping failed")

	// ErrMySQLNotFoundIdentity generates an error that the element is not found.
	ErrMySQLNotFoundIdentity = &MySQLNotFoundIdentityError{
		err: New("error mysql element not found"),
	}

	// ErrMySQLConnectionClosed represents a function to generate an error that the connection closed.
	ErrMySQLConnectionClosed = New("error MySQL connection closed")

	// ErrMySQLTransactionNotCreated represents an error that the transaction is not closed.
	ErrMySQLTransactionNotCreated = New("error MySQL transaction not created")

	// ErrRequiredElementNotFoundByUUID represents a function to generate an error that the required element is not found.
	ErrRequiredElementNotFoundByUUID = func(uuid string) error {
		return Wrapf(ErrMySQLNotFoundIdentity, "error required element not found, uuid: %s", uuid)
	}

	// ErrMySQLInvalidArgumentIdentity generates an error that the argument is invalid.
	ErrMySQLInvalidArgumentIdentity = &MySQLInvalidArgumentIdentityError{
		err: New("error mysql invalid argument"),
	}

	// ErrRequiredMemberNotFilled represents a function to generate an error that the required member is not filled.
	ErrRequiredMemberNotFilled = func(member string) error {
		return Wrapf(ErrMySQLInvalidArgumentIdentity, "error required member not filled (member: %s)", member)
	}

	// ErrMySQLSessionNil represents a function to generate an error that the MySQL session is nil.
	ErrMySQLSessionNil = New("error MySQL session is nil")
)

// MySQLNotFoundIdentityError represents a custom error type that the element is not found.
type MySQLNotFoundIdentityError struct {
	err error
}

// Error returns the string of internal error.
func (e *MySQLNotFoundIdentityError) Error() string {
	return e.err.Error()
}

// Unwrap returns the internal error.
func (e *MySQLNotFoundIdentityError) Unwrap() error {
	return e.err
}

// IsMySQLNotFoundError returns true when the err type is MySQLNotFoundIdentityError.
func IsMySQLNotFoundError(err error) bool {
	target := new(MySQLNotFoundIdentityError)
	return As(err, &target)
}

// MySQLInvalidArgumentIdentityError represents a custom error type that the argument is not found.
type MySQLInvalidArgumentIdentityError struct {
	err error
}

// Error returns the string of internal error.
func (e *MySQLInvalidArgumentIdentityError) Error() string {
	return e.err.Error()
}

// Unwrap returns the internal error.
func (e *MySQLInvalidArgumentIdentityError) Unwrap() error {
	return e.err
}

// IsMySQLInvalidArgumentError returns true when the err type is MySQLInvalidArgumentIdentityError.
func IsMySQLInvalidArgumentError(err error) bool {
	target := new(MySQLInvalidArgumentIdentityError)
	return As(err, &target)
}

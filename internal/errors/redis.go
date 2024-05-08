//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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

	// ErrRedisInvalidKVVKPrefix represents a function to generate an error that kv index and vk prefix are invalid.
	ErrRedisInvalidKVVKPrefix = func(kv, vk string) error {
		return Errorf("kv index and vk prefix must be defferent.\t(kv: %s,\tvk: %s)", kv, vk)
	}

	// ErrRedisNotFoundIdentity generates an RedisNotFoundIdentityError error.
	ErrRedisNotFoundIdentity = &RedisNotFoundIdentityError{
		err: New("error redis entry not found"),
	}

	// ErrRedisNotFound represents a function to wrap Redis key not found error and err.
	ErrRedisNotFound = func(key string) error {
		return Wrapf(ErrRedisNotFoundIdentity, "error redis key '%s' not found", key)
	}

	// ErrRedisInvalidOption generates a new error of Redis invalid option.
	ErrRedisInvalidOption = New("error redis invalid option")

	// ErrRedisGetOperationFailed represents a function to wrap failed to fetch key error and err.
	ErrRedisGetOperationFailed = func(key string, err error) error {
		return Wrapf(err, "Failed to fetch key (%s)", key)
	}

	// ErrRedisSetOperationFailed represents a function to wrap failed to set key error and err.
	ErrRedisSetOperationFailed = func(key string, err error) error {
		return Wrapf(err, "Failed to set key (%s)", key)
	}

	// ErrRedisSetOperationFailed represents a function to wrap failed to delete key error and err.
	ErrRedisDeleteOperationFailed = func(key string, err error) error {
		return Wrapf(err, "Failed to delete key (%s)", key)
	}

	// ErrRedisSetOperationFailed represents a function to generate an error that invalid configuration version.
	ErrInvalidConfigVersion = func(cur, con string) error {
		return Errorf("invalid config version %s not satisfies version constraints %s", cur, con)
	}

	// ErrRedisAddrsNotFound generates a new error of address not found.
	ErrRedisAddrsNotFound = New("error redis addrs not found")

	// ErrRedisConnectionPingFailed generates a new error of Redis connection ping failed.
	ErrRedisConnectionPingFailed = New("error redis connection ping failed")
)

// RedisNotFoundIdentityError represents a struct that includes err and has a method for Redis error handling.
type RedisNotFoundIdentityError struct {
	err error
}

// Error returns the string of ErrRedisNotFoundIdentity.error.
func (e *RedisNotFoundIdentityError) Error() string {
	if e.err == nil {
		e.err = errExpectedErrIsNil("ErrRedisNotFoundIdentity")
	}
	return e.err.Error()
}

// Unwrap returns the error value of RedisNotFoundIdentityError.
func (e *RedisNotFoundIdentityError) Unwrap() error {
	return e.err
}

// IsRedisNotFoundError compares the input error and RedisNotFoundIdentityError.error and returns true or false that is the result of errors.As.
func IsRedisNotFoundError(err error) bool {
	target := new(RedisNotFoundIdentityError)
	return As(err, &target)
}

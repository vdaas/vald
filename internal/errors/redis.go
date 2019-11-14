//
// Copyright (C) 2019 Vdaas.org Vald team ( kpango, kou-m, rinx )
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Package errors provides error types and function
package errors

// "github.com/pkg/errors"

var (

	// Redis
	ErrRedisInvalidKVVKPrefix = func(kv, vk string) error {
		return Errorf("kv index and vk prefix must be defferent.\t(kv: %s,\tvk: %s)", kv, vk)
	}

	ErrRedisGetOperationFailed = func(key string, err error) error {
		return Wrapf(err, "Failed to fetch key (%s)", key)
	}

	ErrRedisSetOperationFailed = func(key string, err error) error {
		return Wrapf(err, "Failed to set key (%s)", key)
	}

	ErrRedisDeleteOperationFailed = func(key string, err error) error {
		return Wrapf(err, "Failed to delete key (%s)", key)
	}

	ErrInvalidConfigVersion = func(cur, con string) error {
		return Errorf("invalid config version %s not satisfies version constraints %s", cur, con)
	}

	ErrRedisAddrsNotFound = New("addrs not found")
)

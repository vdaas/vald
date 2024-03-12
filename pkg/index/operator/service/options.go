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
package service

import (
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

// Option represents the functional option for index.
type Option func(_ *operator) error

var defaultOpts = []Option{
	WithErrGroup(errgroup.Get()),
	WithRotationJobConcurrency(1),
}

func WithErrGroup(eg errgroup.Group) Option {
	return func(o *operator) error {
		if eg != nil {
			o.eg = eg
		}
		return nil
	}
}

func WithReadReplicaEnabled(enabled bool) Option {
	return func(o *operator) error {
		o.readReplicaEnabled = enabled
		return nil
	}
}

func WithReadReplicaLabelKey(key string) Option {
	return func(o *operator) error {
		o.readReplicaLabelKey = key
		return nil
	}
}

func WithRotationJobConcurrency(concurrency uint) Option {
	return func(o *operator) error {
		if concurrency == 0 {
			return errors.NewErrCriticalOption("RotationJobConcurrency", concurrency, errors.New("concurrency should be greater than 0"))
		}
		o.rotationJobConcurrency = concurrency
		return nil
	}
}

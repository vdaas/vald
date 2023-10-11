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

// Package service manages the main logic of benchmark job.
package service

import (
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

// Option represents the functional option for scenario struct.
type Option func(o *operator) error

var defaultOpts = []Option{
	WithJobImage("vdaas/vald-benchmark-job"),
	WithJobImagePullPolicy("Always"),
	WithReconcileCheckDuration("10s"),
	WithJobNamespace("default"),
}

// WithErrGroup sets the error group to scenario.
func WithErrGroup(eg errgroup.Group) Option {
	return func(o *operator) error {
		if eg == nil {
			return errors.NewErrInvalidOption("client", eg)
		}
		o.eg = eg
		return nil
	}
}

// WithReconcileCheckDuration sets the reconcile check duration from input string.
func WithReconcileCheckDuration(ts string) Option {
	return func(o *operator) error {
		t, err := time.ParseDuration(ts)
		if err != nil {
			return err
		}
		o.rcd = t
		return nil
	}
}

// WithJobNamespace sets the namespace for running benchmark job.
func WithJobNamespace(ns string) Option {
	return func(o *operator) error {
		if len(ns) == 0 {
			o.jobNamespace = "default"
		} else {
			o.jobNamespace = ns
		}
		return nil
	}
}

// WithJobImage sets the benchmark job docker image info.
func WithJobImage(image string) Option {
	return func(o *operator) error {
		if len(image) > 0 {
			o.jobImage = image
		}
		return nil
	}
}

// WithJobImagePullPolicy sets the benchmark job docker image pullPolicy.
func WithJobImagePullPolicy(p string) Option {
	return func(o *operator) error {
		if len(p) > 0 {
			o.jobImagePullPolicy = p
		}
		return nil
	}
}

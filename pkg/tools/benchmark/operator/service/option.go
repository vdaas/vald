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

// Package service manages the main logic of benchmark job.
package service

import (
	"time"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
)

// Option represents the functional option for scenario struct.
type Option func(sc *scenario) error

var defaultOpts = []Option{
	WithReconcileCheckDuration("10s"),
}

// WithErrGroup sets the error group to scenario.
func WithErrGroup(eg errgroup.Group) Option {
	return func(sc *scenario) error {
		if eg == nil {
			return errors.NewErrInvalidOption("client", eg)
		}
		sc.eg = eg
		return nil
	}
}

// WithReconcileCheckDuration sets the reconcile check duration from input string.
func WithReconcileCheckDuration(ts string) Option {
	return func(sc *scenario) error {
		t, err := time.ParseDuration(ts)
		if err != nil {
			return err
		}
		sc.rcd = t
		return nil
	}
}

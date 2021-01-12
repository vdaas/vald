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

// Package service
package service

import (
	"time"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/timeutil"
)

type Option func(r *replicator) error

var defaultOptions = []Option{
	WithRecoverCheckDuration("1m"),
	WithErrGroup(errgroup.Get()),
	WithNamespace("vald"),
}

func WithName(name string) Option {
	return func(r *replicator) error {
		if len(name) != 0 {
			r.name = name
		}
		return nil
	}
}

func WithNamespace(ns string) Option {
	return func(r *replicator) error {
		if len(ns) != 0 {
			r.namespace = ns
		}
		return nil
	}
}

func WithRecoverCheckDuration(dur string) Option {
	return func(r *replicator) error {
		if dur == "" {
			return nil
		}
		pd, err := timeutil.Parse(dur)
		if err != nil {
			pd = time.Second
		}
		r.rdur = pd
		return nil
	}
}

func WithErrGroup(eg errgroup.Group) Option {
	return func(r *replicator) error {
		if eg != nil {
			r.eg = eg
		}
		return nil
	}
}

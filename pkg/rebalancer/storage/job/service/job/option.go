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
package job

import (
	"strconv"

	"github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/pkg/rebalancer/storage/job/service/storage"
)

type Option func(r *rebalancer) error

var defaultOpts = []Option{}

func WithStorage(st storage.Storage) Option {
	return func(r *rebalancer) error {
		r.storage = st
		return nil
	}
}

func WithTargetAgentName(name string) Option {
	return func(r *rebalancer) error {
		r.targetAgentName = name
		return nil
	}
}

func WithRate(rate string) Option {
	return func(r *rebalancer) (err error) {
		r.rate, err = strconv.ParseFloat(rate, 64)
		if err != nil {
			return errors.NewErrInvalidOption("rate", rate)
		}
		return nil
	}
}

func WithValdClient(c vald.Client) Option {
	return func(r *rebalancer) error {
		r.client = c
		return nil
	}
}

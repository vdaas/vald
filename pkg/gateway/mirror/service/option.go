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

// Package service represents gateway's service logic
package service

import (
	"time"

	"github.com/vdaas/vald/internal/client/v1/client/mirror"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
)

type Option func(g *gateway) error

var defaultGWOpts = []Option{
	WithErrGroup(errgroup.Get()),
	WithAdvertiseInterval("1s"),
}

func WithMirror(c mirror.Client) Option {
	return func(g *gateway) error {
		if c != nil {
			g.client = c
		}
		return nil
	}
}

func WithSelfMirror(c mirror.Client) Option {
	return func(g *gateway) error {
		if c != nil {
			g.iclient = c
		}
		return nil
	}
}

func WithErrGroup(eg errgroup.Group) Option {
	return func(g *gateway) error {
		if eg != nil {
			g.eg = eg
		}
		return nil
	}
}

func WithAdvertiseInterval(s string) Option {
	return func(g *gateway) error {
		if len(s) == 0 {
			return errors.NewErrInvalidOption("advertiseInterval", s)
		}
		dur, err := time.ParseDuration(s)
		if err != nil {
			return errors.NewErrInvalidOption("advertiseInterval", s, err)
		}
		g.advertiseDur = dur
		return nil
	}
}

func WithPodName(s string) Option {
	return func(g *gateway) error {
		if len(s) == 0 {
			return errors.NewErrCriticalOption("podName", s)
		}
		g.podName = s
		return nil
	}
}

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
	"github.com/vdaas/vald/internal/client/v1/client/mirror"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

// Option represents the functional option for gateway.
type Option func(g *gateway) error

var defaultGatewayOpts = []Option{
	WithErrGroup(errgroup.Get()),
}

// WithMirrorClient returns the option to set the Mirror client.
func WithMirrorClient(c mirror.Client) Option {
	return func(g *gateway) error {
		if c != nil {
			g.client = c
		}
		return nil
	}
}

// WithErrGroup returns the option to set the error group.
func WithErrGroup(eg errgroup.Group) Option {
	return func(g *gateway) error {
		if eg != nil {
			g.eg = eg
		}
		return nil
	}
}

// WithPodName returns the option to set the pod name.
func WithPodName(s string) Option {
	return func(g *gateway) error {
		if s == "" {
			return errors.NewErrCriticalOption("podName", s)
		}
		g.podName = s
		return nil
	}
}

// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

// MirrorOption represents the functional option for mirror.
type MirrorOption func(m *mirr) error

var defaultMirrOpts = []MirrorOption{
	WithRegisterDuration("500ms"),
}

// WithErrorGroup returns the option to set the error group.
func WithErrorGroup(eg errgroup.Group) MirrorOption {
	return func(m *mirr) error {
		if eg != nil {
			m.eg = eg
		}
		return nil
	}
}

// WithGatewayAddrs returns the option to set the gateway addresses.
func WithGatewayAddrs(addrs ...string) MirrorOption {
	return func(m *mirr) error {
		if len(addrs) == 0 {
			return errors.NewErrCriticalOption("lbAddrs", addrs)
		}
		for _, addr := range addrs {
			m.gwAddrl.Store(addr, struct{}{})
		}
		return nil
	}
}

// WithSelfMirrorAddrs returns the option to set the self Mirror Gateway addresses.
func WithSelfMirrorAddrs(addrs ...string) MirrorOption {
	return func(m *mirr) error {
		if len(addrs) == 0 {
			return errors.NewErrCriticalOption("selfMirrorAddrs", addrs)
		}
		for _, addr := range addrs {
			m.selfMirrAddrl.Store(addr, struct{}{})
		}
		return nil
	}
}

// WithGatewayAddrs returns the option to set the Gateway service.
func WithGateway(g Gateway) MirrorOption {
	return func(m *mirr) error {
		if g != nil {
			m.gateway = g
		}
		return nil
	}
}

// WithRegisterDuration returns the option to set the register duration.
func WithRegisterDuration(s string) MirrorOption {
	return func(m *mirr) error {
		if s == "" {
			return errors.NewErrInvalidOption("registerDuration", s)
		}
		dur, err := time.ParseDuration(s)
		if err != nil {
			return errors.NewErrInvalidOption("registerDuration", s, err)
		}
		m.registerDur = dur
		return nil
	}
}

// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
)

type MirrorOption func(m *mirr) error

var defaultMirrOpts = []MirrorOption{
	WithAdvertiseInterval("1s"),
}

func WithErrorGroup(eg errgroup.Group) MirrorOption {
	return func(m *mirr) error {
		if eg != nil {
			m.eg = eg
		}
		return nil
	}
}

func WithValdAddrs(addrs ...string) MirrorOption {
	return func(m *mirr) error {
		if len(addrs) == 0 {
			return errors.NewErrCriticalOption("lbAddrs", addrs)
		}
		if m.gwAddrs == nil {
			m.gwAddrs = make([]string, 0, len(addrs))
		}
		m.gwAddrs = append(m.gwAddrs, addrs...)
		return nil
	}
}

func WithSelfMirrorAddrs(addrs ...string) MirrorOption {
	return func(m *mirr) error {
		if len(addrs) == 0 {
			return errors.NewErrCriticalOption("selfMirrorAddrs", addrs)
		}
		if m.selfMirrAddrs == nil {
			m.selfMirrAddrs = make([]string, 0, len(addrs))
		}
		m.selfMirrAddrs = append(m.selfMirrAddrs, addrs...)
		return nil
	}
}

func WithGateway(g Gateway) MirrorOption {
	return func(m *mirr) error {
		if g != nil {
			m.gateway = g
		}
		return nil
	}
}

func WithAdvertiseInterval(s string) MirrorOption {
	return func(m *mirr) error {
		if len(s) == 0 {
			return errors.NewErrInvalidOption("advertiseInterval", s)
		}
		dur, err := time.ParseDuration(s)
		if err != nil {
			return errors.NewErrInvalidOption("advertiseInterval", s, err)
		}
		m.advertiseDur = dur
		return nil
	}
}

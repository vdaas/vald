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
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/internal/timeutil"
)

// DiscoveryOption represents the functional option for discovery.
type DiscoveryOption func(d *discovery) error

var defaultDiscovererOpts = []DiscoveryOption{
	WithDiscoveryDuration("1s"),
	WithDiscoveryErrGroup(errgroup.Get()),
	WithDiscoveryColocation("dc1"),
}

// WithDiscoveryMirror returns the option to set the Mirror service.
func WithDiscoveryMirror(m Mirror) DiscoveryOption {
	return func(d *discovery) error {
		if m == nil {
			return errors.NewErrCriticalOption("discoveryMirror", m)
		}
		d.mirr = m
		return nil
	}
}

// WithDiscoveryDialer returns the option to set the dialer for controller manager.
func WithDiscoveryDialer(der net.Dialer) DiscoveryOption {
	return func(d *discovery) error {
		if der != nil {
			d.der = der
		}
		return nil
	}
}

// WithDiscoveryNamespace returns the option to set the namespace for discovery.
func WithDiscoveryNamespace(ns string) DiscoveryOption {
	return func(d *discovery) error {
		if len(ns) != 0 {
			d.namespace = ns
		}
		return nil
	}
}

// WithDiscoveryGroup returns the option to set the Mirror group for discovery.
func WithDiscoveryGroup(g string) DiscoveryOption {
	return func(d *discovery) error {
		if len(g) != 0 {
			if d.labels == nil {
				d.labels = make(map[string]string)
			}
			d.labels[groupKey] = g
		}
		return nil
	}
}

// WithDiscoveryColocation returns the option to set the colocation name of datacenter.
func WithDiscoveryColocation(loc string) DiscoveryOption {
	return func(d *discovery) error {
		if len(loc) != 0 {
			d.colocation = loc
		}
		return nil
	}
}

// WithDiscoveryDuration returns the option to set the duration of the discovery.
func WithDiscoveryDuration(s string) DiscoveryOption {
	return func(d *discovery) error {
		if s == "" {
			return nil
		}
		dur, err := timeutil.Parse(s)
		if err != nil {
			return errors.NewErrInvalidOption("discoveryDuration", s, err)
		}
		d.dur = dur
		return nil
	}
}

// WithDiscoveryErrGroup returns the option to set the errgroup.
func WithDiscoveryErrGroup(eg errgroup.Group) DiscoveryOption {
	return func(d *discovery) error {
		if eg != nil {
			d.eg = eg
		}
		return nil
	}
}

// WithDiscoverySelfMirrorAddrs returns the option to set the self Mirror addresses.
func WithDiscoverySelfMirrorAddrs(addrs ...string) DiscoveryOption {
	return func(d *discovery) error {
		if len(addrs) == 0 {
			return errors.NewErrCriticalOption("discoverySelfMirrorAddrs", addrs)
		}
		d.selfMirrAddrs = append(d.selfMirrAddrs, addrs...)
		return nil
	}
}

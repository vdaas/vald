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

// Option represents the functional option for discoverer.
type DiscovererOption func(d *discoverer) error

var defaultDiscovererOpts = []DiscovererOption{
	WithDiscovererDuration("1s"),
	WithDiscovererErrGroup(errgroup.Get()),
	WithDiscovererColocation("dc1"),
}

// WithDiscovererMirror returns the option to set the Mirror service.
func WithDiscovererMirror(m Mirror) DiscovererOption {
	return func(d *discoverer) error {
		if m == nil {
			return errors.NewErrCriticalOption("discovererMirror", m)
		}
		d.mirr = m
		return nil
	}
}

// WithDiscovererDialer returns the option to set the dialer for controller manager.
func WithDiscovererDialer(der net.Dialer) DiscovererOption {
	return func(d *discoverer) error {
		if der != nil {
			d.der = der
		}
		return nil
	}
}

// WithDiscovererNamespace returns the option to set the namespace for discovery.
func WithDiscovererNamespace(ns string) DiscovererOption {
	return func(d *discoverer) error {
		if len(ns) != 0 {
			d.namespace = ns
		}
		return nil
	}
}

// WithDiscovererGroup returns the option to set the Mirror group for discovery.
func WithDiscovererGroup(g string) DiscovererOption {
	return func(d *discoverer) error {
		if len(g) != 0 {
			if d.labels == nil {
				d.labels = make(map[string]string)
			}
			d.labels[groupKey] = g
		}
		return nil
	}
}

// WithDiscovererColocation returns the option to set the colocation name of datacenter.
func WithDiscovererColocation(loc string) DiscovererOption {
	return func(d *discoverer) error {
		if len(loc) != 0 {
			d.colocation = loc
		}
		return nil
	}
}

// WithDiscovererDuration returns the option to set the duration of the discovery.
func WithDiscovererDuration(s string) DiscovererOption {
	return func(d *discoverer) error {
		if s == "" {
			return nil
		}
		dur, err := timeutil.Parse(s)
		if err != nil {
			return errors.NewErrInvalidOption("discovererDuration", s, err)
		}
		d.dur = dur
		return nil
	}
}

// WithDiscovererErrGroup returns the option to set the errgroup.
func WithDiscovererErrGroup(eg errgroup.Group) DiscovererOption {
	return func(d *discoverer) error {
		if eg != nil {
			d.eg = eg
		}
		return nil
	}
}

// WithDiscovererSelfMirrorAddrs returns the option to set the self Mirror addresses.
func WithDiscovererSelfMirrorAddrs(addrs ...string) DiscovererOption {
	return func(d *discoverer) error {
		if len(addrs) == 0 {
			return errors.NewErrCriticalOption("discovererSelfMirrorAddrs", addrs)
		}
		d.selfMirrAddrs = append(d.selfMirrAddrs, addrs...)
		return nil
	}
}

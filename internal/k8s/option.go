//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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

// Package k8s provides kubernetes control functionality
package k8s

import (
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

type Option func(*controller) error

var defaultOptions = []Option{
	WithErrGroup(errgroup.Get()),
}

func WithErrGroup(eg errgroup.Group) Option {
	return func(c *controller) error {
		if eg != nil {
			c.eg = eg
		}
		return nil
	}
}

func WithControllerName(name string) Option {
	return func(c *controller) error {
		c.name = name
		return nil
	}
}

func WithResourceController(rc ResourceController) Option {
	return func(c *controller) error {
		if c.rcs == nil {
			c.rcs = make([]ResourceController, 0, 1)
		}
		c.rcs = append(c.rcs, rc)
		return nil
	}
}

func WithMetricsAddress(addr string) Option {
	return func(c *controller) error {
		c.merticsAddr = addr
		return nil
	}
}

func WithLeaderElection(enabled bool, id, namespace string) Option {
	return func(c *controller) error {
		if enabled && id == "" {
			return errors.NewErrCriticalOption("leaderElectionID", id)
		}
		c.leaderElection = enabled
		c.leaderElectionID = id
		c.leaderElectionNamespace = namespace
		return nil
	}
}

// WithLeaderElectionDetail sets the leader election lease timings applied to
// the manager. Non-positive values are ignored so that controller-runtime
// defaults (15s/10s/2s) remain in effect for the unset fields.
func WithLeaderElectionDetail(leaseDuration, renewDeadline, retryPeriod time.Duration) Option {
	return func(c *controller) error {
		if leaseDuration > 0 {
			c.leaseDuration = &leaseDuration
		}
		if renewDeadline > 0 {
			c.renewDeadline = &renewDeadline
		}
		if retryPeriod > 0 {
			c.retryPeriod = &retryPeriod
		}
		return nil
	}
}

// WithSyncPeriod sets the informer cache resync interval. A non-positive
// value keeps the controller-runtime default (10h with jitter).
func WithSyncPeriod(dur time.Duration) Option {
	return func(c *controller) error {
		if dur > 0 {
			c.syncPeriod = &dur
		}
		return nil
	}
}

// WithCacheNamespaces restricts the manager cache to the given namespaces.
// Empty entries are skipped; when no namespace remains, the cache watches the
// whole cluster as before.
func WithCacheNamespaces(namespaces ...string) Option {
	return func(c *controller) error {
		for _, ns := range namespaces {
			if ns != "" {
				c.cacheNamespaces = append(c.cacheNamespaces, ns)
			}
		}
		return nil
	}
}

func WithDialer(der net.Dialer) Option {
	return func(c *controller) error {
		if der != nil {
			c.der = der
		}
		return nil
	}
}

//
// Copyright (C) 2019 kpango (Yusuke Kato)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Package k8s provides kubernetes controll functionallity
package k8s

import "sigs.k8s.io/controller-runtime/pkg/manager"

type Option func(*controller) error

var (
	defaultOpts = []Option{}
)

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

func WithManager(mgr manager.Manager) Option {
	return func(c *controller) error {
		c.mgr = mgr
		return nil
	}
}

func WithMetricsAddress(addr string) Option {
	return func(c *controller) error {
		c.merticsAddr = addr
		return nil
	}
}

func WithEnableLeaderElection() Option {
	return func(c *controller) error {
		c.leaderElection = true
		return nil
	}
}

func WithDisableLeaderElection() Option {
	return func(c *controller) error {
		c.leaderElection = false
		return nil
	}
}

//go:build e2e

//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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

// package portforward provides a persistent port forwarding daemon for Kubernetes services.
package portforward

import (
	"net/http"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/sync/errgroup"
	k8s "github.com/vdaas/vald/tests/v2/e2e/kubernetes"
)

// Option represents the functional option for backoff.
type Option func(*portForward)

var defaultOptions = []Option{
	WithHTTPClient(http.DefaultClient),
	WithBackoff(backoff.New()),
	WithNamespace("default"),
}

func WithClient(client k8s.Client) Option {
	return func(pf *portForward) {
		if client != nil {
			pf.client = client
		}
	}
}

func WithBackoff(bo backoff.Backoff) Option {
	return func(pf *portForward) {
		if bo != nil {
			pf.backoff = bo
		}
	}
}

func WithErrGroup(eg errgroup.Group) Option {
	return func(pf *portForward) {
		if eg != nil {
			pf.eg = eg
		}
	}
}

func WithNamespace(ns string) Option {
	return func(pf *portForward) {
		if ns != "" {
			pf.namespace = ns
		}
	}
}

func WithServiceName(name string) Option {
	return func(pf *portForward) {
		if name != "" {
			pf.serviceName = name
		}
	}
}

func WithAddress(addrs ...string) Option {
	return func(pf *portForward) {
		if addrs != nil {
			pf.addresses = addrs
		}
	}
}

func WithHTTPClient(c *http.Client) Option {
	return func(pf *portForward) {
		if c != nil {
			pf.httpClient = c
		}
	}
}

func WithPorts(pairs map[uint16]uint16) Option {
	return func(pf *portForward) {
		if pairs != nil {
			pf.ports = pairs
		}
	}
}

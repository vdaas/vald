//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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
package service

import (
	"fmt"
	"time"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/timeutil"
)

type GWOption func(g *gateway) error

var (
	defaultGWOpts = []GWOption{
		WithErrGroup(errgroup.Get()),
		WithAgentNamespace("default"),
	}
)

func WithDiscovererClient(client grpc.Client) GWOption {
	return func(g *gateway) error {
		g.dscClient = client
		return nil
	}
}

func WithDiscovererAddr(addr string) GWOption {
	return func(g *gateway) error {
		g.dscAddr = addr
		return nil
	}
}

func WithDiscovererHostPort(host string, port int) GWOption {
	return func(g *gateway) error {
		g.dscAddr = fmt.Sprintf("%s:%d", host, port)
		return nil
	}
}

func WithDiscoverDuration(dur string) GWOption {
	return func(g *gateway) error {
		d, err := timeutil.Parse(dur)
		if err != nil {
			d = time.Second
		}
		g.dscDur = d
		return nil
	}
}

func WithAgentOptions(opts ...grpc.Option) GWOption {
	return func(g *gateway) error {
		g.agentOpts = append(g.agentOpts, opts...)
		return nil
	}
}

func WithAgentName(name string) GWOption {
	return func(g *gateway) error {
		if name != "" {
			g.agentName = name
		}
		return nil
	}
}

func WithAgentNamespace(ns string) GWOption {
	return func(g *gateway) error {
		if ns != "" {
			g.namespace = ns
		}
		return nil
	}
}

func WithAgentPort(port int) GWOption {
	return func(g *gateway) error {
		g.agentPort = port
		return nil
	}
}

func WithAgentServiceDNSARecord(a string) GWOption {
	return func(g *gateway) error {
		g.agentARecord = a
		return nil
	}
}

func WithNodeName(nn string) GWOption {
	return func(g *gateway) error {
		if nn != "" {
			g.nodeName = nn
		}
		return nil
	}
}

func WithErrGroup(eg errgroup.Group) GWOption {
	return func(g *gateway) error {
		if eg != nil {
			g.eg = eg
		}
		return nil
	}
}

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

// Package jaeger provides a jaeger exporter.
package jaeger

import (
	"net"
)

type Option func(*export) error

var jaegerDefaultOpts = []Option{
	WithServiceName("vald"),
	// For tracing over network, packets must fit in MTU 1500, which has a
	// payload size of 1472. But if network is local, we can use default package size (MTU: 65000)
	// TODO: Implement the handling appropriately later.
	WithMaxPacketSize(1472),
}

func WithCollectorEndpoint(cep string) Option {
	return func(exp *export) error {
		if cep != "" {
			exp.collectorEndpoint = cep
		}
		return nil
	}
}

func WithAgentEndpoint(aep string) Option {
	return func(exp *export) error {
		if aep != "" {
			host, port, err := net.SplitHostPort(aep)
			if err != nil {
				return err
			}
			exp.agentHost = host
			exp.agentPort = port
		}
		return nil
	}
}

func WithUsername(username string) Option {
	return func(exp *export) error {
		if username != "" {
			exp.collectorUserName = username
		}
		return nil
	}
}

func WithPassword(password string) Option {
	return func(exp *export) error {
		if password != "" {
			exp.collectorPassword = password
		}
		return nil
	}
}

func WithServiceName(serviceName string) Option {
	return func(exp *export) error {
		if serviceName != "" {
			exp.serviceName = serviceName
		}
		return nil
	}
}

func WithBufferMaxCount(cnt int) Option {
	return func(exp *export) error {
		exp.agentMaxPacketSize = cnt
		return nil
	}
}

func WithMaxPacketSize(size int) Option {
	return func(exp *export) error {
		exp.agentMaxPacketSize = size
		return nil
	}
}

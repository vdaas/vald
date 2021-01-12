//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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

// Package config providers configuration type and load configuration logic
package config

type Discoverer struct {
	Name              string `json:"name" yaml:"name"`
	Namespace         string `json:"namespace" yaml:"namespace"`
	DiscoveryDuration string `json:"discovery_duration" yaml:"discovery_duration"`
}

func (d *Discoverer) Bind() *Discoverer {
	d.Name = GetActualValue(d.Name)
	d.Namespace = GetActualValue(d.Namespace)
	d.DiscoveryDuration = GetActualValue(d.DiscoveryDuration)
	return d
}

type DiscovererClient struct {
	Duration           string      `json:"duration" yaml:"duration"`
	Client             *GRPCClient `json:"client" yaml:"client"`
	AgentClientOptions *GRPCClient `json:"agent_client_options" yaml:"agent_client_options"`
}

func (d *DiscovererClient) Bind() *DiscovererClient {
	d.Duration = GetActualValue(d.Duration)
	if d.Client != nil {
		d.Client.Bind()
	} else {
		d.Client = newGRPCClientConfig()
	}
	if d.AgentClientOptions != nil {
		d.AgentClientOptions.Bind()
	} else {
		d.AgentClientOptions = newGRPCClientConfig()
	}
	return d
}

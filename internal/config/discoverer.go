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

// Package config providers configuration type and load configuration logic
package config

import "fmt"

type Discoverer struct {
	Name              string `json:"name" yaml:"name"`
	Namespace         string `json:"namespace" yaml:"namespace"`
	CacheSyncDuration string `json:"cache_sync_duration" yaml:"cache_sync_duration"`
}

func (d *Discoverer) Bind() *Discoverer {
	d.Name = GetActualValue(d.Name)
	d.Namespace = GetActualValue(d.Namespace)
	d.CacheSyncDuration = GetActualValue(d.CacheSyncDuration)
	return d
}

type DiscovererClient struct {
	Host        string      `json:"host" yaml:"host"`
	Port        int         `json:"port" yaml:"port"`
	Duration    string      `json:"duration" yaml:"duration"`
	Client      *GRPCClient `json:"discover_client" yaml:"discover_client"`
	AgentClient *GRPCClient `json:"agent_client" yaml:"agent_client"`
}

func (d *DiscovererClient) Bind() *DiscovererClient {
	d.Host = GetActualValue(d.Host)
	d.Duration = GetActualValue(d.Duration)
	if d.Client != nil {
		d.Client.Bind()
	} else {
		d.Client = newGRPCClientConfig()
	}
	if len(d.Host) != 0 {
		d.Client.Addrs = append(d.Client.Addrs, fmt.Sprintf("%s:%d", d.Host, d.Port))
	}
	if d.AgentClient != nil {
		d.AgentClient.Bind()
	} else {
		d.AgentClient = newGRPCClientConfig()
	}
	return d
}

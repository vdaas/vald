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

type Observability struct {
	Prometheus *Prometheus `json:"prometheus" yaml:"prometheus"`
	Jaeger     *Jaeger     `json:"jaeger" yaml:"jaeger"`
}

type Prometheus struct {
	Enabled   bool   `json:"enabled" yaml:"enabled"`
	Namespace string `json:"namespace" yaml:"namespace"`
}

type Jaeger struct {
	Enabled bool `json:"enabled" yaml:"enabled"`

	CollectorEndpoint string `json:"collector_endpoint" yaml:"collector_endpoint"`
	AgentEndpoint     string `json:"agent_endpoint" yaml:"agent_endpoint"`

	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`

	ServiceName string `json:"service_name" yaml:"service_name"`

	BufferMaxCount int `json:"buffer_max_count" yaml:"buffer_max_count"`
}

func (c *Observability) Bind() *Observability {
	if c.Prometheus != nil {
		c.Prometheus.Namespace = GetActualValue(c.Prometheus.Namespace)
	}

	if c.Jaeger != nil {
		c.Jaeger.CollectorEndpoint = GetActualValue(c.Jaeger.CollectorEndpoint)
		c.Jaeger.AgentEndpoint = GetActualValue(c.Jaeger.AgentEndpoint)
		c.Jaeger.Username = GetActualValue(c.Jaeger.Username)
		c.Jaeger.Password = GetActualValue(c.Jaeger.Password)
		c.Jaeger.ServiceName = GetActualValue(c.Jaeger.ServiceName)
	}

	return c
}

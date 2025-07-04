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

package config

// Discoverer represents the Discoverer configurations.
type Discoverer struct {
	Name              string     `json:"name,omitempty"               yaml:"name"`
	Namespace         string     `json:"namespace,omitempty"          yaml:"namespace"`
	DiscoveryDuration string     `json:"discovery_duration,omitempty" yaml:"discovery_duration"`
	Net               *Net       `json:"net,omitempty"                yaml:"net"`
	Selectors         *Selectors `json:"selectors,omitempty"          yaml:"selectors"`
}

type Selectors struct {
	Pod         *Selector `json:"pod,omitempty"          yaml:"pod"`
	Node        *Selector `json:"node,omitempty"         yaml:"node"`
	NodeMetrics *Selector `json:"node_metrics,omitempty" yaml:"node_metrics"`
	PodMetrics  *Selector `json:"pod_metrics,omitempty"  yaml:"pod_metrics"`
	Service     *Selector `json:"service,omitempty"      yaml:"service"`
}

func (s *Selectors) GetPodFields() map[string]string {
	if s == nil {
		return nil
	}
	return s.Pod.GetFields()
}

func (s *Selectors) GetPodLabels() map[string]string {
	if s == nil {
		return nil
	}
	return s.Pod.GetLabels()
}

func (s *Selectors) GetNodeFields() map[string]string {
	if s == nil {
		return nil
	}
	return s.Node.GetFields()
}

func (s *Selectors) GetNodeLabels() map[string]string {
	if s == nil {
		return nil
	}
	return s.Node.GetLabels()
}

func (s *Selectors) GetPodMetricsFields() map[string]string {
	if s == nil {
		return nil
	}
	return s.PodMetrics.GetFields()
}

func (s *Selectors) GetPodMetricsLabels() map[string]string {
	if s == nil {
		return nil
	}
	return s.PodMetrics.GetLabels()
}

func (s *Selectors) GetNodeMetricsFields() map[string]string {
	if s == nil {
		return nil
	}
	return s.NodeMetrics.GetFields()
}

func (s *Selectors) GetNodeMetricsLabels() map[string]string {
	if s == nil {
		return nil
	}
	return s.NodeMetrics.GetLabels()
}

func (s *Selectors) GetServiceFields() map[string]string {
	if s == nil {
		return nil
	}
	return s.Service.GetFields()
}

func (s *Selectors) GetServiceLabels() map[string]string {
	if s == nil {
		return nil
	}
	return s.Service.GetLabels()
}

type Selector struct {
	Labels map[string]string `json:"labels,omitempty" yaml:"labels"`
	Fields map[string]string `json:"fields,omitempty" yaml:"fields"`
}

func (s *Selector) GetLabels() map[string]string {
	if s == nil {
		return nil
	}
	return s.Labels
}

func (s *Selector) GetFields() map[string]string {
	if s == nil {
		return nil
	}
	return s.Fields
}

type ReadReplica struct {
	Enabled bool   `json:"enabled,omitempty" yaml:"enabled"`
	IDKey   string `json:"id_key,omitempty"  yaml:"id_key"`
}

func (r *ReadReplica) GetEnabled() bool {
	if r == nil {
		return false
	}
	return r.Enabled
}

func (r *ReadReplica) GetIDKey() string {
	if r == nil {
		return ""
	}
	return r.IDKey
}

// Bind binds the actual data from the Discoverer receiver field.
func (d *Discoverer) Bind() *Discoverer {
	d.Name = GetActualValue(d.Name)
	d.Namespace = GetActualValue(d.Namespace)
	d.DiscoveryDuration = GetActualValue(d.DiscoveryDuration)
	if d.Net != nil {
		d.Net.Bind()
	}

	if d.Selectors != nil {
		d.Selectors.Bind()
	}

	return d
}

// Bind binds the actual data from the Selectors receiver field.
func (s *Selectors) Bind() *Selectors {
	// Initialization of s itself should be handled by the parent struct's Bind method.
	if s.Pod == nil {
		s.Pod = new(Selector)
	}
	s.Pod.Bind()

	if s.Node == nil {
		s.Node = new(Selector)
	}
	s.Node.Bind()

	if s.PodMetrics == nil {
		s.PodMetrics = new(Selector)
	}
	s.PodMetrics.Bind()

	if s.NodeMetrics == nil {
		s.NodeMetrics = new(Selector)
	}
	s.NodeMetrics.Bind()

	if s.Service == nil {
		s.Service = new(Selector)
	}
	s.Service.Bind()
	return s
}

// Bind binds the actual data from the Selector receiver field.
func (s *Selector) Bind() *Selector {
	// Initialization of s itself should be handled by the parent struct's Bind method.
	for k, v := range s.Labels {
		s.Labels[k] = GetActualValue(v)
	}
	for k, v := range s.Fields {
		s.Fields[k] = GetActualValue(v)
	}
	return s
}

func (r *ReadReplica) Bind() *ReadReplica {
	// Initialization of r itself should be handled by the parent struct's Bind method.
	r.IDKey = GetActualValue(r.IDKey)
	return r
}

// DiscovererClient represents the DiscovererClient configurations.
type DiscovererClient struct {
	Duration           string      `json:"duration"             yaml:"duration"`
	Client             *GRPCClient `json:"client"               yaml:"client"`
	AgentClientOptions *GRPCClient `json:"agent_client_options" yaml:"agent_client_options"`
}

// Bind binds the actual data from the DiscovererClient receiver field.
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

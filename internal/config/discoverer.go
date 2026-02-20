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

package config

// Discoverer represents the Discoverer configurations.
type Discoverer struct {
	// Net represents the network configuration.
	Net *Net `json:"net,omitempty" yaml:"net"`
	// Selectors represents the selector configuration.
	Selectors *Selectors `json:"selectors,omitempty" yaml:"selectors"`
	// Name represents the discoverer name.
	Name string `json:"name,omitempty" yaml:"name"`
	// Namespace represents the namespace.
	Namespace string `json:"namespace,omitempty" yaml:"namespace"`
	// DiscoveryDuration represents the discovery duration.
	DiscoveryDuration string `json:"discovery_duration,omitempty" yaml:"discovery_duration"`
}

// Selectors represents the selector configuration.
type Selectors struct {
	// Pod represents the pod selector.
	Pod *Selector `json:"pod,omitempty" yaml:"pod"`
	// Node represents the node selector.
	Node *Selector `json:"node,omitempty" yaml:"node"`
	// NodeMetrics represents the node metrics selector.
	NodeMetrics *Selector `json:"node_metrics,omitempty" yaml:"node_metrics"`
	// PodMetrics represents the pod metrics selector.
	PodMetrics *Selector `json:"pod_metrics,omitempty" yaml:"pod_metrics"`
	// Service represents the service selector.
	Service *Selector `json:"service,omitempty" yaml:"service"`
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

// Selector represents the selector configuration.
type Selector struct {
	// Labels represents the labels.
	Labels map[string]string `json:"labels,omitempty" yaml:"labels"`
	// Fields represents the fields.
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

// ReadReplica represents the read replica configuration.
type ReadReplica struct {
	// IDKey represents the ID key.
	IDKey string `json:"id_key,omitempty" yaml:"id_key"`
	// Enabled enables read replica.
	Enabled bool `json:"enabled,omitempty" yaml:"enabled"`
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
	// Client represents the client configuration.
	Client *GRPCClient `json:"client" yaml:"client"`
	// AgentClientOptions represents the agent client options.
	AgentClientOptions *GRPCClient `json:"agent_client_options" yaml:"agent_client_options"`
	// Duration represents the duration.
	Duration string `json:"duration" yaml:"duration"`
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

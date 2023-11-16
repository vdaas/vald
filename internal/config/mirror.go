// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package config

// Mirror represents the Mirror Gateway configuration.
type Mirror struct {
	// Net represents the network configuration tcp, udp, unix domain socket.
	Net *Net `json:"net,omitempty" yaml:"net"`

	// GRPCClient represents the configurations for gRPC client.
	Client *GRPCClient `json:"client" yaml:"client"`

	// SelfMirrorAddr represents the address for the self Mirror Gateway.
	SelfMirrorAddr string `json:"self_mirror_addr" yaml:"self_mirror_addr"`

	// GatewayAddr represents the address for the Vald Gateway (e.g lb-gateway).
	GatewayAddr string `json:"gateway_addr" yaml:"gateway_addr"`

	// PodName represents the mirror gateway pod name.
	PodName string `json:"pod_name" yaml:"pod_name"`

	// RegisterDuration represents the duration to register Mirror Gateway.
	RegisterDuration string `json:"register_duration" yaml:"register_duration"`

	// Namespace represents the target namespace to discover ValdMirrorTarget resource.
	Namespace string `json:"namespace" yaml:"namespace"`

	// DiscoveryDuration represents the duration to discover.
	DiscoveryDuration string `json:"discovery_duration" yaml:"discovery_duration"`

	// Colocation represents the colocation name.
	Colocation string `json:"colocation" yaml:"colocation"`

	// Group represents the group name of the Mirror Gateways.
	// It is used to discover ValdMirrorTarget resources with the same group name.
	Group string `json:"group" yaml:"group"`
}

// Bind binds the actual data from the Mirror receiver fields.
func (m *Mirror) Bind() *Mirror {
	m.SelfMirrorAddr = GetActualValue(m.SelfMirrorAddr)
	m.GatewayAddr = GetActualValue(m.GatewayAddr)
	m.PodName = GetActualValue(m.PodName)
	m.RegisterDuration = GetActualValue(m.RegisterDuration)
	m.Namespace = GetActualValue(m.Namespace)
	m.DiscoveryDuration = GetActualValue(m.DiscoveryDuration)
	m.Colocation = GetActualValue(m.Colocation)
	m.Group = GetActualValue(m.Group)

	if m.Net != nil {
		m.Net = m.Net.Bind()
	} else {
		m.Net = new(Net).Bind()
	}
	if m.Client != nil {
		m.Client = m.Client.Bind()
	} else {
		m.Client = new(GRPCClient).Bind()
	}
	return m
}

// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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
	Net *Net `json:"net,omitempty"                yaml:"net"`
	// Client represents the gRPC client configuration for connecting the LB Gateway.
	Client *GRPCClient `json:"client" yaml:"client"`
	// SelfMirrorAddr represents the address for the self Mirror Gateway.
	SelfMirrorAddr string `json:"self_mirror_addr" yaml:"self_mirror_addr"`
	// GatewayAddr represents the address for the Vald Gateway (e.g lb-gateway).
	GatewayAddr string `json:"gateway_addr" yaml:"gateway_addr"`
	// PodName represents self Mirror Gateway Pod name.
	PodName string `json:"pod_name" yaml:"pod_name"`
	// AdvertiseInterval represents interval to advertise Mirror Gateway information to other mirror gateway.
	AdvertiseInterval string `json:"advertise_interval" yaml:"advertise_interval"`
	// Namespace represents namespace to be reconciled.
	Namespace string `json:"namespace,omitempty"          yaml:"namespace"`
	// DiscoverDuration represents namespace to be reconciled.
	DiscoveryDuration string `json:"discovery_duration" yaml:"discovery_duration"`
	// Colocation represents colocation name.
	Colocation string `json:"colocation" yaml:"colocation"`
}

// Bind binds the actual data from the Mirror receiver fields.
func (m *Mirror) Bind() *Mirror {
	m.SelfMirrorAddr = GetActualValue(m.SelfMirrorAddr)
	m.GatewayAddr = GetActualValue(m.GatewayAddr)
	m.PodName = GetActualValue(m.PodName)
	m.AdvertiseInterval = GetActualValue(m.AdvertiseInterval)
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

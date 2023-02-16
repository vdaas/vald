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

// Package config providers configuration type and load configuration logic
package config

// Mirror represents the Mirror Gateway configuration.
type Mirror struct {
	// Client represents the gRPC client configuration for connecting the LB Gateway.
	Client *GRPCClient `json:"client" yaml:"client"`
	// MirrorClient represents the gRPC client configuration for connecting the self Mirror Gateway.
	SelfMirrorAddrs []string `json:"self_mirror_addrs" yaml:"self_mirror_addrs"`
	// MirrorClient represents the gRPC client configuration for connecting the Mirror Gateway.
	MirrorClient *GRPCClient `json:"mirror_client" yaml:"mirror_client"`
	// PodName represents self Mirror Gateway Pod name.
	PodName string `json:"pod_name" yaml:"pod_name"`
	// AdvertiseInterval represents interval to advertise Mirror Gateway information to other mirror gateway.
	AdvertiseInterval string `json:"advertise_interval" yaml:"advertise_interval"`
}

// Bind binds the actual data from the Mirror receiver fields.
func (m *Mirror) Bind() *Mirror {
	m.SelfMirrorAddrs = GetActualValues(m.SelfMirrorAddrs)
	m.PodName = GetActualValue(m.PodName)
	m.AdvertiseInterval = GetActualValue(m.AdvertiseInterval)
	if m.Client != nil {
		m.Client = m.Client.Bind()
	} else {
		m.Client = new(GRPCClient).Bind()
	}
	if m.MirrorClient != nil {
		m.MirrorClient = m.MirrorClient.Bind()
	} else {
		m.MirrorClient = new(GRPCClient).Bind()
	}
	return m
}

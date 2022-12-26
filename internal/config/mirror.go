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
	// MirrorClient represent the gRPC client configuration for connecting the LB Gateway.
	LB *GRPCClient `json:"lb" yaml:"lb"`
	// MirrorClient represent the gRPC client configuration for connecting the Mirror Gateway.
	Mirror *GRPCClient `json:"mirror" yaml:"mirror"`
	// MirrorClient represent the gRPC client configuration for connecting the self Mirror Gateway.
	SelfMirror *GRPCClient `json:"self_mirror" yaml:"self_mirror"`
}

// Bind binds the actual data from the Mirror receiver fields.
func (m *Mirror) Bind() *Mirror {
	if m.LB != nil {
		m.LB = m.LB.Bind()
	} else {
		m.LB = new(GRPCClient).Bind()
	}
	if m.Mirror != nil {
		m.Mirror = m.Mirror.Bind()
	} else {
		m.Mirror = new(GRPCClient).Bind()
	}
	return m
}

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

// SearchJob represents the configuration for the internal benchmark search job.
type SearchJob struct {
	Dimension     int         `json:"dimension" yaml: "dimension"`
	Iter          uint32      `json:"iter" yaml: "iter"`
	Num           uint32      `json:"num" yaml: "num"`
	MinNum        uint32      `json:"min_num" yaml:"min_num"`
	Radius        float64     `json:"radius" yaml:"radius"`
	Epsilon       float64     `json:"epsilon" yaml:"epsilon"`
	Timeout       string      `json:"timeout" yaml:"timeout"`
	GatewayClient *GRPCClient `json:"gateway_client" yaml:"gateway_client"`
}

// Bind binds the actual data from the Job search receiver fields.
func (s *SearchJob) Bind() *SearchJob {
	s.Timeout = GetActualValue(s.Timeout)

	if s.GatewayClient != nil {
		s.GatewayClient = s.GatewayClient.Bind()
	}
	return s
}

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

// CircuitBreaker represents the configuration for the internal circuitbreaker package.
type CircuitBreaker struct {
	ClosedErrorRate      float32 `yaml:"closed_error_rate"      json:"closed_error_rate,omitempty"`
	HalfOpenErrorRate    float32 `yaml:"half_open_error_rate"   json:"half_open_error_rate,omitempty"`
	MinSamples           int64   `yaml:"min_samples"            json:"min_samples,omitempty"`
	OpenTimeout          string  `yaml:"open_timeout"           json:"open_timeout,omitempty"`
	ClosedRefreshTimeout string  `yaml:"closed_refresh_timeout" json:"closed_refresh_timeout,omitempty"`
}

func (cb *CircuitBreaker) Bind() *CircuitBreaker {
	cb.OpenTimeout = GetActualValue(cb.OpenTimeout)
	cb.ClosedRefreshTimeout = GetActualValue(cb.ClosedRefreshTimeout)
	return cb
}

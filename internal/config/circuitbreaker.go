// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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
	// OpenTimeout represents the open timeout.
	OpenTimeout          string  `json:"open_timeout,omitempty"           yaml:"open_timeout"`
	// ClosedRefreshTimeout represents the closed refresh timeout.
	ClosedRefreshTimeout string  `json:"closed_refresh_timeout,omitempty" yaml:"closed_refresh_timeout"`
	// MinSamples represents the minimum samples.
	MinSamples           int64   `json:"min_samples,omitempty"            yaml:"min_samples"`
	// ClosedErrorRate represents the closed error rate.
	ClosedErrorRate      float32 `json:"closed_error_rate,omitempty"      yaml:"closed_error_rate"`
	// HalfOpenErrorRate represents the half-open error rate.
	HalfOpenErrorRate    float32 `json:"half_open_error_rate,omitempty"   yaml:"half_open_error_rate"`
}

func (cb *CircuitBreaker) Bind() *CircuitBreaker {
	cb.OpenTimeout = GetActualValue(cb.OpenTimeout)
	cb.ClosedRefreshTimeout = GetActualValue(cb.ClosedRefreshTimeout)
	return cb
}

// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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

package grpc

import "github.com/vdaas/vald/pkg/operator/benchmark/service"

type Option func(*server) error

var defaultOpts = []Option{}

// WithOperator sets the benchmark Operator that backs the handler's status
// queries. It must be supplied: the server promotes GetScenarioStatus and
// GetBenchmarkJobStatus from the embedded Operator, so without this the REST
// status endpoints would call methods on a nil interface and panic.
func WithOperator(op service.Operator) Option {
	return func(s *server) error {
		if op != nil {
			s.Operator = op
		}
		return nil
	}
}

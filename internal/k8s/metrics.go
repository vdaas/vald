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

package k8s

import metricsv1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"

// metrics.k8s.io/v1beta1 API types.
type (
	NodeMetrics     = metricsv1beta1.NodeMetrics
	NodeMetricsList = metricsv1beta1.NodeMetricsList
	PodMetrics      = metricsv1beta1.PodMetrics
	PodMetricsList  = metricsv1beta1.PodMetricsList
)

// AddMetricsToScheme registers the metrics.k8s.io/v1beta1 types with the given scheme.
//
//nolint:gochecknoglobals // immutable function alias confining k8s.io imports to this package
var AddMetricsToScheme = metricsv1beta1.AddToScheme

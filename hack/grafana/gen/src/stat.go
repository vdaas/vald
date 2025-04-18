//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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

package main

import (
	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/common"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/stat"
)

func addStatPanels(builder *dashboard.DashboardBuilder) {
	for _, panel := range []cog.Builder[dashboard.Panel]{
		statPanel("Vald Version", "vald_version"),
		statPanel("Go Version", "go_version"),
		statPanel("Go OS", "go_os"),
		stat.NewPanelBuilder().
			Title("Pods ($ReplicaSet)").
			WithTarget(prometheusQuery(
				`count(kube_pod_info{namespace="$Namespace", pod=~"$ReplicaSet.*"})`,
			).Format("table")).
			Span(widthQuarter).Height(heightMedium),
		stat.NewPanelBuilder().
			Title("Total memory working set ($ReplicaSet)").
			WithTarget(prometheusQuery(
				`sum(container_memory_working_set_bytes{namespace="$Namespace", container=~"$ReplicaSet", image!=""})`,
			).Format("time_series")).
			Unit("decbytes").
			Thresholds(
				dashboard.NewThresholdsConfigBuilder().
					Mode(dashboard.ThresholdsModeAbsolute).
					Steps(memThresholds),
			).
			Span(widthQuarter).Height(heightMedium),
		statPanel("Git Commit", "git_commit").Span(witdhOneThird),
		statPanel("Built at", "build_time"),
	} {
		builder.WithPanel(panel)
	}
}

func statPanel(title string, field string) *stat.PanelBuilder {
	return stat.NewPanelBuilder().
		Title(title).
		WithTarget(prometheusQuery(
			`app_version_info{exported_kubernetes_namespace="$Namespace", kubernetes_name=~"$ReplicaSet", target_pod=~"$PodName"}`,
		).Format("table")).
		ReduceOptions(common.NewReduceDataOptionsBuilder().Calcs([]string{"lastNotNull"}).Fields(field)).
		Span(widthOneSixth).Height(heightShort)
}

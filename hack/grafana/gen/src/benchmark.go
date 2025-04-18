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
	"fmt"

	"github.com/grafana/grafana-foundation-sdk/go/common"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/stat"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
)

func addSumStatPanel(builder *dashboard.DashboardBuilder, title string, metric string) {
	builder.
		WithPanel(
			stat.NewPanelBuilder().
				Title(title).
				WithTarget(prometheusQuery(
					fmt.Sprintf(
						`sum(%s{exported_kubernetes_namespace="$Namespace", kubernetes_name=~"$ReplicaSet", target_pod=~"$PodName"})`,
						metric,
					),
				).Format("table")).
				ReduceOptions(common.NewReduceDataOptionsBuilder().Calcs([]string{"lastNotNull"}).Fields("/^Value$/")).
				Span(widthOneSixth).Height(heightShort),
		)
}

func addBenchmarkStatPanel(builder *dashboard.DashboardBuilder, title string, field string, width uint32) {
	builder.
		WithPanel(
			stat.NewPanelBuilder().
				Title(title).
				WithTarget(prometheusQuery(
					`benchmark_operator_info{exported_kubernetes_namespace="$Namespace", kubernetes_name=~"$ReplicaSet", target_pod=~"$PodName"}`,
				).Format("table")).
				ReduceOptions(common.NewReduceDataOptionsBuilder().Calcs([]string{"lastNotNull"}).Fields(field)).
				Span(width).Height(heightShort),
		)
}

func addBenchmarkStatPanels(dashboard *dashboard.DashboardBuilder) {
	addSumStatPanel(dashboard, "All Scenario Count", "benchmark_operator_applied_scenario")
	addSumStatPanel(dashboard, "Running Scenario Count", "benchmark_operator_running_scenario")
	addSumStatPanel(dashboard, "Running Scenario Count", "benchmark_operator_complete_scenario")
	addBenchmarkStatPanel(dashboard, "Job Image Name", "repository", witdhOneThird)
	addBenchmarkStatPanel(dashboard, "Job Image Version", "tag", widthOneSixth)
	addSumStatPanel(dashboard, "All Benchmark Job Count", "benchmark_operator_applied_benchmark_job")
	addSumStatPanel(dashboard, "Running Benchmark Job Count", "benchmark_operator_running_benchmark_job")
	addSumStatPanel(dashboard, "Completed Benchmark Job Count", "benchmark_operator_complete_benchmark_job")
}

func addBenchmarkJobCPUPanel(builder *dashboard.DashboardBuilder) {
	panel := timeseries.NewPanelBuilder().
		Title("CPU").
		Height(heightTall).Span(widthHalf).
		WithTarget(prometheusQuery(
			`sum(irate(container_cpu_usage_seconds_total{namespace="$Namespace", pod=~"$JobPodName", image=~".*$JobReplicaSet.*"}[$interval])) by (pod) and on() count(kube_pod_created{pod=~"$JobPodName"}) >= 1`,
		).Format("time_series").LegendFormat("{{pod}}"))
	builder.WithPanel(panel)
}

func addBenchmarkJobMemoryPanel(builder *dashboard.DashboardBuilder) {
	panel := timeseries.NewPanelBuilder().
		Title("Memory Working Set").
		Height(heightTall).Span(widthHalf).
		WithTarget(prometheusQuery(
			`sum(container_memory_working_set_bytes{namespace="$Namespace", pod=~"$JobPodName", image=~".*$JobReplicaSet.*"}) by (pod) and on() count(kube_pod_created{pod=~"$JobPodName"}) >= 1`,
		).Format("time_series").LegendFormat("{{pod}}")).
		Unit("decbytes")
	builder.WithPanel(panel)
}

func addBenchmarkJobMetrics(builder *dashboard.DashboardBuilder) {
	addBenchmarkJobCPUPanel(builder)
	addBenchmarkJobMemoryPanel(builder)
	addCompletedRPCPanel(builder, "", "$JobReplicaSet", "$JobPodName")
	addLatencyPanel(builder, "", "$JobReplicaSet", "$JobPodName", ".*", `=~".*"`)
	addGoroutinePanel(builder, "$JobReplicaSet", "$JobPodName")
	addGCPanel(builder, "$JobReplicaSet")
}

func addOperatorMetrics(builder *dashboard.DashboardBuilder) {
	addCPUPanel(builder)
	addMemoryPanel(builder)
	addGoroutinePanel(builder, "$ReplicaSet", "$PodName")
	addGCPanel(builder, "$ReplicaSet")
}

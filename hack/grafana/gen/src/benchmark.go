//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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
	"github.com/grafana/grafana-foundation-sdk/go/common"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/stat"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
	"github.com/grafana/promql-builder/go/promql"
	"github.com/vdaas/vald/internal/observability/metrics/tools/benchmark"
	"github.com/vdaas/vald/pkg/tools/benchmark/operator/config"
)

func addSumStatPanel(builder *dashboard.DashboardBuilder, title string, metric string) {
	builder.
		WithPanel(
			stat.NewPanelBuilder().
				Title(title).
				WithTarget(prometheusQuery(
					promql.Sum(addBasicLabel(promql.Vector(metric))).String(),
				).Format("table")).
				ReduceOptions(common.NewReduceDataOptionsBuilder().Calcs([]string{"lastNotNull"}).Fields("/^Value$/")).
				Span(widthOneSixth).Height(heightShort),
		)
}

func addBenchmarkStatPanel(
	builder *dashboard.DashboardBuilder, title string, field string, width uint32,
) {
	builder.
		WithPanel(
			stat.NewPanelBuilder().
				Title(title).
				WithTarget(prometheusQuery(
					addBasicLabel(promql.Vector(config.BenchmarkOperatorInfo)).String(),
				).Format("table")).
				ReduceOptions(common.NewReduceDataOptionsBuilder().Calcs([]string{"lastNotNull"}).Fields(field)).
				Span(width).Height(heightShort),
		)
}

func addBenchmarkStatPanels(dashboard *dashboard.DashboardBuilder) {
	addSumStatPanel(dashboard, benchmark.AppliedScenarioCountDescription, benchmark.AppliedScenarioCount)
	addSumStatPanel(dashboard, benchmark.RunningScenarioCountDescription, benchmark.RunningScenarioCount)
	addSumStatPanel(dashboard, benchmark.CompleteScenarioCountDescription, benchmark.CompleteScenarioCount)
	addBenchmarkStatPanel(dashboard, "Job Image Name", "repository", witdhOneThird)
	addBenchmarkStatPanel(dashboard, "Job Image Version", "tag", widthOneSixth)
	addSumStatPanel(dashboard, benchmark.AppliedBenchmarkJobCountDescription, benchmark.AppliedBenchmarkJobCount)
	addSumStatPanel(dashboard, benchmark.RunningBenchmarkJobCountDescription, benchmark.RunningBenchmarkJobCount)
	addSumStatPanel(dashboard, benchmark.CompleteBenchmarkJobCountDescription, benchmark.CompleteBenchmarkJobCount)
}

func addBenchmarkJobCPUPanel(builder *dashboard.DashboardBuilder) {
	panel := timeseries.NewPanelBuilder().
		Title("CPU").
		Height(heightTall).Span(widthHalf).
		WithTarget(prometheusQuery(
			promql.Sum(promql.Irate(
				promql.Vector(cpuMetric).
					Label("namespace", namespaceVariable).
					LabelMatchRegexp("pod", "$JobPodName").
					LabelMatchRegexp(imageKey, ".*$JobReplicaSet.*").
					Range(intervalVariable),
			)).By([]string{"pod"}).String(),
		).Format("time_series").LegendFormat("{{pod}}"))
	builder.WithPanel(panel)
}

func addBenchmarkJobMemoryPanel(builder *dashboard.DashboardBuilder) {
	panel := timeseries.NewPanelBuilder().
		Title("Memory Working Set").
		Height(heightTall).Span(widthHalf).
		WithTarget(prometheusQuery(
			promql.Sum(
				promql.Vector(memMetric).
					Label("namespace", namespaceVariable).
					LabelMatchRegexp("pod", "$JobPodName").
					LabelMatchRegexp(imageKey, ".*$JobReplicaSet.*"),
			).By([]string{"pod"}).String(),
		).Format("time_series").LegendFormat("{{pod}}")).
		Unit("decbytes")
	builder.WithPanel(panel)
}

func addBenchmarkJobMetrics(builder *dashboard.DashboardBuilder) {
	addBenchmarkJobCPUPanel(builder)
	addBenchmarkJobMemoryPanel(builder)
	addCompletedRPCPanel(builder, "", "$JobReplicaSet", "$JobPodName")
	addLatencyPanel(builder, "", "$JobReplicaSet", "$JobPodName", ".*", ".*", true)
	addGoroutinePanel(builder, "$JobReplicaSet", "$JobPodName")
	addGCPanel(builder, "$JobReplicaSet")
}

func addOperatorMetrics(builder *dashboard.DashboardBuilder) {
	addCPUPanel(builder)
	addMemoryPanel(builder)
	addGoroutinePanel(builder, "$ReplicaSet", "$PodName")
	addGCPanel(builder, "$ReplicaSet")
}

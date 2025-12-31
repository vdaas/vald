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
	"github.com/grafana/promql-builder/go/promql"
	"github.com/vdaas/vald/internal/observability/metrics/agent/core/ngt"
	"github.com/vdaas/vald/pkg/agent/core/ngt/config"
)

func addAgentPanels(builder *dashboard.DashboardBuilder) {
	builder.
		WithPanel(stat.NewPanelBuilder().
			Title("Indices").
			WithTarget(prometheusQuery(
				promql.Sum(
					addBasicLabel(promql.Vector(ngt.IndexCountMetricsName)),
				).String(),
			).Format("table")).
			Thresholds(
				dashboard.NewThresholdsConfigBuilder().
					Mode(dashboard.ThresholdsModeAbsolute).
					Steps(indiceThresholds),
			).
			Span(witdhOneThird).Height(heightMedium)).
		WithPanel(stat.NewPanelBuilder().
			Title("Indexing Pods").
			WithTarget(prometheusQuery(
				promql.Sum(
					addBasicLabel(promql.Vector(ngt.IsIndexingMetricsName)),
				).String(),
			).Format("table")).
			Thresholds(
				dashboard.NewThresholdsConfigBuilder().
					Mode(dashboard.ThresholdsModeAbsolute).
					Steps(podThresholds),
			).
			Span(widthOneSixth).Height(heightMedium)).
		WithPanel(stat.NewPanelBuilder().
			Title("Uncommitted Indices").
			WithTarget(prometheusQuery(
				promql.Sum(
					addBasicLabel(promql.Vector(ngt.UncommittedIndexCountMetricsName)),
				).String(),
			).Format("table")).
			Thresholds(
				dashboard.NewThresholdsConfigBuilder().
					Mode(dashboard.ThresholdsModeAbsolute).
					Steps(queueThresholds),
			).
			Span(widthOneSixth).Height(heightMedium)).
		WithPanel(stat.NewPanelBuilder().
			Title("Insert Vqueue").
			WithTarget(prometheusQuery(
				promql.Sum(
					addBasicLabel(promql.Vector(ngt.InsertVQueueCountMetricsName)),
				).String(),
			).Format("table")).
			Thresholds(
				dashboard.NewThresholdsConfigBuilder().
					Mode(dashboard.ThresholdsModeAbsolute).
					Steps(queueThresholds),
			).
			Span(widthOneSixth).Height(heightMedium)).
		WithPanel(stat.NewPanelBuilder().
			Title("Delete Vqueue").
			WithTarget(prometheusQuery(
				promql.Sum(
					addBasicLabel(promql.Vector(ngt.DeleteVQueueCountMetricsName)),
				).String(),
			).Format("table")).
			Thresholds(
				dashboard.NewThresholdsConfigBuilder().
					Mode(dashboard.ThresholdsModeAbsolute).
					Steps(queueThresholds),
			).
			Span(widthOneSixth).Height(heightMedium)).
		WithPanel(statPanel("Algorithm Version", "algorithm_info").Span(widthQuarter)).
		WithPanel(agentStatPanel("Dimension", "dimension")).
		WithPanel(agentStatPanel("Distance Type", "distance_type")).
		WithPanel(agentStatPanel("Object Type", "object_type").Span(widthOneEighth)).
		WithPanel(stat.NewPanelBuilder().
			Title("Broken Index Store Count").
			WithTarget(prometheusQuery(
				promql.Sum(
					addBasicLabel(promql.Vector(ngt.BrokenIndexStoreCountMetricsName)),
				).String(),
			).Format("table")).
			Span(widthOneEighth).Height(heightShort))
}

func agentStatPanel(title string, field string) *stat.PanelBuilder {
	return stat.NewPanelBuilder().
		Title(title).
		WithTarget(prometheusQuery(
			addBasicLabel(promql.Vector(config.AgentNGTInfo)).String(),
		).Format("table")).
		ReduceOptions(common.NewReduceDataOptionsBuilder().Calcs([]string{"lastNotNull"}).Fields(field)).
		Span(widthQuarter).Height(heightShort)
}

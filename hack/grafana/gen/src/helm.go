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
	"fmt"

	"github.com/grafana/grafana-foundation-sdk/go/common"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/stat"
	"github.com/grafana/grafana-foundation-sdk/go/table"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
	"github.com/grafana/promql-builder/go/promql"
)

func addReconcileDurationPanel(dashboard *dashboard.DashboardBuilder) {
	panel := timeseries.NewPanelBuilder().
		Title("Reconcile Duration").
		Span(widthHalf).Height(heightTall)
	for _, quantile := range quntiles {
		panel.WithTarget(prometheusQuery(
			promql.HistogramQuantile(
				quantile,
				promql.Sum(promql.Rate(promql.Vector(
					reconcileSec,
				).Label("namespace", namespaceVariable).
					LabelMatchRegexp(appNameKey, nameVariable).
					LabelMatchRegexp(instanceKey, podVariable).
					Range(intervalVariable),
				)).By([]string{"le"}),
			).String(),
			// fmt.Sprintf(
			// 	`histogram_quantile(%f, sum(rate(controller_runtime_reconcile_time_seconds_bucket{namespace=~"$Namespace", app_kubernetes_io_name=~"$ReplicaSet", instance=~"$PodName"}[$interval])) by (le))`,
			// 	quantile,
			// ),
		).
			Format("time_series").LegendFormat(fmt.Sprintf("p%d", int(quantile*100))))
	}
	dashboard.WithPanel(panel)
}

func addK8SAPILantencyPanel(dashboard *dashboard.DashboardBuilder) {
	panel := timeseries.NewPanelBuilder().
		Title("K8S Server API Latency").
		Span(widthHalf).Height(heightTall)
	for _, quantile := range quntiles {
		panel.WithTarget(prometheusQuery(
			promql.HistogramQuantile(
				quantile,
				promql.Sum(promql.Rate(promql.Vector(
					restDurationSec,
				).Label("namespace", namespaceVariable).
					Range(intervalVariable),
				)).By([]string{"le"}),
			).String(),
		).
			Format("time_series").LegendFormat(fmt.Sprintf("p%d", int(quantile*100))))
	}
	dashboard.WithPanel(panel)
}

func addHelmOperatorPanels(builder *dashboard.DashboardBuilder) {
	builder.
		WithPanel(
			stat.NewPanelBuilder().
				Title("Operator SDK Version").
				WithTarget(prometheusQuery(
					promql.Vector(helmInfo).
						Label("namespace", namespaceVariable).
						Label(appNameKey, nameVariable).
						LabelMatchRegexp(instanceKey, podVariable).String(),
				).Format("table")).
				ReduceOptions(common.NewReduceDataOptionsBuilder().Calcs([]string{"lastNotNull"}).Fields("version")).
				Span(widthQuarter).Height(heightMedium),
		).
		WithPanel(
			stat.NewPanelBuilder().
				Title("Go Version").
				WithTarget(prometheusQuery(
					promql.Vector(goInfo).
						Label("namespace", namespaceVariable).
						LabelMatchRegexp(appNameKey, nameVariable).
						LabelMatchRegexp(instanceKey, podVariable).String(),
				).Format("table")).
				ReduceOptions(common.NewReduceDataOptionsBuilder().Calcs([]string{"lastNotNull"}).Fields("version")).
				Span(widthQuarter).Height(heightMedium),
		).
		WithPanel(
			table.NewPanelBuilder().
				Title("Resources").
				WithTarget(prometheusQuery(
					promql.Vector(resCreatedSec).
						Label("namespace", namespaceVariable).
						LabelMatchRegexp(appNameKey, nameVariable).
						LabelMatchRegexp(instanceKey, podVariable).String(),
				).Format("table")).
				Span(widthHalf).Height(heightMedium),
		).
		WithPanel(
			timeseries.NewPanelBuilder().
				Title("Reconcile Total").
				Span(widthHalf).Height(heightTall).
				WithTarget(prometheusQuery(
					promql.Sum(promql.Irate(promql.Vector(reconcileTotal).
						Label("namespace", namespaceVariable).
						LabelMatchRegexp(appNameKey, nameVariable).
						LabelMatchRegexp(instanceKey, podVariable).
						Range(intervalVariable),
					)).By([]string{"controller", "result"}).String(),
				).Format("time_series").LegendFormat("{{controller}} ({{result}})")),
		).
		WithPanel(
			timeseries.NewPanelBuilder().
				Title("Reconcile Errors").
				Span(widthHalf).Height(heightTall).
				WithTarget(prometheusQuery(
					promql.Sum(promql.Irate(promql.Vector(reconcileErrorTotal).
						Label("namespace", namespaceVariable).
						LabelMatchRegexp(appNameKey, nameVariable).
						LabelMatchRegexp(instanceKey, podVariable).
						Range(intervalVariable),
					)).By([]string{"controller"}).String(),
				).Format("time_series").LegendFormat("{{controller}}")),
		)
	addReconcileDurationPanel(builder)
	addK8SAPILantencyPanel(builder)
}

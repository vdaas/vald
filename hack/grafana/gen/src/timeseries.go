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

	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
	"github.com/grafana/promql-builder/go/promql"
	"github.com/vdaas/vald/internal/observability/metrics/agent/core/ngt"
	"github.com/vdaas/vald/internal/observability/metrics/mem"
	"github.com/vdaas/vald/internal/observability/metrics/mem/index"
	"github.com/vdaas/vald/internal/observability/metrics/runtime/goroutine"
)

func addCPUPanel(dashboard *dashboard.DashboardBuilder) {
	panel := timeseries.NewPanelBuilder().
		Title("CPU").
		Span(widthHalf).Height(heightTall).
		WithTarget(prometheusQuery(
			promql.Sum(promql.Irate(promql.Vector(cpuMetric).
				Label("namespace", namespaceVariable).
				LabelMatchRegexp("container", "$ReplicaSet").
				LabelMatchRegexp("pod", "$PodName").
				LabelNeq("image", "").
				Range(intervalVariable))).By([]string{"pod"}).String(),
		).Format("time_series").LegendFormat("{{pod}}"))
	dashboard.WithPanel(panel.Min(0))
}

func addMemoryPanel(builder *dashboard.DashboardBuilder) {
	panel := timeseries.NewPanelBuilder().
		Title("Memory").
		Span(widthHalf).Height(heightTall)
	panel.WithTarget(prometheusQuery(
		promql.Sum(promql.Vector(memMetric).
			Label("namespace", namespaceVariable).
			Label("container", nameVariable).
			LabelMatchRegexp("pod", podVariable).
			LabelNeq("image", "")).By([]string{"pod"}).String()).
		Format("time_series").LegendFormat("{{target_pod}}"))
	builder.WithPanel(panel.Min(0))
}

func addJobCPUPanel(dashboard *dashboard.DashboardBuilder) {
	panel := timeseries.NewPanelBuilder().
		Title("CPU").
		Span(widthHalf).Height(heightTall).
		WithTarget(prometheusQuery(
			promql.Sum(promql.Irate(promql.Vector(cpuMetric).
				Label("namespace", namespaceVariable).
				Label("container", nameVariable).
				LabelMatchRegexp("pod", podVariable).
				LabelNeq("image", "").
				Range(intervalVariable))).By([]string{"pod"}).String()).
			Format("time_series").LegendFormat("{{pod}}"))
	dashboard.WithPanel(panel.Min(0))
}

func addJobMemoryPanel(builder *dashboard.DashboardBuilder) {
	panel := timeseries.NewPanelBuilder().
		Title("Memory").
		Span(widthHalf).Height(heightTall).
		WithTarget(prometheusQuery(
			promql.Sum(promql.Vector(memMetric).
				Label("namespace", namespaceVariable).
				Label("container", nameVariable).
				LabelMatchRegexp("pod", podVariable).
				LabelNeq("image", "")).By([]string{"pod"}).String()).
			Format("time_series").LegendFormat("{{pod}}"))
	builder.WithPanel(panel.Min(0))
}

func addLatencyPanel(
	builder *dashboard.DashboardBuilder,
	subTitle string,
	kubernetesName string,
	targetPod string,
	container string,
	method string,
	match bool,
) {
	panel := timeseries.NewPanelBuilder().
		Title(fmt.Sprintf("Latency (%s%s)", subTitle, targetPod)).
		Span(widthHalf).Height(heightTall)
	for _, quantile := range quntiles {
		panel.WithTarget(prometheusQuery(
			promql.HistogramQuantile(quantile, promql.Sum(promql.Rate(
				addGRPCMatch(
					addBasicLabel(promql.Vector(serverLatencyBucket)).
						LabelMatchRegexp("container", container),
					method, match,
				).Range(intervalVariable))).By([]string{"le", grpcServerMethod})).String(),
		).
			Format("time_series").LegendFormat(fmt.Sprintf("{{%s}} p%d", grpcServerMethod, int(quantile*100)))).Min(0)
	}
	builder.WithPanel(panel)
}

func addCompletedRPCPanel(
	dashboard *dashboard.DashboardBuilder, subTitle string, kubernetesName string, targetPod string,
) {
	panel := timeseries.NewPanelBuilder().
		Title(fmt.Sprintf("Completed RPCs (%s%s)", subTitle, targetPod)).
		Span(widthHalf).Height(heightTall).
		WithTarget(prometheusQuery(
			promql.Sum(promql.Irate(
				addBasicLabel(promql.Vector(serverCompletedRPCs)).
					Range(intervalVariable))).By([]string{grpcServerMethod, grpcServerStatus}).String(),
		).Format("time_series").LegendFormat("{{grpc_server_method}} ({{grpc_server_status}})")).
		WithTarget(prometheusQuery(
			promql.Sum(promql.Irate(
				addBasicLabel(promql.Vector(serverCompletedRPCs)).
					Range(intervalVariable))).By([]string{grpcServerStatus}).String(),
		).Format("time_series").LegendFormat("Total ({{grpc_server_status}})"))
	dashboard.WithPanel(panel)
}

func addGoroutinePanel(
	dashboard *dashboard.DashboardBuilder, kubernetesName string, targetPod string,
) {
	panel := timeseries.NewPanelBuilder().
		Title(goroutine.MetricsDescription).
		Span(widthHalf).Height(heightTall).
		WithTarget(prometheusQuery(
			addBasicLabel(promql.Vector(goroutine.MetricsName)).String(),
		).Format("time_series").LegendFormat("{{target_pod}}"))
	dashboard.WithPanel(panel)
}

func addGCPanel(dashboard *dashboard.DashboardBuilder, kubernetesName string) {
	panel := timeseries.NewPanelBuilder().
		Title(mem.NumGCMetricsDescription).
		Span(widthHalf).Height(heightTall).
		WithTarget(prometheusQuery(
			promql.Increase(promql.Vector(mem.NumGCMetricsName).
				Label(namespaceKey, namespaceVariable).
				LabelMatchRegexp(nameKey, kubernetesName).
				LabelMatchRegexp("target_node", ".+").
				Range(intervalVariable)).String(),
		).Format("time_series").LegendFormat("{{target_pod}}"))
	dashboard.WithPanel(panel)
}

func addIndexPanel(dashboard *dashboard.DashboardBuilder) {
	panel := timeseries.NewPanelBuilder().
		Title("Total Indices").
		Span(widthHalf).Height(heightTall).
		WithTarget(prometheusQuery(
			promql.Sum(addBasicLabel(
				promql.Vector(ngt.IndexCountMetricsName))).String(),
		).Format("time_series").LegendFormat("indices"))
	dashboard.WithPanel(panel)
}

func addUncommitedIndexPanel(dashboard *dashboard.DashboardBuilder) {
	panel := timeseries.NewPanelBuilder().
		Title("Uncommitted Indices").
		Span(widthHalf).Height(heightTall).
		WithTarget(prometheusQuery(
			promql.Sum(addBasicLabel(promql.Vector(ngt.UncommittedIndexCountMetricsName))).String(),
		).Format("time_series").LegendFormat("uncommitted-indices"))
	dashboard.WithPanel(panel)
}

func addIndexLatencyPanel(dashboard *dashboard.DashboardBuilder) {
	panel := timeseries.NewPanelBuilder().
		Title("SaveIndex Latency").
		Span(widthHalf).Height(heightTall)
	for _, quantile := range quntiles {
		panel.WithTarget(prometheusQuery(
			promql.HistogramQuantile(quantile, promql.Sum(promql.Rate(addBasicLabel(
				promql.Vector(serverLatencyBucket),
			).LabelMatchRegexp(grpcServerMethod, ".*Index$").
				Range(intervalVariable))).By([]string{"le", grpcServerMethod})).String(),
		).
			Format("time_series").LegendFormat(fmt.Sprintf("{{%s}} p%d", grpcServerMethod, int(quantile*100))))
	}
	dashboard.WithPanel(panel)
}

func addIndexPerPodPanel(dashboard *dashboard.DashboardBuilder) {
	panel := timeseries.NewPanelBuilder().
		Title("Indices Per Pod").
		Span(widthHalf).Height(heightTall).
		WithTarget(prometheusQuery(
			promql.Sum(addBasicLabel(promql.Vector(
				ngt.IndexCountMetricsName))).
				By([]string{"target_pod"}).String(),
		).Format("time_series").LegendFormat("{{hostname}}"))
	dashboard.WithPanel(panel)
}

func addMemstatsPanels(dashboard *dashboard.DashboardBuilder) {
	addMetricPanel(dashboard, index.AllocMetricsDescription, index.AllocMetricsName, cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, index.TotalAllocMetricsDescription, index.TotalAllocMetricsName, cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, index.SysMetricsDescription, index.SysMetricsName, cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, mem.LookupsMetricsDescription, mem.LookupsMetricsName, nil)
	addMetricPanel(dashboard, mem.MallocsMetricsDescription, mem.MallocsMetricsName, nil)
	addMetricPanel(dashboard, mem.FreesMetricsDescription, mem.FreesMetricsName, nil)
	addMetricPanel(dashboard, mem.HeapAllocMetricsDescription, mem.HeapAllocMetricsName, cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, mem.HeapSysMetricsDescription, mem.HeapSysMetricsName, cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, mem.HeapIdleMetricsDescription, mem.HeapIdleMetricsName, cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, mem.HeapInuseMetricsDescription, mem.HeapInuseMetricsName, cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, mem.HeapReleasedMetricsDescription, mem.HeapReleasedMetricsName, cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, mem.HeapObjectsMetricsDescription, mem.HeapObjectsMetricsName, nil)
	addMetricPanel(dashboard, mem.StackInuseMetricsDescription, mem.StackInuseMetricsName, cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, mem.StackSysMetricsDescription, mem.StackSysMetricsName, cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, mem.MspanInuseMetricsDescription, mem.MspanInuseMetricsName, cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, mem.MspanSysMetricsDescription, mem.MspanSysMetricsName, cog.ToPtr("decbytes"))

	addMetricPanel(dashboard, mem.McacheInuseMetricsDescription, mem.McacheInuseMetricsName, cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, mem.McacheSysMetricsDescription, mem.McacheSysMetricsName, cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, mem.BuckHashSysMetricsDescription, mem.BuckHashSysMetricsName, cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, mem.GcSysMetricsDescription, mem.GcSysMetricsName, cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, mem.OtherSysMetricsDescription, mem.OtherSysMetricsName, cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, mem.NextGcSysMetricsDescription, mem.NextGcSysMetricsName, cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, mem.PauseTotalMsMetricsDescription, mem.PauseTotalMsMetricsName, cog.ToPtr("ms"))
	addMetricPanel(dashboard, mem.NumGCMetricsDescription, mem.NumGCMetricsName, nil)
	addMetricPanel(dashboard, mem.NumForcedGCMetricsDescription, mem.NumGCMetricsName, nil)
	addMetricPanel(dashboard, mem.HeapWillReturnMetricsDescription, mem.HeapWillReturnMetricsName, cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, mem.LiveObjectsMetricsDescription, mem.LiveObjectsMetricsName, nil)
}

func addProcStatusPanels(dashboard *dashboard.DashboardBuilder) {
	addMetricPanel(dashboard, mem.VmpeakMetricsDescription, mem.VmpeakMetricsName, cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, mem.VmsizeMetricsDescription, mem.VmsizeMetricsName, cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, mem.VmdataMetricsDescription, mem.VmdataMetricsName, cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, mem.VmrssMetricsDescription, mem.VmrssMetricsName, cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, mem.VmhwmMetricsDescription, mem.VmhwmMetricsName, cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, mem.VmstkMetricsDescription, mem.VmstkMetricsName, cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, mem.VmswapMetricsDescription, mem.VmswapMetricsName, cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, mem.VmexeMetricsDescription, mem.VmexeMetricsName, cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, mem.VmlibMetricsDescription, mem.VmlibMetricsName, cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, mem.VmlckMetricsDescription, mem.VmlckMetricsName, cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, mem.VmpinMetricsDescription, mem.VmpinMetricsName, cog.ToPtr("decbytes"))
	addMetricPanel(dashboard, mem.VmpteMetricsDescription, mem.VmpteMetricsName, cog.ToPtr("decbytes"))
}

func addMetricPanel(
	dashboard *dashboard.DashboardBuilder, title string, metric string, unit *string,
) {
	panel := timeseries.NewPanelBuilder().
		Title(title).
		Span(widthQuarter).Height(heightTall).
		WithTarget(prometheusQuery(
			addBasicLabel(promql.Vector(metric)).String(),
		).Format("time_series").LegendFormat("{{target_pod}}")).Min(0)
	if unit != nil {
		panel.Unit(*unit)
	}
	dashboard.WithPanel(panel)
}

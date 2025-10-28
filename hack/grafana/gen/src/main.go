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
package main

import (
	"bytes"
	"path/filepath"
	"time"

	"github.com/grafana/grafana-foundation-sdk/go/cog/plugins"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"
	"github.com/vdaas/vald/internal/encoding/json"
	"github.com/vdaas/vald/internal/os"
	"github.com/vdaas/vald/internal/strings"
	"gopkg.in/yaml.v3"
)

func prometheusQuery(query string) *prometheus.DataqueryBuilder {
	return prometheus.NewDataqueryBuilder().
		Expr(query)
}

func ValdClusterOverview() *dashboard.DashboardBuilder {
	builder := dashboard.NewDashboardBuilder("Vald Cluster Overview").Uid("00-vald-cluster-overview")
	addOverviewVariables(builder)
	addStatPanels(builder)
	addOverviewIndexPanel(builder)
	addNodeCPUPanel(builder)
	addNodeMemoryPanel(builder)

	addCompletedRPCPanel(builder, "Vald LB Gateway: ", "$ReplicaSet", "$ValdGatewayPodName")
	addLatencyPanel(builder, "Vald LB Gateway: ", "$ReplicaSet", "$ValdGatewayPodName", "$ValdGatewayContainerName", ".*", true)
	addCompletedRPCPanel(builder, "Vald Agent: ", "$ReplicaSet", "$PodName")
	addLatencyPanel(builder, "Vald Agent: ", "$ReplicaSet", "$PodName", "$ContainerName", ".*Index$", false)
	addBackoffPanel(builder)
	addBackoffPerRPCPanel(builder)
	addCircuitBreakerState(builder)
	addLatencyPanel(builder, "Vald Agent Index: ", "$ReplicaSet", "$PodName", "$ContainerName", ".*Index$", true)
	repeatOverview(builder)
	return builder.Time("now-5m", "now")
}

func ValdLBGateway() *dashboard.DashboardBuilder {
	dashboard := dashboard.NewDashboardBuilder("Vald LB Gateway").Uid("08-vald-lb-gateway")
	addVariables(dashboard, "gateway lb.*", "1m")
	addStatPanels(dashboard)
	addCPUPanel(dashboard)
	addMemoryPanel(dashboard)
	addCompletedRPCPanel(dashboard, "", "$ReplicaSet", "$PodName")
	addLatencyPanel(dashboard, "", "$ReplicaSet", "$PodName", "$ContainerName", ".*", true)
	addGoroutinePanel(dashboard, "$ReplicaSet", "$PodName")
	addGCPanel(dashboard, "$ReplicaSet")
	return dashboard.Time("now-3h", "now")
}

func ValdDiscoverer() *dashboard.DashboardBuilder {
	dashboard := dashboard.NewDashboardBuilder("Vald Discoverer").Uid("02-vald-discoverer")
	addVariables(dashboard, "discoverer.*", "5m")
	addStatPanels(dashboard)
	addCPUPanel(dashboard)
	addMemoryPanel(dashboard)
	addCompletedRPCPanel(dashboard, "", "$ReplicaSet", "$PodName")
	addLatencyPanel(dashboard, "", "$ReplicaSet", "$PodName", "$ContainerName", ".*", true)
	addGoroutinePanel(dashboard, "$ReplicaSet", "$PodName")
	addGCPanel(dashboard, "$ReplicaSet")
	return dashboard.Time("now-3h", "now")
}

func ValdIndexManager() *dashboard.DashboardBuilder {
	dashboard := dashboard.NewDashboardBuilder("Vald Index Manager").Uid("05-vald-index-manager")
	addVariables(dashboard, "index manager", "5m")
	addStatPanels(dashboard)
	addCPUPanel(dashboard)
	addMemoryPanel(dashboard)
	addCompletedRPCPanel(dashboard, "", "$ReplicaSet", "$PodName")
	addLatencyPanel(dashboard, "", "$ReplicaSet", "$PodName", "$ContainerName", ".*", true)
	addGoroutinePanel(dashboard, "$ReplicaSet", "$PodName")
	addGCPanel(dashboard, "$ReplicaSet")
	return dashboard.Time("now-3h", "now")
}

func ValdIndexCorrection() *dashboard.DashboardBuilder {
	builder := dashboard.NewDashboardBuilder("Vald Index Correction").Uid("09-vald-index-correction")
	addVariables(builder, "index correction job", "5m")
	builder.
		WithRow(dashboard.NewRowBuilder("Stats"))
	addStatPanels(builder)
	builder.
		WithRow(dashboard.NewRowBuilder("Index Correction Stats"))
	addIndexCorrectionPanels(builder)
	builder.
		WithRow(dashboard.NewRowBuilder("Time Series"))
	addJobCPUPanel(builder)
	addJobMemoryPanel(builder)
	addGoroutinePanel(builder, "$ReplicaSet", "$PodName")
	addGCPanel(builder, "$ReplicaSet")
	return builder.Time("now-3h", "now")
}

func ValdAgent() *dashboard.DashboardBuilder {
	builder := dashboard.NewDashboardBuilder("Vald Agent").Uid("01-vald-agent")
	builder.
		WithRow(dashboard.NewRowBuilder("Stats"))
	addVariables(builder, "agent .*", "5m")
	addStatPanels(builder)
	builder.
		WithRow(dashboard.NewRowBuilder("Agent Stats"))
	addAgentPanels(builder)
	builder.
		WithRow(dashboard.NewRowBuilder("Time Series"))
	addCPUPanel(builder)
	addMemoryPanel(builder)
	addIndexPanel(builder)
	addUncommitedIndexPanel(builder)
	addIndexPerPodPanel(builder)
	addIndexLatencyPanel(builder)
	addCompletedRPCPanel(builder, "", "$ReplicaSet", "$PodName")
	addLatencyPanel(builder, "", "$ReplicaSet", "$PodName", "$ContainerName", ".*SaveIndex", false)
	addGCPanel(builder, "$ReplicaSet")
	addGoroutinePanel(builder, "$ReplicaSet", "$PodName")
	return builder.Time("now-3h", "now")
}

func ValdAgentMemory() *dashboard.DashboardBuilder {
	builder := dashboard.NewDashboardBuilder("Vald Agent Memory").Uid("99-vald-agent-memory")
	builder.
		WithRow(dashboard.NewRowBuilder("runtime.Memstats"))
	addVariables(builder, "agent .*", "5m")
	addMemstatsPanels(builder)
	builder.
		WithRow(dashboard.NewRowBuilder("/proc/<pid>/status"))
	addProcStatusPanels(builder)
	return builder.Time("now-15m", "now").Refresh("5s")
}

func ValdHelmOperator() *dashboard.DashboardBuilder {
	builder := dashboard.NewDashboardBuilder("Vald Helm Operator").Uid("07-vald-helm-operator")
	addHelmOperatorVariables(builder)
	addHelmOperatorPanels(builder)
	return builder.Time("now-5m", "now")
}

func ValdBenchmarkOperator() *dashboard.DashboardBuilder {
	builder := dashboard.NewDashboardBuilder("Vald Benchmark Operator").Uid("10-vald-benchmark-operator")
	addBenchmarkOperatorVariables(builder)
	addStatPanels(builder)
	addBenchmarkStatPanels(builder)
	builder.
		WithRow(dashboard.NewRowBuilder("Benchmark Job Metrics"))
	addBenchmarkJobMetrics(builder)
	builder.
		WithRow(dashboard.NewRowBuilder("Operator Metrics"))
	addOperatorMetrics(builder)
	return builder.Time("now-3h", "now")
}

func ValdOtelCollector() *dashboard.DashboardBuilder {
	builder := dashboard.NewDashboardBuilder("Vald Otel Collector").Uid("11-vald-otel-collector")
	addVariables(builder, "opentelemetry-collector.*", "5m")

	builder.WithRow(dashboard.NewRowBuilder("Stats"))
	addStatPanels(builder)

	builder.WithRow(dashboard.NewRowBuilder("Time Series"))
	addCPUPanel(builder)
	addMemoryPanel(builder)
	addOtelQueuePanel(builder)
	addOtelSpanDropPanel(builder)
	addOtelMetricDropPanel(builder)
	addOtelNetworkPanel(builder)

	return builder.Time("now-3h", "now")
}

func configMapTemplate(header string, id string, content string) *yaml.Node {
	return &yaml.Node{
		HeadComment: header,
		Kind:        yaml.MappingNode,
		Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Value: "apiVersion"},
			{Kind: yaml.ScalarNode, Value: "v1"},

			{Kind: yaml.ScalarNode, Value: "kind"},
			{Kind: yaml.ScalarNode, Value: "ConfigMap"},

			{Kind: yaml.ScalarNode, Value: "metadata"},
			{
				Kind: yaml.MappingNode,
				Content: []*yaml.Node{
					{Kind: yaml.ScalarNode, Value: "name"},
					{Kind: yaml.ScalarNode, Value: "grafana-dashboards-" + strings.Join(strings.Split(id, "-")[1:], "-")},
				},
			},

			{Kind: yaml.ScalarNode, Value: "data"},
			{
				Kind: yaml.MappingNode,
				Content: []*yaml.Node{
					{Kind: yaml.ScalarNode, Value: id + ".json"},
					{Kind: yaml.ScalarNode, Value: content},
				},
			},
		},
	}
}

func main() {
	maintainer := os.Getenv(maintainerKey)
	if maintainer == "" {
		maintainer = defaultMaintainer
	}
	rootDir := os.Getenv(rootKey)
	var header bytes.Buffer
	err := license.Execute(&header, struct {
		Year       int
		Maintainer string
	}{
		Year:       time.Now().Year(),
		Maintainer: maintainer,
	})
	if err != nil {
		panic(err)
	}

	// Required to correctly unmarshal panels and dataqueries
	plugins.RegisterDefaultPlugins()

	for _, builder := range []func() *dashboard.DashboardBuilder{ValdClusterOverview, ValdLBGateway, ValdDiscoverer, ValdIndexManager, ValdAgent, ValdIndexCorrection, ValdAgentMemory, ValdHelmOperator, ValdBenchmarkOperator, ValdOtelCollector} {
		dashboardModel, err := builder().Build()
		if err != nil {
			panic(err)
		}

		ArrangePanels(&dashboardModel)

		dashboardJson, err := json.MarshalIndent(dashboardModel, "", "  ")
		if err != nil {
			panic(err)
		}

		filePath := filepath.Join(rootDir, "k8s/metrics/grafana/dashboards", *dashboardModel.Uid+".yaml")
		file, err := os.Create(filePath)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		enc := yaml.NewEncoder(file)
		enc.SetIndent(2)
		err = enc.Encode(configMapTemplate(header.String(), *dashboardModel.Uid, string(dashboardJson)))
		if err != nil {
			panic(err)
		}
	}
}

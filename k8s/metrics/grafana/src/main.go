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
	"encoding/json"
	"net/url"

	"github.com/go-openapi/strfmt"
	"github.com/grafana/grafana-foundation-sdk/go/cog/plugins"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"
	goapi "github.com/grafana/grafana-openapi-client-go/client"
	"github.com/grafana/grafana-openapi-client-go/models"
)

func prometheusQuery(query string) *prometheus.DataqueryBuilder {
	return prometheus.NewDataqueryBuilder().
		Expr(query)
}

func ValdLBGateway() *dashboard.DashboardBuilder {
	dashboard := dashboard.NewDashboardBuilder("Vald LB Gateway").Uid("vald_lb_gateway")
	addVariables(dashboard, "gateway lb.*")
	addStatPanels(dashboard)
	addMemoryPanel(dashboard)
	addCPUPanel(dashboard)
	addCompletedRPCPanel(dashboard)
	addLatencyPanel(dashboard)
	addGCPanel(dashboard)
	addGoroutinePanel(dashboard)
	return dashboard.Time("now-3h", "now")
}

func ValdDiscoverer() *dashboard.DashboardBuilder {
	dashboard := dashboard.NewDashboardBuilder("Vald Discoverer").Uid("vald_discoverer")
	addVariables(dashboard, "discoverer.*")
	addStatPanels(dashboard)
	addMemoryPanel(dashboard)
	addCPUPanel(dashboard)
	addCompletedRPCPanel(dashboard)
	addLatencyPanel(dashboard)
	addGCPanel(dashboard)
	addGoroutinePanel(dashboard)
	return dashboard.Time("now-3h", "now")
}

func ValdIndexManager() *dashboard.DashboardBuilder {
	dashboard := dashboard.NewDashboardBuilder("Vald Index Manager").Uid("vald_index_manager")
	addVariables(dashboard, "index manager")
	addStatPanels(dashboard)
	addMemoryPanel(dashboard)
	addCPUPanel(dashboard)
	addCompletedRPCPanel(dashboard)
	addLatencyPanel(dashboard)
	addGCPanel(dashboard)
	addGoroutinePanel(dashboard)
	return dashboard.Time("now-3h", "now")
}

func ValdIndexCorrection() *dashboard.DashboardBuilder {
	builder := dashboard.NewDashboardBuilder("Vald Index Correction").Uid("vald_index_correction")
	addVariables(builder, "index correction job")
	builder.
		WithRow(dashboard.NewRowBuilder("Stats"))
	addStatPanels(builder)
	builder.
		WithRow(dashboard.NewRowBuilder("Index Correction Stats"))
	addIndexCorrectionPanels(builder)
	builder.
		WithRow(dashboard.NewRowBuilder("Time Series"))
	addMemoryPanel(builder)
	addCPUPanel(builder)
	addGCPanel(builder)
	addGoroutinePanel(builder)
	return builder.Time("now-3h", "now")
}

func ValdAgent() *dashboard.DashboardBuilder {
	builder := dashboard.NewDashboardBuilder("Vald Agent").Uid("vald_agent")
	builder.
		WithRow(dashboard.NewRowBuilder("Stats"))
	addVariables(builder, "agent .*")
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
	addCompletedRPCPanel(builder)
	addLatencyPanel(builder)
	addGCPanel(builder)
	addGoroutinePanel(builder)
	return builder.Time("now-3h", "now")
}

func ValdAgentMemory() *dashboard.DashboardBuilder {
	builder := dashboard.NewDashboardBuilder("Vald Agent Memory").Uid("vald_agent_memory")
	builder.
		WithRow(dashboard.NewRowBuilder("runtime.Memstats"))
	addVariables(builder, "agent .*")
	addAgentMemoryPanels(builder)
	return builder.Time("now-15m", "now").Refresh("5s")
}

func main() {
	// Required to correctly unmarshal panels and dataqueries
	plugins.RegisterDefaultPlugins()

	cfg := &goapi.TransportConfig{
		// Host is the doman name or IP address of the host that serves the API.
		Host: "localhost:3000",
		// BasePath is the URL prefix for all API paths, relative to the host root.
		BasePath: "/api",
		// Schemes are the transfer protocols used by the API (http or https).
		Schemes: []string{"http"},
		// BasicAuth is optional basic auth credentials.
		BasicAuth: url.UserPassword("admin", "admin"),
	}

	client := goapi.NewHTTPClientWithConfig(strfmt.Default, cfg)

	for _, dashboard := range []func() *dashboard.DashboardBuilder{ValdLBGateway, ValdDiscoverer, ValdIndexManager, ValdAgent, ValdIndexCorrection, ValdAgentMemory} {
		dashboardModel, err := dashboard().Build()
		if err != nil {
			panic(err)
		}
		_, err = json.MarshalIndent(dashboardModel, "", "  ")
		if err != nil {
			panic(err)
		}

		// fmt.Println(string(dashboardJson))
		_, err = client.Dashboards.PostDashboard(&models.SaveDashboardCommand{
			Dashboard: dashboardModel,
			Overwrite: true,
		})
		if err != nil {
			panic(err)
		}
	}

	// fmt.Printf("response: %#v\n", response.Payload)
}

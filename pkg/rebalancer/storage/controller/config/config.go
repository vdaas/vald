//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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

// Package setting stores all server application settings
package config

import (
	"strings"

	"github.com/vdaas/vald/internal/config"
)

type GlobalConfig = config.GlobalConfig

type RebalanceReason uint8

const (
	DEVIATION RebalanceReason = iota
	RECOVERY
	MANUAL
)

func (r RebalanceReason) String() string {
	switch r {
	case DEVIATION:
		return "deviation"
	case RECOVERY:
		return "recovery"
	case MANUAL:
		return "manual"
	default:
		return "unknown"
	}
}

type AgentResourceType uint8

const (
	UNKNOWN_RESOURCE_TYPE AgentResourceType = iota
	STATEFULSET
	DAEMONSET
	REPLICASET
)

func (t AgentResourceType) String() string {
	switch t {
	case STATEFULSET:
		return "statefulset"
	case REPLICASET:
		return "replicaset"
	case DAEMONSET:
		return "daemonset"
	default:
		return "unknown"
	}
}

func AToAgentResourceType(t string) AgentResourceType {
	switch strings.ToLower(t) {
	case "statefulset":
		return STATEFULSET
	case "replicaset":
		return REPLICASET
	case "daemonset":
		return DAEMONSET
	default:
		return UNKNOWN_RESOURCE_TYPE
	}
}

// Config represent a application setting data content (config.yaml).
// In K8s environment, this configuration is stored in K8s ConfigMap.
type Data struct {
	config.GlobalConfig `json:",inline" yaml:",inline"`

	// Server represent all server configurations
	Server *config.Servers `json:"server_config" yaml:"server_config"`

	// Observability represent observability configurations
	Observability *config.Observability `json:"observability" yaml:"observability"`

	// Rebalancer represents rebalance controller configurations
	Rebalancer *config.RebalanceController `json:"rebalancer" yaml:"rebalancer"`
}

func NewConfig(path string) (cfg *Data, err error) {
	err = config.Read(path, &cfg)

	if err != nil {
		return nil, err
	}

	if cfg != nil {
		cfg.Bind()
	}

	if cfg.Server != nil {
		cfg.Server = cfg.Server.Bind()
	}

	if cfg.Observability != nil {
		cfg.Observability = cfg.Observability.Bind()
	}

	if cfg.Rebalancer != nil {
		cfg.Rebalancer = cfg.Rebalancer.Bind()
	} else {
		cfg.Rebalancer = new(config.RebalanceController)
	}
	return cfg, nil
}

// func FakeData() {
// 	d := Data{
// 		Version: "v0.0.1",
// 		Server: &config.Servers{
// 			Servers: []*config.Server{
// 				{
// 					Name:              "agent-rest",
// 					Host:              "127.0.0.1",
// 					Port:              8080,
// 					Mode:              "REST",
// 					ProbeWaitTime:     "3s",
// 					ShutdownDuration:  "5s",
// 					HandlerTimeout:    "5s",
// 					IdleTimeout:       "2s",
// 					ReadHeaderTimeout: "1s",
// 					ReadTimeout:       "1s",
// 					WriteTimeout:      "1s",
// 				},
// 				{
// 					Name: "agent-grpc",
// 					Host: "127.0.0.1",
// 					Port: 8082,
// 					Mode: "GRPC",
// 				},
// 			},
// 			MetricsServers: []*config.Server{
// 				{
// 					Name:              "pprof",
// 					Host:              "127.0.0.1",
// 					Port:              6060,
// 					Mode:              "REST",
// 					ProbeWaitTime:     "3s",
// 					ShutdownDuration:  "5s",
// 					HandlerTimeout:    "5s",
// 					IdleTimeout:       "2s",
// 					ReadHeaderTimeout: "1s",
// 					ReadTimeout:       "1s",
// 					WriteTimeout:      "1s",
// 				},
// 			},
// 			HealthCheckServers: []*config.Server{
// 				{
// 					Name: "livenesss",
// 					Host: "127.0.0.1",
// 					Port: 3000,
// 				},
// 				{
// 					Name: "readiness",
// 					Host: "127.0.0.1",
// 					Port: 3001,
// 				},
// 			},
// 			StartUpStrategy: []string{
// 				"livenesss",
// 				"pprof",
// 				"agent-grpc",
// 				"agent-rest",
// 				"readiness",
// 			},
// 			ShutdownStrategy: []string{
// 				"readiness",
// 				"agent-rest",
// 				"agent-grpc",
// 				"pprof",
// 				"livenesss",
// 			},
// 			FullShutdownDuration: "30s",
// 			TLS: &config.TLS{
// 				Enabled: false,
// 				Cert:    "/path/to/cert",
// 				Key:     "/path/to/key",
// 				CA:      "/path/to/ca",
// 			},
// 		},
// 		Rebalancer: &config.RebalanceController{
// 			RebalanceJobName:        "agent-rebalance-job",
// 			RebalanceJobNamespace:   "vald",
// 			RebalanceJobTemplateKey: "job.tpl",
// 			ConfigMapName:           "agent-rebalance-job-template",
// 			ConfigMapNamespace:      "vald",
// 			AgentName:               "vald-agent-ngt",
// 			AgentNamespace:          "vald",
// 			AgentResourceType:       "statefulset",
// 			ReconcileCheckDuration:  "5m",
// 			Tolerance:               0.1,
// 			RateThreshold:           0.1,
// 			LeaderElectionID:        "agent-rebalance-controller",
// 		},
// 	}
// 	fmt.Println(config.ToRawYaml(d))
// }

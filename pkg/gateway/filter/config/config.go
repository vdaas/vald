//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errors"
)

type (
	GlobalConfig = config.GlobalConfig
	Server       = config.Server
)

// Config represent a application setting data content (config.yaml).
// In K8s environment, this configuration is stored in K8s ConfigMap.
type Data struct {
	config.GlobalConfig `json:",inline" yaml:",inline"`

	// Server represent all server configurations
	Server *config.Servers `json:"server_config" yaml:"server_config"`

	// Observability represent observability configurations
	Observability *config.Observability `json:"observability" yaml:"observability"`

	// Client represent gateway client configuration
	Client *config.GRPCClient `json:"client" yaml:"client"`

	// IngressFilters represent filter configuration for ingress payload
	IngressFilters *config.IngressFilter `json:"ingress_filter" yaml:"ingress_filter"`

	// EgressFilters represent filter configuration for egress payload
	EgressFilters *config.EgressFilter `json:"egress_filter" yaml:"egress_filter"`
}

func NewConfig(path string) (cfg *Data, err error) {
	cfg = new(Data)

	err = config.Read(path, &cfg)

	if err != nil {
		return nil, err
	}

	if cfg != nil {
		cfg.Bind()
	} else {
		return nil, errors.ErrInvalidConfig
	}

	if cfg.Server != nil {
		cfg.Server = cfg.Server.Bind()
	} else {
		return nil, errors.ErrInvalidConfig
	}

	if cfg.Observability != nil {
		cfg.Observability = cfg.Observability.Bind()
	} else {
		cfg.Observability = new(config.Observability).Bind()
	}

	if cfg.Client != nil {
		cfg.Client = cfg.Client.Bind()
	} else {
		return nil, errors.ErrInvalidConfig
	}

	if cfg.IngressFilters != nil {
		cfg.IngressFilters = cfg.IngressFilters.Bind()
	} else {
		cfg.IngressFilters = new(config.IngressFilter).Bind()
	}

	if cfg.EgressFilters != nil {
		cfg.EgressFilters = cfg.EgressFilters.Bind()
	} else {
		cfg.EgressFilters = new(config.EgressFilter).Bind()
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
// 		Gateway: &config.Gateway{
// 			AgentPort:           8080,
// 			AgentName:           "vald-agent",
// 			BackoffEnabled:      false,,
// 		},
// 	}
// 	fmt.Println(config.ToRawYaml(d))
// }

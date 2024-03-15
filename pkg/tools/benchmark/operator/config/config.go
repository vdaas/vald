//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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

// Package config stores all server application settings
package config

import (
	"github.com/vdaas/vald/internal/config"
)

// GlobalConfig is type alias for config.GlobalConfig.
type GlobalConfig = config.GlobalConfig

// Config represent a application setting data content (config.yaml).
// In K8s environment, this configuration is stored in K8s ConfigMap.
type Config struct {
	config.GlobalConfig `json:",inline" yaml:",inline"`

	// Server represent all server configurations
	Server *config.Servers `json:"server_config" yaml:"server_config"`

	// Observability represent observability configurations
	Observability *config.Observability `json:"observability" yaml:"observability"`

	// JobImage represents the location of Docker image for benchmark job and its ImagePullPolicy
	JobImage *config.BenchmarkJobImageInfo `json:"job_image" yaml:"job_image"`
}

// NewConfig represents the set config from the given setting file path.
func NewConfig(path string) (cfg *Config, err error) {
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

	if cfg.JobImage != nil {
		cfg.JobImage = cfg.JobImage.Bind()
	} else {
		cfg.JobImage = new(config.BenchmarkJobImageInfo)
	}
	return cfg, nil
}

// func FakeData() {
// 	d := Config{
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
// 		Scenario: &config.BenchmarkScenario{
// 			Target: &config.BenchmarkTarget{
// 				Host: "localhost",
// 				Port: 8081,
// 			},
// 			Dataset: &config.BenchmarkDataset{
// 				Name:    "fashion-mnist-784-euc",
// 				Group:   "Train",
// 				Indexes: 10000,
// 				Range: &config.BenchmarkDatasetRange{
// 					Start: 100000,
// 					End:   200000,
// 				},
// 			},
// 			Jobs: []*config.BenchmarkJob{
// 				{
// 					JobType:    "search",
// 					Replica:    1,
// 					Repetition: 1,
// 					Dimension:  784,
// 					Iter:       10000,
// 					Num:        10,
// 					MinNum:     100,
// 					Radius:     -1,
// 					Epsilon:    0.1,
// 					Timeout:    "5s",
// 				},
// 			},
// 		},
// 	}
// 	fmt.Println(config.ToRawYaml(d))
// }

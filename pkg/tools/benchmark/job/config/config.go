//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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

// Package config stores all server application settings
package config

import (
	"context"
	"os"

	"github.com/vdaas/vald/internal/config"
	v1 "github.com/vdaas/vald/internal/k8s/vald/benchmark/api/v1"
	"github.com/vdaas/vald/internal/log"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// GlobalConfig is type alias for config.GlobalConfig
type GlobalConfig = config.GlobalConfig

// Config represent a application setting data content (config.yaml).
// In K8s environment, this configuration is stored in K8s ConfigMap.
type Config struct {
	config.GlobalConfig `json:",inline" yaml:",inline"`

	// Server represent all server configurations
	Server *config.Servers `json:"server_config" yaml:"server_config"`

	// Observability represent observability configurations
	Observability *config.Observability `json:"observability" yaml:"observability"`

	// Job represents benchmark job configurations
	Job *config.BenchmarkJob `json:"job" yaml:"job"`
}

var (
	NAMESPACE = os.Getenv("POD_NAMESPACE")
	NAME      = os.Getenv("POD_NAME")
	mgr       manager.Manager
)

// NewConfig represents the set config from the given setting file path.
func NewConfig(ctx context.Context, path string) (cfg *Config, err error) {
	err = config.Read(path, &cfg)
	if err != nil {
		return nil, err
	}

	if cfg != nil {
		cfg.Bind()
	}

	if cfg.Server != nil {
		cfg.Server = cfg.Server.Bind()
	} else {
		cfg.Server = new(config.Servers)
	}

	if cfg.Observability != nil {
		cfg.Observability = cfg.Observability.Bind()
	} else {
		cfg.Observability = new(config.Observability)
	}

	if cfg.Job != nil {
		cfg.Job = cfg.Job.Bind()
	} else {
		cfg.Job = new(config.BenchmarkJob)
	}

	if cfg.Job.GatewayClient == nil {
		cfg.Job.GatewayClient = new(config.GRPCClient)
		cfg.Job.GatewayClient.Addrs = []string{
			"vald-lb-gateway.default.svc.cluster.local:8081",
		}
	}

	// Get config from applied ValdBenchmarkJob custom resource
	scheme := runtime.NewScheme()
	v1.AddToScheme(scheme)
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), manager.Options{
		Scheme: scheme,
	})
	if err != nil {
		log.Warn(err.Error())
	}
	cli := mgr.GetAPIReader()
	var jobResource v1.ValdBenchmarkJob
	err = cli.Get(ctx, client.ObjectKey{
		Name:      NAME,
		Namespace: NAMESPACE,
	}, &jobResource)
	if err != nil {
		log.Warn(err.Error())
	}

	cfg.Job.Target = jobResource.Spec.Target
	cfg.Job.Dataset = jobResource.Spec.Dataset
	cfg.Job.Replica = jobResource.Spec.Replica
	cfg.Job.Repetition = jobResource.Spec.Repetition
	cfg.Job.JobType = jobResource.Spec.JobType
	cfg.Job.Dimension = jobResource.Spec.Dimension
	cfg.Job.Epsilon = float64(jobResource.Spec.Epsilon)
	cfg.Job.Radius = float64(jobResource.Spec.Radius)
	cfg.Job.Num = uint32(jobResource.Spec.Num)
	cfg.Job.MinNum = uint32(jobResource.Spec.MinNum)
	cfg.Job.Timeout = jobResource.Spec.Timeout
	cfg.Job.Rules = jobResource.Spec.Rules
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
// 		Job: &config.BenchmarkJob{
//			JobType:   "search",
//			Dataset:   "fashion-mnist-784-euc",
//			Dimension: 784,
//			Iter:      100,
//			Num:       10,
//			MinNum:    5,
//			Radius:    -1,
//			Epsilon:   0.1,
//			Timeout:   "5s",
// 		},
// 	}
// 	fmt.Println(config.ToRawYaml(d))
// }

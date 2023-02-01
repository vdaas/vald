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
	"time"

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

	// Server represent all server configuration
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

	if cfg.Job.ClientConfig == nil {
		cfg.Job.ClientConfig = new(config.GRPCClient)
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
	} else {
		cfg.Job.Target = jobResource.Spec.Target
		cfg.Job.Dataset = jobResource.Spec.Dataset
		cfg.Job.Replica = jobResource.Spec.Replica
		cfg.Job.Repetition = jobResource.Spec.Repetition
		cfg.Job.JobType = jobResource.Spec.JobType
		cfg.Job.Rules = jobResource.Spec.Rules
		cfg.Job.InsertConfig = jobResource.Spec.InsertConfig
		cfg.Job.UpdateConfig = jobResource.Spec.UpdateConfig
		cfg.Job.UpsertConfig = jobResource.Spec.UpsertConfig
		cfg.Job.SearchConfig = jobResource.Spec.SearchConfig
		cfg.Job.RemoveConfig = jobResource.Spec.RemoveConfig
		cfg.Job.ClientConfig = convertClientOpts(jobResource.Spec.ClientConfig)
	}

	return cfg, nil
}

func convertClientOpts(g *v1.ClientConfig) *config.GRPCClient {
	cfg := new(config.GRPCClient)
	if g == nil {
		return cfg
	}

	if len(g.Addrs) > 0 {
		cfg.Addrs = g.Addrs
	}

	_, err := time.ParseDuration(g.HealthCheckDuration)
	if err == nil {
		cfg.HealthCheckDuration = g.HealthCheckDuration
	} else {
		cfg.HealthCheckDuration = "10s"
	}

	if g.ConnectionPool != nil {
		cfg.ConnectionPool = (*config.ConnectionPool)(g.ConnectionPool)
	}

	if g.Backoff != nil {
		cfg.Backoff = (*config.Backoff)(g.Backoff)
	}

	if g.CircuitBreaker != nil {
		cfg.CircuitBreaker = (*config.CircuitBreaker)(g.CircuitBreaker)
	}

	if g.CallOption != nil {
		cfg.CallOption = (*config.CallOption)(g.CallOption)
	}

	if g.DialOption != nil {
		cfg.DialOption.WriteBufferSize = g.DialOption.WriteBufferSize
		cfg.DialOption.ReadBufferSize = g.DialOption.ReadBufferSize
		cfg.DialOption.InitialWindowSize = g.DialOption.InitialWindowSize
		cfg.DialOption.MaxMsgSize = g.DialOption.MaxMsgSize
		cfg.DialOption.BackoffMaxDelay = g.DialOption.BackoffMaxDelay
		cfg.DialOption.BackoffBaseDelay = g.DialOption.BackoffBaseDelay
		cfg.DialOption.BackoffMultiplier = g.DialOption.BackoffMultiplier
		cfg.DialOption.MinimumConnectionTimeout = g.DialOption.MinimumConnectionTimeout
		cfg.DialOption.EnableBackoff = g.DialOption.EnableBackoff
		cfg.DialOption.Insecure = g.DialOption.Insecure
		cfg.DialOption.Timeout = g.DialOption.Timeout
		cfg.DialOption.Interceptors = g.DialOption.Interceptors
		if g.DialOption.Net != nil {
			if g.DialOption.Net.DNS != nil {
				cfg.DialOption.Net.DNS = (*config.DNS)(g.DialOption.Net.DNS)
			}
			if g.DialOption.Net.Dialer != nil {
				cfg.DialOption.Net.Dialer = (*config.Dialer)(g.DialOption.Net.Dialer)
			}
			if g.DialOption.Net.SocketOption != nil {
				cfg.DialOption.Net.SocketOption = (*config.SocketOption)(g.DialOption.Net.SocketOption)
			}
			if g.DialOption.Net.TLS != nil && g.DialOption.Net.TLS.Enabled {
				cfg.DialOption.Net.TLS = (*config.TLS)(g.DialOption.Net.TLS)
			}
		}
	}

	if g.TLS != nil && g.TLS.Enabled {
		cfg.TLS = (*config.TLS)(g.TLS)
	}

	return cfg
}

// func FakeData() {
// 	d := Config{
// 		Version: "v0.0.1",
// 		Server: &config.Servers{
// 			Servers: []*config.Server{
// 				{
// 					Name:          "agent-rest",
// 					Host:          "127.0.0.1",
// 					Port:          8080,
// 					Mode:          "REST",
// 					ProbeWaitTime: "3s",
// 					HTTP: &config.HTTP{
// 						ShutdownDuration:  "5s",
// 						HandlerTimeout:    "5s",
// 						IdleTimeout:       "2s",
// 						ReadHeaderTimeout: "1s",
// 						ReadTimeout:       "1s",
// 						WriteTimeout:      "1s",
// 					},
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
// 					Name:          "pprof",
// 					Host:          "127.0.0.1",
// 					Port:          6060,
// 					Mode:          "REST",
// 					ProbeWaitTime: "3s",
// 					HTTP: &config.HTTP{
// 						ShutdownDuration:  "5s",
// 						HandlerTimeout:    "5s",
// 						IdleTimeout:       "2s",
// 						ReadHeaderTimeout: "1s",
// 						ReadTimeout:       "1s",
// 						WriteTimeout:      "1s",
// 					},
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
// 			Target: &v1.BenchmarkTarget{
// 				Host: "vald-lb-gateway.svc.local",
// 				Port: 8081,
// 			},
// 			Dataset: &v1.BenchmarkDataset{
// 				Name:    "fashion-mnist",
// 				Group:   "train",
// 				Indexes: 10000,
// 				Range: &v1.BenchmarkDatasetRange{
// 					Start: 0,
// 					End:   10000,
// 				},
// 			},
// 			Replica:      1,
// 			Repetition:   1,
// 			JobType:      "search",
// 			InsertConfig: &v1.InsertJobConfig{},
// 			UpdateConfig: &v1.UpdateJobConfig{},
// 			UpsertConfig: &v1.UpsertJobConfig{},
// 			SearchConfig: &v1.SearchJobConfig{},
// 			RemoveConfig: &v1.RemoveJobConfig{},
// 			ClientConfig: &config.GRPCClient{
// 				Addrs:               []string{},
// 				HealthCheckDuration: "1s",
// 				ConnectionPool: &config.ConnectionPool{
// 					ResolveDNS:           true,
// 					EnableRebalance:      true,
// 					RebalanceDuration:    "30m",
// 					Size:                 3,
// 					OldConnCloseDuration: "2m",
// 				},
// 				Backoff: &config.Backoff{
// 					InitialDuration:  "5ms",
// 					BackoffTimeLimit: "5s",
// 					MaximumDuration:  "5s",
// 					JitterLimit:      "100ms",
// 					BackoffFactor:    1.1,
// 					RetryCount:       100,
// 					EnableErrorLog:   true,
// 				},
// 				CircuitBreaker: &config.CircuitBreaker{
// 					ClosedErrorRate:      0.7,
// 					HalfOpenErrorRate:    0.5,
// 					MinSamples:           1000,
// 					OpenTimeout:          "1s",
// 					ClosedRefreshTimeout: "10s",
// 				},
// 				CallOption: &config.CallOption{
// 					WaitForReady:          true,
// 					MaxRetryRPCBufferSize: 0,
// 					MaxRecvMsgSize:        0,
// 					MaxSendMsgSize:        0,
// 				},
// 				DialOption: &config.DialOption{
// 					WriteBufferSize:             0,
// 					ReadBufferSize:              0,
// 					InitialWindowSize:           0,
// 					InitialConnectionWindowSize: 0,
// 					MaxMsgSize:                  0,
// 					BackoffMaxDelay:             "120s",
// 					BackoffBaseDelay:            "1s",
// 					BackoffJitter:               0.2,
// 					BackoffMultiplier:           1.6,
// 					MinimumConnectionTimeout:    "20s",
// 					EnableBackoff:               true,
// 					Insecure:                    true,
// 					Timeout:                     "",
// 					Interceptors:                []string{},
// 					Net: &config.Net{
// 						DNS: &config.DNS{
// 							CacheEnabled:    true,
// 							RefreshDuration: "30m",
// 							CacheExpiration: "1h",
// 						},
// 						Dialer: &config.Dialer{
// 							Timeout:          "",
// 							Keepalive:        "",
// 							FallbackDelay:    "",
// 							DualStackEnabled: true,
// 						},
// 						TLS: &config.TLS{
// 							Enabled:            false,
// 							Cert:               "path/to/cert",
// 							Key:                "path/to/key",
// 							CA:                 "path/to/ca",
// 							InsecureSkipVerify: false,
// 						},
// 						SocketOption: &config.SocketOption{
// 							ReusePort:                true,
// 							ReuseAddr:                true,
// 							TCPFastOpen:              true,
// 							TCPNoDelay:               true,
// 							TCPQuickAck:              true,
// 							TCPCork:                  false,
// 							TCPDeferAccept:           true,
// 							IPTransparent:            false,
// 							IPRecoverDestinationAddr: false,
// 						},
// 					},
// 					Keepalive: &config.GRPCClientKeepalive{
// 						Time:                "120s",
// 						Timeout:             "30s",
// 						PermitWithoutStream: true,
// 					},
// 				},
// 				TLS: &config.TLS{
// 					Enabled:            false,
// 					Cert:               "path/to/cert",
// 					Key:                "path/to/key",
// 					CA:                 "path/to/ca",
// 					InsecureSkipVerify: false,
// 				},
// 			},
// 			Rules: []*v1.BenchmarkJobRule{},
// 		},
// 	}
// 	fmt.Println(config.ToRawYaml(d))
// }

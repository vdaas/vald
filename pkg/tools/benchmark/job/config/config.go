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
	"context"
	"encoding/json"
	"os"

	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/k8s/client"
	v1 "github.com/vdaas/vald/internal/k8s/vald/benchmark/api/v1"
	"github.com/vdaas/vald/internal/log"
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

	// K8sClient represents kubernetes clients
	K8sClient client.Client `json:"k8s_client" yaml:"k8s_client"`
}

var (
	NAMESPACE               = os.Getenv("CRD_NAMESPACE")
	NAME                    = os.Getenv("CRD_NAME")
	JOBNAME_ANNOTATION      = "before-job-name"
	JOBNAMESPACE_ANNOTATION = "before-job-namespace"
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
	var jobResource v1.ValdBenchmarkJob
	if cfg.K8sClient == nil {
		c, err := client.New(client.WithSchemeBuilder(*v1.SchemeBuilder))
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}
		cfg.K8sClient = c
	}
	err = cfg.K8sClient.Get(ctx, NAME, NAMESPACE, &jobResource)
	if err != nil {
		log.Warn(err.Error())
	} else {
		// create override Config
		overrideCfg := new(Config)
		if jobResource.Spec.GlobalConfig != nil {
			overrideCfg.GlobalConfig = *jobResource.Spec.Bind()
		}
		if jobResource.Spec.ServerConfig != nil {
			overrideCfg.Server = (*jobResource.Spec.ServerConfig).Bind()
		}
		// jobResource.Spec has another field comparering Config.Job, so json.Marshal and Unmarshal are used for embedding field value of Config.Job from jobResource.Spec
		var overrideJobCfg config.BenchmarkJob
		b, err := json.Marshal(*jobResource.Spec.DeepCopy())
		if err == nil {
			err = json.Unmarshal([]byte(b), &overrideJobCfg)
			if err != nil {
				log.Warn(err.Error())
			}
			overrideCfg.Job = overrideJobCfg.Bind()
		}
		if annotations := jobResource.GetAnnotations(); annotations != nil {
			overrideCfg.Job.BeforeJobName = annotations[JOBNAME_ANNOTATION]
			overrideCfg.Job.BeforeJobNamespace = annotations[JOBNAMESPACE_ANNOTATION]
		}
		return config.Merge(cfg, overrideCfg)
	}
	return cfg, nil
}

// func FakeData() {
// 	d := Config{
// 		GlobalConfig: config.GlobalConfig{
// 			Version: "v0.0.1",
// 			TZ:      "JST",
// 			Logging: &config.Logging{
// 				Format: "raw",
// 				Level:  "debug",
// 				Logger: "glg",
// 			},
// 		},
// 		Server: &config.Servers{
// 			Servers: []*config.Server{
// 				{
// 					Name:          "grpc",
// 					Host:          "0.0.0.0",
// 					Port:          8081,
// 					Mode:          "GRPC",
// 					ProbeWaitTime: "3s",
// 					SocketPath:    "",
// 					GRPC: &config.GRPC{
// 						BidirectionalStreamConcurrency: 20,
// 						MaxReceiveMessageSize:          0,
// 						MaxSendMessageSize:             0,
// 						InitialWindowSize:              1048576,
// 						InitialConnWindowSize:          2097152,
// 						Keepalive: &config.GRPCKeepalive{
// 							MaxConnIdle:         "",
// 							MaxConnAge:          "",
// 							MaxConnAgeGrace:     "",
// 							Time:                "3h",
// 							Timeout:             "60s",
// 							MinTime:             "10m",
// 							PermitWithoutStream: true,
// 						},
// 						WriteBufferSize:   0,
// 						ReadBufferSize:    0,
// 						ConnectionTimeout: "",
// 						MaxHeaderListSize: 0,
// 						HeaderTableSize:   0,
// 						Interceptors: []string{
// 							"RecoverInterceptor",
// 						},
// 						EnableReflection: true,
// 					},
// 					SocketOption: &config.SocketOption{
// 						ReusePort:                true,
// 						ReuseAddr:                true,
// 						TCPFastOpen:              false,
// 						TCPNoDelay:               false,
// 						TCPCork:                  false,
// 						TCPQuickAck:              false,
// 						TCPDeferAccept:           false,
// 						IPTransparent:            false,
// 						IPRecoverDestinationAddr: false,
// 					},
// 					Restart: true,
// 				},
// 			},
// 			HealthCheckServers: []*config.Server{
// 				{
// 					Name:          "livenesss",
// 					Host:          "0.0.0.0",
// 					Port:          3000,
// 					Mode:          "",
// 					Network:       "tcp",
// 					ProbeWaitTime: "3s",
// 					SocketPath:    "",
// 					HTTP: &config.HTTP{
// 						HandlerTimeout:    "",
// 						IdleTimeout:       "",
// 						ReadHeaderTimeout: "",
// 						ReadTimeout:       "",
// 						ShutdownDuration:  "5s",
// 						WriteTimeout:      "",
// 					},
// 					SocketOption: &config.SocketOption{
// 						ReusePort:                true,
// 						ReuseAddr:                true,
// 						TCPFastOpen:              true,
// 						TCPNoDelay:               true,
// 						TCPCork:                  false,
// 						TCPQuickAck:              true,
// 						TCPDeferAccept:           false,
// 						IPTransparent:            false,
// 						IPRecoverDestinationAddr: false,
// 					},
// 				},
// 				{
// 					Name:          "readiness",
// 					Host:          "0.0.0.0",
// 					Port:          3001,
// 					Mode:          "",
// 					Network:       "tcp",
// 					ProbeWaitTime: "3s",
// 					SocketPath:    "",
// 					HTTP: &config.HTTP{
// 						HandlerTimeout:    "",
// 						IdleTimeout:       "",
// 						ReadHeaderTimeout: "",
// 						ReadTimeout:       "",
// 						ShutdownDuration:  "0s",
// 						WriteTimeout:      "",
// 					},
// 					SocketOption: &config.SocketOption{
// 						ReusePort:                true,
// 						ReuseAddr:                true,
// 						TCPFastOpen:              true,
// 						TCPNoDelay:               true,
// 						TCPCork:                  false,
// 						TCPQuickAck:              true,
// 						TCPDeferAccept:           false,
// 						IPTransparent:            false,
// 						IPRecoverDestinationAddr: false,
// 					},
// 				},
// 			},
// 			MetricsServers: []*config.Server{
// 				{
// 					Name:          "pprof",
// 					Host:          "0.0.0.0",
// 					Port:          6060,
// 					Mode:          "REST",
// 					Network:       "tcp",
// 					ProbeWaitTime: "3s",
// 					SocketPath:    "",
// 					HTTP: &config.HTTP{
// 						HandlerTimeout:    "5s",
// 						IdleTimeout:       "2s",
// 						ReadHeaderTimeout: "1s",
// 						ReadTimeout:       "1s",
// 						ShutdownDuration:  "5s",
// 						WriteTimeout:      "1m",
// 					},
// 					SocketOption: &config.SocketOption{
// 						ReusePort:                true,
// 						ReuseAddr:                true,
// 						TCPFastOpen:              false,
// 						TCPNoDelay:               false,
// 						TCPCork:                  false,
// 						TCPQuickAck:              false,
// 						TCPDeferAccept:           false,
// 						IPTransparent:            false,
// 						IPRecoverDestinationAddr: false,
// 					},
// 				},
// 			},
// 			StartUpStrategy: []string{
// 				"livenesss",
// 				"readiness",
// 				"pprof",
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
// 		Observability: &config.Observability{
// 			Enabled: true,
// 			OTLP: &config.OTLP{
// 				CollectorEndpoint: "",
// 				Attribute: &config.OTLPAttribute{
// 					Namespace:   NAMESPACE,
// 					PodName:     NAME,
// 					NodeName:    "",
// 					ServiceName: "vald",
// 				},
// 				TraceBatchTimeout:       "1s",
// 				TraceExportTimeout:      "1m",
// 				TraceMaxExportBatchSize: 1024,
// 				TraceMaxQueueSize:       256,
// 				MetricsExportInterval:   "1s",
// 				MetricsExportTimeout:    "1m",
// 			},
// 			Metrics: &config.Metrics{
// 				EnableVersionInfo: true,
// 				EnableMemory:      true,
// 				EnableGoroutine:   true,
// 				EnableCGO:         true,
// 				VersionInfoLabels: []string{
// 					"vald_version",
// 					"server_name",
// 					"git_commit",
// 					"build_time",
// 					"go_version",
// 					"go_arch",
// 					"algorithm_info",
// 				},
// 			},
// 			Trace: &config.Trace{
// 				Enabled: true,
// 			},
// 		},
// 		Job: &config.BenchmarkJob{
// 			Target: &config.BenchmarkTarget{
// 				Host: "vald-lb-gateway.svc.local",
// 				Port: 8081,
// 			},
// 			Dataset: &config.BenchmarkDataset{
// 				Name:    "fashion-mnist",
// 				Group:   "train",
// 				Indexes: 10000,
// 				Range: &config.BenchmarkDatasetRange{
// 					Start: 0,
// 					End:   10000,
// 				},
// 			},
// 			Replica:      1,
// 			Repetition:   1,
// 			JobType:      "search",
// 			InsertConfig: &config.InsertConfig{},
// 			UpdateConfig: &config.UpdateConfig{},
// 			UpsertConfig: &config.UpsertConfig{},
// 			SearchConfig: &config.SearchConfig{},
// 			RemoveConfig: &config.RemoveConfig{},
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
// 					InitialWindowSize:           1048576,
// 					InitialConnectionWindowSize: 2097152,
// 					MaxMsgSize:                  0,
// 					BackoffMaxDelay:             "120s",
// 					BackoffBaseDelay:            "1s",
// 					BackoffJitter:               0.2,
// 					BackoffMultiplier:           1.6,
// 					MinimumConnectionTimeout:    "20s",
// 					EnableBackoff:               false,
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
// 							TCPFastOpen:              false,
// 							TCPNoDelay:               false,
// 							TCPQuickAck:              false,
// 							TCPCork:                  false,
// 							TCPDeferAccept:           false,
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
// 			Rules: []*config.BenchmarkJobRule{},
// 		},
// 	}
// 	fmt.Println(config.ToRawYaml(d))
// }

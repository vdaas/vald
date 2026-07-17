// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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
package config

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/os"
)

func TestNewConfig(t *testing.T) {
	t.Parallel()

	const (
		testNamespace   = "default"
		testRotatorName = "vald-readreplica-rotate"
	)

	const validYAML = `
version: v1.0.0
server_config:
  full_shutdown_duration: 10ms
observability:
  enabled: true
operator:
  namespace: default
  agent_name: vald-agent
  agent_namespace: default
  rotator_name: vald-readreplica-rotate
  target_read_replica_id_annotations_key: vald.vdaas.org/target-read-replica-id
  rotation_job_concurrency: 2
  read_replica_enabled: true
  read_replica_label_key: vald-readreplica-id
  job_templates:
    rotate:
      metadata:
        name: vald-readreplica-rotate
        namespace: default
`

	const missingObservabilityYAML = `
version: v1.0.0
server_config:
  full_shutdown_duration: 10ms
operator:
  namespace: default
  agent_name: vald-agent
  rotator_name: vald-readreplica-rotate
  target_read_replica_id_annotations_key: vald.vdaas.org/target-read-replica-id
  rotation_job_concurrency: 1
`

	const missingServerYAML = `
version: v1.0.0
operator:
  namespace: default
  rotation_job_concurrency: 1
`

	const missingOperatorYAML = `
version: v1.0.0
server_config:
  full_shutdown_duration: 10ms
`

	wantBoundServer := &config.Servers{
		FullShutdownDuration: "10ms",
		StartUpStrategy:      make([]string, 0),
		ShutdownStrategy:     make([]string, 0),
		TLS:                  &config.TLS{Enabled: false},
	}

	tests := []struct {
		wantCfg   *Data
		wantErr   error
		checkErr  func(*testing.T, error)
		name      string
		fileName  string
		data      string
		skipWrite bool
	}{
		{
			name:     "returns the bound Data when the yaml is valid",
			fileName: "valid.yaml",
			data:     validYAML,
			wantCfg: &Data{
				GlobalConfig: config.GlobalConfig{Version: "v1.0.0"},
				Server:       wantBoundServer,
				Observability: &config.Observability{
					Enabled: true,
					OTLP:    &config.OTLP{Attribute: new(config.OTLPAttribute)},
					Metrics: new(config.Metrics),
					Trace:   new(config.Trace),
				},
				Operator: &config.IndexOperator{
					Namespace:                         testNamespace,
					AgentName:                         "vald-agent",
					AgentNamespace:                    testNamespace,
					RotatorName:                       testRotatorName,
					TargetReadReplicaIDAnnotationsKey: "vald.vdaas.org/target-read-replica-id",
					RotationJobConcurrency:            2,
					ReadReplicaEnabled:                true,
					ReadReplicaLabelKey:               "vald-readreplica-id",
					JobTemplates: config.IndexJobTemplates{
						Rotate: &k8s.Job{
							ObjectMeta: k8s.ObjectMeta{
								Name:      testRotatorName,
								Namespace: testNamespace,
							},
						},
					},
				},
			},
		},
		{
			name:     "defaults Observability when it is omitted from the file",
			fileName: "no_observability.yaml",
			data:     missingObservabilityYAML,
			wantCfg: &Data{
				GlobalConfig: config.GlobalConfig{Version: "v1.0.0"},
				Server:       wantBoundServer,
				Observability: &config.Observability{
					OTLP:    &config.OTLP{Attribute: new(config.OTLPAttribute)},
					Metrics: new(config.Metrics),
					Trace:   new(config.Trace),
				},
				Operator: &config.IndexOperator{
					Namespace:                         testNamespace,
					AgentName:                         "vald-agent",
					RotatorName:                       testRotatorName,
					TargetReadReplicaIDAnnotationsKey: "vald.vdaas.org/target-read-replica-id",
					RotationJobConcurrency:            1,
				},
			},
		},
		{
			name:     "returns ErrInvalidConfig when server_config is missing",
			fileName: "no_server.yaml",
			data:     missingServerYAML,
			wantErr:  errors.ErrInvalidConfig,
		},
		{
			name:     "returns ErrInvalidConfig when operator is missing",
			fileName: "no_operator.yaml",
			data:     missingOperatorYAML,
			wantErr:  errors.ErrInvalidConfig,
		},
		{
			name:     "returns ErrInvalidConfig when the yaml file is empty",
			fileName: "empty.yaml",
			data:     "",
			wantErr:  errors.ErrInvalidConfig,
		},
		{
			name:      "returns an error when the file does not exist",
			fileName:  "does_not_exist.yaml",
			skipWrite: true,
			checkErr: func(t *testing.T, err error) {
				t.Helper()
				require.True(t, errors.Is(err, os.ErrNotExist), "expected a not-exist error, got: %v", err)
			},
		},
		{
			name:     "returns ErrUnsupportedConfigFileType for an unsupported extension",
			fileName: "config.txt",
			data:     "irrelevant contents",
			wantErr:  errors.ErrUnsupportedConfigFileType(".txt"),
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()

			dir := tt.TempDir()
			path := filepath.Join(dir, test.fileName)
			if !test.skipWrite {
				require.NoError(tt, os.WriteFile(path, []byte(test.data), 0o600))
			}

			gotCfg, err := NewConfig(path)

			switch {
			case test.checkErr != nil:
				test.checkErr(tt, err)
			case test.wantErr != nil:
				require.True(tt, errors.Is(err, test.wantErr), "got_error: %v, want: %v", err, test.wantErr)
				require.Nil(tt, gotCfg)
			default:
				require.NoError(tt, err)
				require.Equal(tt, test.wantCfg, gotCfg)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW

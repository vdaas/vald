//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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

package config

import (
	"bytes"
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/file"
)

func writeTempFile(t *testing.T, name, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if _, err := file.WriteFile(context.Background(), path, bytes.NewReader([]byte(content)), 0o600); err != nil {
		t.Fatalf("failed to write %s: %v", name, err)
	}
	return path
}

const minimalServerYAML = `server_config:
  servers:
    - name: grpc
      port: 8081
`

//nolint:maintidx
func TestNew(t *testing.T) {
	t.Parallel()

	type test struct {
		check   func(t *testing.T, cfg *Data)
		name    string
		yaml    string
		wantErr bool
	}

	tests := []test{
		{
			name: "returns config when all sections are given",
			yaml: `---
version: v0.0.0
` + minimalServerYAML + `
operator:
  name: vald-operator
  namespace: default
  controller:
    leader_election:
      enabled: true
      id: vald-operator
      namespace: default
      lease_duration: 30s
      renew_deadline: 20s
      retry_period: 5s
    metrics_address: ":9090"
    max_concurrent_reconciles: 2
    sync_period: 10h
    cache_namespaces:
      - default
    requeue:
      success: ""
      on_error: 100ms
      not_found: 1s
  vrs:
    default_vrs_path: /opt/valdoperatorrelease/config/vrs.yaml
    internal_host_domain: ""
    log_level: warn
  node_pool:
    require_match: true
    label_prefix: vald.vdaas.org
    agent_pods_per_node: 2
  persistent_volume:
    default_storage_class: standard
    default_access_mode: ReadWriteOnce
    buffer_ratio: 1.5
    min_size_bytes: 1073741824
  networking:
    enable_ingress: true
    gateway_ingress_annotations: {}
    gateway_service_type: NodePort
    discoverer_daemonset_max_surge: 30%
    discoverer_daemonset_max_unavailable: 0%
`,
			check: func(t *testing.T, cfg *Data) {
				t.Helper()
				o := cfg.Operator
				if o.Name != "vald-operator" {
					t.Errorf("Operator.Name = %q, want %q", o.Name, "vald-operator")
				}
				le := o.Controller.LeaderElection
				if le == nil || !le.Enabled || le.ID != "vald-operator" {
					t.Errorf("LeaderElection = %+v, want enabled with id vald-operator", le)
				}
				if le.LeaseDuration != "30s" || le.RenewDeadline != "20s" || le.RetryPeriod != "5s" {
					t.Errorf("LeaderElection timings = %q/%q/%q, want 30s/20s/5s",
						le.LeaseDuration, le.RenewDeadline, le.RetryPeriod)
				}
				if o.Controller.MaxConcurrentReconciles != 2 {
					t.Errorf("MaxConcurrentReconciles = %d, want 2", o.Controller.MaxConcurrentReconciles)
				}
				if o.Controller.MetricsAddress != ":9090" {
					t.Errorf("MetricsAddress = %q, want :9090", o.Controller.MetricsAddress)
				}
				if len(o.Controller.CacheNamespaces) != 1 || o.Controller.CacheNamespaces[0] != "default" {
					t.Errorf("CacheNamespaces = %v, want [default]", o.Controller.CacheNamespaces)
				}
				if o.Controller.Requeue.OnError != "100ms" || o.Controller.Requeue.NotFound != "1s" {
					t.Errorf("Requeue = %+v, want on_error 100ms not_found 1s", o.Controller.Requeue)
				}
				if !o.NodePool.RequireMatch || o.NodePool.AgentPodsPerNode != 2 {
					t.Errorf("NodePool = %+v, want require_match with 2 pods per node", o.NodePool)
				}
				if o.PersistentVolume.BufferRatio != 1.5 || o.PersistentVolume.MinSizeBytes != 1073741824 {
					t.Errorf("PersistentVolume = %+v, want ratio 1.5 min 1Gi", o.PersistentVolume)
				}
				if !o.Networking.EnableIngress || o.Networking.GatewayServiceType != "NodePort" {
					t.Errorf("Networking = %+v, want ingress enabled NodePort", o.Networking)
				}
				if cfg.Observability == nil {
					t.Error("Observability should be defaulted when omitted")
				}
			},
		},
		{
			name: "binds omitted sections to defaulted structs",
			yaml: `---
version: v0.0.0
` + minimalServerYAML + `
operator:
  name: vald-operator
`,
			check: func(t *testing.T, cfg *Data) {
				t.Helper()
				o := cfg.Operator
				if o.Controller == nil || o.Controller.LeaderElection == nil || o.Controller.Requeue == nil {
					t.Fatal("Controller section should be non-nil after Bind")
				}
				if o.Vrs == nil || o.Vrs.Agent == nil || o.Vrs.Gateway == nil || o.Vrs.Indexer == nil {
					t.Fatal("Vrs section should be non-nil after Bind")
				}
				if o.NodePool == nil || o.PersistentVolume == nil || o.Networking == nil {
					t.Fatal("NodePool/PersistentVolume/Networking should be non-nil after Bind")
				}
				if o.Vrs.LogFormat != "raw" || o.Vrs.Logger != "glg" {
					t.Errorf("Vrs logging defaults = %q/%q, want raw/glg", o.Vrs.LogFormat, o.Vrs.Logger)
				}
				if o.Vrs.ManagedGenerationLabel != "managed-generation" {
					t.Errorf("ManagedGenerationLabel = %q, want managed-generation", o.Vrs.ManagedGenerationLabel)
				}
				if o.Vrs.Agent.MaxSurge != "1" || o.Vrs.Agent.MaxUnavailable != "1" {
					t.Errorf("Agent rolling update defaults = %q/%q, want 1/1",
						o.Vrs.Agent.MaxSurge, o.Vrs.Agent.MaxUnavailable)
				}
				if o.Vrs.Agent.EnableInMemoryMode == nil || !*o.Vrs.Agent.EnableInMemoryMode {
					t.Error("Agent.EnableInMemoryMode should default to true")
				}
				if o.Vrs.Gateway.HpaTargetCPUUtilization != 80 {
					t.Errorf("Gateway.HpaTargetCPUUtilization = %d, want 80", o.Vrs.Gateway.HpaTargetCPUUtilization)
				}
				if o.Vrs.Gateway.IngressServicePort != "grpc" || o.Vrs.Gateway.IngressPathType != "Prefix" {
					t.Errorf("Gateway ingress defaults = %q/%q, want grpc/Prefix",
						o.Vrs.Gateway.IngressServicePort, o.Vrs.Gateway.IngressPathType)
				}
				if o.Vrs.Indexer.AutoIndexDurationLimit != "1h" ||
					o.Vrs.Indexer.AutoSaveIndexDurationLimit != "-1h" {
					t.Errorf("Indexer defaults = %q/%q, want 1h/-1h",
						o.Vrs.Indexer.AutoIndexDurationLimit, o.Vrs.Indexer.AutoSaveIndexDurationLimit)
				}
			},
		},
		{
			name: "returns error when operator section is missing",
			yaml: `---
version: v0.0.0
` + minimalServerYAML,
			wantErr: true,
		},
		{
			name: "returns error when server_config section is missing",
			yaml: `---
version: v0.0.0
operator:
  name: vald-operator
`,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			path := writeTempFile(t, "config.yaml", tc.yaml)
			cfg, err := New(path)
			if tc.wantErr {
				if err == nil {
					t.Error("New() error = nil, want error")
				}
				return
			}
			if err != nil {
				t.Fatalf("New() error = %v, want nil", err)
			}
			if tc.check != nil {
				tc.check(t, cfg)
			}
		})
	}
}

// TestOperatorBind_EnvExpansion cannot run in parallel because it mutates the
// process environment via t.Setenv.
func TestOperatorBind_EnvExpansion(t *testing.T) {
	t.Setenv("OPERATOR_NAME", "vald-operator-from-env")
	t.Setenv("LEADER_ELECTION_ID", "le-from-env")
	t.Setenv("VRS_PATH", "/from/env/vrs.yaml")
	t.Setenv("CACHE_NS", "ns-from-env")

	o := (&Operator{
		Name: "_OPERATOR_NAME_",
		Controller: &Controller{
			LeaderElection:  &LeaderElection{ID: "_LEADER_ELECTION_ID_"},
			CacheNamespaces: []string{"_CACHE_NS_"},
		},
		Vrs: &Vrs{DefaultVrsPath: "_VRS_PATH_"},
	}).Bind()

	if o.Name != "vald-operator-from-env" {
		t.Errorf("Name = %q, want env-expanded value", o.Name)
	}
	if o.Controller.LeaderElection.ID != "le-from-env" {
		t.Errorf("LeaderElection.ID = %q, want env-expanded value", o.Controller.LeaderElection.ID)
	}
	if o.Vrs.DefaultVrsPath != "/from/env/vrs.yaml" {
		t.Errorf("Vrs.DefaultVrsPath = %q, want env-expanded value", o.Vrs.DefaultVrsPath)
	}
	if len(o.Controller.CacheNamespaces) != 1 || o.Controller.CacheNamespaces[0] != "ns-from-env" {
		t.Errorf("CacheNamespaces = %v, want env-expanded value", o.Controller.CacheNamespaces)
	}
}

func TestOperator_Load(t *testing.T) {
	t.Parallel()

	writeVrs := func(t *testing.T) string {
		t.Helper()
		return writeTempFile(t, "vrs.yaml", `---
apiVersion: vald.vdaas.org/v1
kind: ValdRelease
metadata:
  name: vald-cluster
`)
	}

	t.Run("loads and parses the default vrs template", func(t *testing.T) {
		t.Parallel()

		path := writeVrs(t)
		o := &Operator{
			Vrs: &Vrs{
				DefaultVrsPath: path,
				LogLevel:       "warn",
			},
			NodePool: &NodePool{AgentPodsPerNode: 3},
		}
		cfg, err := o.Load()
		if err != nil {
			t.Fatalf("Load() error = %v, want nil", err)
		}
		if cfg.DefaultVrs.Us == nil || cfg.DefaultVrs.Us.GetKind() != "ValdRelease" {
			t.Errorf("DefaultVrs.Us kind = %v, want ValdRelease", cfg.DefaultVrs.Us)
		}
		if cfg.AgentPodsPerNode != 3 {
			t.Errorf("AgentPodsPerNode = %d, want 3", cfg.AgentPodsPerNode)
		}
		if cfg.GatewayIngressAnnotations == nil {
			t.Error("GatewayIngressAnnotations should be non-nil")
		}
		if cfg.VrsLogFormat != "raw" || cfg.VrsLogger != "glg" {
			t.Errorf("Vrs logging defaults = %q/%q, want raw/glg", cfg.VrsLogFormat, cfg.VrsLogger)
		}
		if cfg.ManagedGenerationLabel != "managed-generation" {
			t.Errorf("ManagedGenerationLabel = %q, want managed-generation", cfg.ManagedGenerationLabel)
		}
		if !cfg.AgentEnableInMemoryMode {
			t.Error("AgentEnableInMemoryMode should default to true")
		}
		if cfg.GatewayHpaTargetCPUUtilization != 80 {
			t.Errorf("GatewayHpaTargetCPUUtilization = %d, want 80", cfg.GatewayHpaTargetCPUUtilization)
		}
	})

	t.Run("parses the requeue intervals into durations", func(t *testing.T) {
		t.Parallel()

		o := &Operator{
			Vrs: &Vrs{DefaultVrsPath: writeVrs(t)},
			Controller: &Controller{
				MaxConcurrentReconciles: 4,
				Requeue: &Requeue{
					Success:  "5m",
					OnError:  "100ms",
					NotFound: "1s",
				},
			},
		}
		cfg, err := o.Load()
		if err != nil {
			t.Fatalf("Load() error = %v, want nil", err)
		}
		if cfg.RequeueAfterSuccess != 5*time.Minute {
			t.Errorf("RequeueAfterSuccess = %v, want 5m", cfg.RequeueAfterSuccess)
		}
		if cfg.RequeueAfterError != 100*time.Millisecond {
			t.Errorf("RequeueAfterError = %v, want 100ms", cfg.RequeueAfterError)
		}
		if cfg.RequeueAfterNotFound != time.Second {
			t.Errorf("RequeueAfterNotFound = %v, want 1s", cfg.RequeueAfterNotFound)
		}
		if cfg.MaxConcurrentReconciles != 4 {
			t.Errorf("MaxConcurrentReconciles = %d, want 4", cfg.MaxConcurrentReconciles)
		}
	})

	t.Run("returns error when a requeue interval is invalid", func(t *testing.T) {
		t.Parallel()

		o := &Operator{
			Vrs: &Vrs{DefaultVrsPath: writeVrs(t)},
			Controller: &Controller{
				Requeue: &Requeue{OnError: "not-a-duration"},
			},
		}
		if _, err := o.Load(); err == nil {
			t.Error("Load() error = nil, want error")
		}
	})

	t.Run("returns error when the template file does not exist", func(t *testing.T) {
		t.Parallel()

		o := &Operator{
			Vrs: &Vrs{DefaultVrsPath: filepath.Join(t.TempDir(), "missing.yaml")},
		}
		if _, err := o.Load(); err == nil {
			t.Error("Load() error = nil, want error")
		}
	})

	t.Run("returns error when the template is invalid yaml", func(t *testing.T) {
		t.Parallel()

		path := writeTempFile(t, "vrs.yaml", "\tinvalid: [yaml")
		o := &Operator{Vrs: &Vrs{DefaultVrsPath: path}}
		if _, err := o.Load(); err == nil {
			t.Error("Load() error = nil, want error")
		}
	})
}

// NOT IMPLEMENTED BELOW
//
// func TestOperator_Bind(t *testing.T) {
// 	type fields struct {
// 		Controller       *Controller
// 		Vrs              *Vrs
// 		NodePool         *NodePool
// 		PersistentVolume *PersistentVolume
// 		Networking       *Networking
// 		Name             string
// 		Namespace        string
// 	}
// 	type want struct {
// 		want *Operator
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *Operator) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got *Operator) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           Controller:Controller{},
// 		           Vrs:Vrs{},
// 		           NodePool:NodePool{},
// 		           PersistentVolume:PersistentVolume{},
// 		           Networking:Networking{},
// 		           Name:"",
// 		           Namespace:"",
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           Controller:Controller{},
// 		           Vrs:Vrs{},
// 		           NodePool:NodePool{},
// 		           PersistentVolume:PersistentVolume{},
// 		           Networking:Networking{},
// 		           Name:"",
// 		           Namespace:"",
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			o := &Operator{
// 				Controller:       test.fields.Controller,
// 				Vrs:              test.fields.Vrs,
// 				NodePool:         test.fields.NodePool,
// 				PersistentVolume: test.fields.PersistentVolume,
// 				Networking:       test.fields.Networking,
// 				Name:             test.fields.Name,
// 				Namespace:        test.fields.Namespace,
// 			}
//
// 			got := o.Bind()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestOperator_bindController(t *testing.T) {
// 	type fields struct {
// 		Controller       *Controller
// 		Vrs              *Vrs
// 		NodePool         *NodePool
// 		PersistentVolume *PersistentVolume
// 		Networking       *Networking
// 		Name             string
// 		Namespace        string
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           Controller:Controller{},
// 		           Vrs:Vrs{},
// 		           NodePool:NodePool{},
// 		           PersistentVolume:PersistentVolume{},
// 		           Networking:Networking{},
// 		           Name:"",
// 		           Namespace:"",
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           Controller:Controller{},
// 		           Vrs:Vrs{},
// 		           NodePool:NodePool{},
// 		           PersistentVolume:PersistentVolume{},
// 		           Networking:Networking{},
// 		           Name:"",
// 		           Namespace:"",
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			o := &Operator{
// 				Controller:       test.fields.Controller,
// 				Vrs:              test.fields.Vrs,
// 				NodePool:         test.fields.NodePool,
// 				PersistentVolume: test.fields.PersistentVolume,
// 				Networking:       test.fields.Networking,
// 				Name:             test.fields.Name,
// 				Namespace:        test.fields.Namespace,
// 			}
//
// 			o.bindController()
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestOperator_bindVrs(t *testing.T) {
// 	type fields struct {
// 		Controller       *Controller
// 		Vrs              *Vrs
// 		NodePool         *NodePool
// 		PersistentVolume *PersistentVolume
// 		Networking       *Networking
// 		Name             string
// 		Namespace        string
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           Controller:Controller{},
// 		           Vrs:Vrs{},
// 		           NodePool:NodePool{},
// 		           PersistentVolume:PersistentVolume{},
// 		           Networking:Networking{},
// 		           Name:"",
// 		           Namespace:"",
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           Controller:Controller{},
// 		           Vrs:Vrs{},
// 		           NodePool:NodePool{},
// 		           PersistentVolume:PersistentVolume{},
// 		           Networking:Networking{},
// 		           Name:"",
// 		           Namespace:"",
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			o := &Operator{
// 				Controller:       test.fields.Controller,
// 				Vrs:              test.fields.Vrs,
// 				NodePool:         test.fields.NodePool,
// 				PersistentVolume: test.fields.PersistentVolume,
// 				Networking:       test.fields.Networking,
// 				Name:             test.fields.Name,
// 				Namespace:        test.fields.Namespace,
// 			}
//
// 			o.bindVrs()
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestOperator_bindNodePool(t *testing.T) {
// 	type fields struct {
// 		Controller       *Controller
// 		Vrs              *Vrs
// 		NodePool         *NodePool
// 		PersistentVolume *PersistentVolume
// 		Networking       *Networking
// 		Name             string
// 		Namespace        string
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           Controller:Controller{},
// 		           Vrs:Vrs{},
// 		           NodePool:NodePool{},
// 		           PersistentVolume:PersistentVolume{},
// 		           Networking:Networking{},
// 		           Name:"",
// 		           Namespace:"",
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           Controller:Controller{},
// 		           Vrs:Vrs{},
// 		           NodePool:NodePool{},
// 		           PersistentVolume:PersistentVolume{},
// 		           Networking:Networking{},
// 		           Name:"",
// 		           Namespace:"",
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			o := &Operator{
// 				Controller:       test.fields.Controller,
// 				Vrs:              test.fields.Vrs,
// 				NodePool:         test.fields.NodePool,
// 				PersistentVolume: test.fields.PersistentVolume,
// 				Networking:       test.fields.Networking,
// 				Name:             test.fields.Name,
// 				Namespace:        test.fields.Namespace,
// 			}
//
// 			o.bindNodePool()
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestOperator_bindPersistentVolume(t *testing.T) {
// 	type fields struct {
// 		Controller       *Controller
// 		Vrs              *Vrs
// 		NodePool         *NodePool
// 		PersistentVolume *PersistentVolume
// 		Networking       *Networking
// 		Name             string
// 		Namespace        string
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           Controller:Controller{},
// 		           Vrs:Vrs{},
// 		           NodePool:NodePool{},
// 		           PersistentVolume:PersistentVolume{},
// 		           Networking:Networking{},
// 		           Name:"",
// 		           Namespace:"",
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           Controller:Controller{},
// 		           Vrs:Vrs{},
// 		           NodePool:NodePool{},
// 		           PersistentVolume:PersistentVolume{},
// 		           Networking:Networking{},
// 		           Name:"",
// 		           Namespace:"",
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			o := &Operator{
// 				Controller:       test.fields.Controller,
// 				Vrs:              test.fields.Vrs,
// 				NodePool:         test.fields.NodePool,
// 				PersistentVolume: test.fields.PersistentVolume,
// 				Networking:       test.fields.Networking,
// 				Name:             test.fields.Name,
// 				Namespace:        test.fields.Namespace,
// 			}
//
// 			o.bindPersistentVolume()
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestOperator_bindNetworking(t *testing.T) {
// 	type fields struct {
// 		Controller       *Controller
// 		Vrs              *Vrs
// 		NodePool         *NodePool
// 		PersistentVolume *PersistentVolume
// 		Networking       *Networking
// 		Name             string
// 		Namespace        string
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           Controller:Controller{},
// 		           Vrs:Vrs{},
// 		           NodePool:NodePool{},
// 		           PersistentVolume:PersistentVolume{},
// 		           Networking:Networking{},
// 		           Name:"",
// 		           Namespace:"",
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           Controller:Controller{},
// 		           Vrs:Vrs{},
// 		           NodePool:NodePool{},
// 		           PersistentVolume:PersistentVolume{},
// 		           Networking:Networking{},
// 		           Name:"",
// 		           Namespace:"",
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			o := &Operator{
// 				Controller:       test.fields.Controller,
// 				Vrs:              test.fields.Vrs,
// 				NodePool:         test.fields.NodePool,
// 				PersistentVolume: test.fields.PersistentVolume,
// 				Networking:       test.fields.Networking,
// 				Name:             test.fields.Name,
// 				Namespace:        test.fields.Namespace,
// 			}
//
// 			o.bindNetworking()
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }

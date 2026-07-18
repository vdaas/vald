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

package service

import (
	"context"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/k8s/resource"
	"github.com/vdaas/vald/internal/k8s/vald/operator/api/metadata"
	v1 "github.com/vdaas/vald/internal/k8s/vald/operator/api/v1"
	"github.com/vdaas/vald/internal/k8s/vald/operator/api/valdrelease"
	mock "github.com/vdaas/vald/internal/test/mock/k8s"
	"github.com/vdaas/vald/pkg/operator/vald/config"
)

const (
	testIngressHost         = "foo.example.com"
	testServiceTypeLoadBlnc = "LoadBalancer"
	testVrNamespaceA        = "ns-a"
	testVrName1             = "vr-1"
	testVrName2             = "vr-2"
)

func loadValdOperatorReleaseFromYAML(tb testing.TB) *v1.ValdOperatorRelease {
	tb.Helper()
	data, err := os.ReadFile("testdata/vor.yaml") //nolint:gosec
	assert.NoError(tb, err)

	var cr v1.ValdOperatorRelease
	err = k8s.YAMLUnmarshal(data, &cr)
	assert.NoError(tb, err)
	return &cr
}

//nolint:goconst
func TestVrsBuilder_Validate(t *testing.T) {
	validCR := func() *v1.ValdOperatorRelease {
		return &v1.ValdOperatorRelease{
			Spec: v1.ValdOperatorReleaseSpec{
				Infrastructure: []v1.ValdOperatorReleaseInfra{
					{
						Role: testRoleGreen,
						Clusters: []v1.DestClusters{
							{ID: "abc-123", Name: "cluster-a"},
						},
					},
				},
			},
		}
	}

	tests := []struct {
		cr      *v1.ValdOperatorRelease
		name    string
		wantErr bool
	}{
		{
			name:    "nil CR",
			cr:      nil,
			wantErr: true,
		},
		{
			name: "empty infrastructure",
			cr: &v1.ValdOperatorRelease{
				Spec: v1.ValdOperatorReleaseSpec{Infrastructure: nil},
			},
			wantErr: true,
		},
		{
			name: "infra with empty clusters",
			cr: &v1.ValdOperatorRelease{
				Spec: v1.ValdOperatorReleaseSpec{
					Infrastructure: []v1.ValdOperatorReleaseInfra{
						{Role: testRoleGreen, Clusters: []v1.DestClusters{}},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "cluster with empty Name",
			cr: func() *v1.ValdOperatorRelease {
				cr := validCR()
				cr.Spec.Infrastructure[0].Clusters[0].Name = ""
				return cr
			}(),
			wantErr: true,
		},
		{
			name:    "valid CR",
			cr:      validCR(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &vrsBuilder{cr: tt.cr}
			err := b.validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestVrsBuilder_Build(t *testing.T) {
	// Initialize config for testing
	cfg := initConfig(t)

	cr := loadValdOperatorReleaseFromYAML(t)

	b := newVrsBuilder(cr, cfg, alwaysAvailable())
	objList, err := b.Build(context.Background())
	assert.NoError(t, err)

	// path to the golden file
	goldenPath := "testdata/vrs.golden.yaml"

	// number of generated items
	assert.NotNil(t, objList)
	items := objList
	assert.Equal(t, 4, len(items))

	topItem := items[0].(*valdrelease.ValdRelease)
	assert.Equal(t, cr.Namespace, topItem.GetNamespace())

	itemData, err := k8s.YAMLMarshal(topItem)
	assert.NoError(t, err)
	// Regenerate the golden file only when explicitly requested:
	//   UPDATE_GOLDEN=1 go test ./pkg/operator/vald/service/...
	if os.Getenv("UPDATE_GOLDEN") != "" {
		assert.NoError(t, os.WriteFile(goldenPath, itemData, 0o644)) //nolint:gosec
	}
	expected, err := os.ReadFile(goldenPath)
	assert.NoError(t, err)

	// Compare parsed structures instead of raw bytes: the repository format
	// pipeline (license header + prettier) restyles the golden file, so
	// byte-equality would break on formatting-only differences.
	var want, got map[string]any
	assert.NoError(t, k8s.YAMLUnmarshal(expected, &want))
	assert.NoError(t, k8s.YAMLUnmarshal(itemData, &got))
	assert.Equal(t, want, got)
}

func TestVrsBuilder_Build_ManagedLabelsPreserved(t *testing.T) {
	cfg := initConfig(t)
	cr := loadValdOperatorReleaseFromYAML(t)

	items, err := newVrsBuilder(cr, cfg, alwaysAvailable()).Build(context.Background())
	assert.NoError(t, err)
	assert.NotEmpty(t, items)

	for _, obj := range items {
		labels := obj.GetLabels()
		assert.Equal(t, v1.GroupVersion.Group, labels[metadata.ManagedByLabel],
			"%s: managed-by label must survive the per-infra label merge", obj.GetName())
		assert.Equal(t, valdrelease.GVK.Kind,
			labels[v1.GroupVersion.Group+metadata.SubResourceLabelSuffix],
			"%s: managed-resource label must survive the per-infra label merge", obj.GetName())
		assert.NotEmpty(t, labels[nodePoolLabelType], "%s: per-infra type label must be present", obj.GetName())
		assert.NotEmpty(t, labels[nodePoolLabelRole], "%s: per-infra role label must be present", obj.GetName())
	}
}

func TestVrsBuilder_Build_OverlayAppliedPerCluster(t *testing.T) {
	cfg := initConfig(t)
	cr := loadValdOperatorReleaseFromYAML(t)

	items, err := newVrsBuilder(cr, cfg, alwaysAvailable()).Build(context.Background())
	assert.NoError(t, err)
	assert.Len(t, items, 4)

	names := make(map[string]bool, len(items))
	for _, obj := range items {
		vr, ok := obj.(*valdrelease.ValdRelease)
		assert.True(t, ok)
		names[vr.GetName()] = true
		// vor.yaml overlays agent.ngt.dimension=256 over the spec value 128;
		// the shared parsed overlay must reach every cluster.
		assert.Equal(t, 256, *vr.Spec.Agent.Ngt.Dimension, "%s: overlay must apply", vr.GetName())
	}
	assert.Len(t, names, 4, "every cluster must keep its own name")
}

// Merged test components

//nolint:goconst
func minimalVor() *v1.ValdOperatorRelease {
	return &v1.ValdOperatorRelease{
		ObjectMeta: k8s.ObjectMeta{
			Name:      "test",
			Namespace: "test-ns",
		},
		Spec: v1.ValdOperatorReleaseSpec{
			Infrastructure: []v1.ValdOperatorReleaseInfra{
				{
					Role:   testRoleGreen,
					Active: true,
					Clusters: []v1.DestClusters{
						{Name: "cluster-1"},
					},
					NodePools: v1.NodePools{
						v1.NodePoolTypeGeneral: v1.NodePool{
							Replicas:        1,
							MachineResource: v1.MachineResource{Cpu: "2", Memory: "4Gi"},
						},
					},
				},
			},
			VectorEngine: v1.VectorEngine{
				Vald: v1.Vald{
					Agent: v1.Agent{
						Ngt: v1.Ngt{Dimension: 128},
					},
				},
			},
		},
	}
}

// initConfig loads the runtime Config from the same default operator settings
// that the former env-based configuration provided, so the golden assets stay
// stable. Builder values omitted here (log format, HPA target, ingress port,
// indexer limits, agent rolling-update) come from the Bind defaults.
func initConfig(tb testing.TB) *config.Config {
	tb.Helper()
	return initConfigWith(tb, nil)
}

// initConfigWith loads the runtime Config after applying mutate to the base
// operator settings used by initConfig, so tests can tweak single knobs
// (networking, node pool, ...) without repeating the whole fixture.
func initConfigWith(tb testing.TB, mutate func(o *config.Operator)) *config.Config {
	tb.Helper()
	op := &config.Operator{
		Vrs: &config.Vrs{
			DefaultVrsPath: "testdata/vrs.yaml",
			LogLevel:       "warn",
		},
		NodePool: &config.NodePool{
			AgentPodsPerNode: 2,
		},
		PersistentVolume: &config.PersistentVolume{
			DefaultStorageClass: "standard",
			DefaultAccessMode:   "ReadWriteOnce",
			BufferRatio:         1.5,
			MinSizeBytes:        1073741824,
		},
		Networking: &config.Networking{
			EnableIngress:                     true,
			GatewayServiceType:                "NodePort",
			DiscovererDaemonSetMaxSurge:       "30%",
			DiscovererDaemonSetMaxUnavailable: "0%",
		},
	}
	if mutate != nil {
		mutate(op)
	}
	cfg, err := op.Load()
	if err != nil {
		tb.Fatalf("config load: %v", err)
	}
	return cfg
}

func newTestBuilder(cr *v1.ValdOperatorRelease, cfg *config.Config) *vrsBuilder {
	return newVrsBuilder(cr, cfg, alwaysAvailable())
}

func TestBuildLb_IngressEnabled(t *testing.T) {
	cfg := initConfig(t)
	cr := minimalVor()
	cr.Spec.VectorEngine.Vald.Gateway.Ingress = &v1.GatewayIngress{
		Enabled: true,
		Host:    testIngressHost,
	}
	lb := newTestBuilder(cr, cfg).buildLb()

	assert.True(t, *lb.Ingress.Enabled)
	assert.Equal(t, testIngressHost, *lb.Ingress.Host)
}

func TestBuildLb_IngressDisabledWhenNil(t *testing.T) {
	cfg := initConfig(t)
	cr := minimalVor()
	lb := newTestBuilder(cr, cfg).buildLb()

	assert.Nil(t, lb.Ingress.Enabled)
	assert.Empty(t, lb.Ingress.Host)
}

func TestBuildLb_ServiceTypeDefault(t *testing.T) {
	cfg := initConfig(t)
	cr := minimalVor()
	lb := newTestBuilder(cr, cfg).buildLb()

	assert.Equal(t, valdrelease.GatewayLbServiceTypeNodePort, *lb.ServiceType)
}

func TestBuildLb_ServiceTypeLoadBalancer(t *testing.T) {
	cfg := initConfig(t)
	cr := minimalVor()
	cr.Spec.VectorEngine.Vald.Gateway.ServiceType = testServiceTypeLoadBlnc
	lb := newTestBuilder(cr, cfg).buildLb()

	assert.Equal(t, valdrelease.GatewayLbServiceTypeLoadBalancer, *lb.ServiceType)
}

func TestBuildLb_ServiceTypeClusterIP(t *testing.T) {
	cfg := initConfig(t)
	cr := minimalVor()
	cr.Spec.VectorEngine.Vald.Gateway.ServiceType = "ClusterIP"
	lb := newTestBuilder(cr, cfg).buildLb()

	assert.Equal(t, valdrelease.GatewayLbServiceTypeClusterIP, *lb.ServiceType)
}

func TestBuildLb_ServiceTypeFromConfig(t *testing.T) {
	cfg := initConfigWith(t, func(o *config.Operator) {
		o.Networking.GatewayServiceType = testServiceTypeLoadBlnc
	})
	cr := minimalVor() // CR spec leaves ServiceType empty

	lb := newTestBuilder(cr, cfg).buildLb()

	assert.Equal(t, valdrelease.GatewayLbServiceTypeLoadBalancer, *lb.ServiceType,
		"operator config must apply when the CR omits the service type")
}

func TestBuildLb_ServiceTypeCRWinsOverConfig(t *testing.T) {
	cfg := initConfigWith(t, func(o *config.Operator) {
		o.Networking.GatewayServiceType = testServiceTypeLoadBlnc
	})
	cr := minimalVor()
	cr.Spec.VectorEngine.Vald.Gateway.ServiceType = "ClusterIP"

	lb := newTestBuilder(cr, cfg).buildLb()

	assert.Equal(t, valdrelease.GatewayLbServiceTypeClusterIP, *lb.ServiceType,
		"CR spec must win over the operator config")
}

func TestBuildIngress_ConfigGateDisablesCREnabled(t *testing.T) {
	cfg := initConfigWith(t, func(o *config.Operator) {
		o.Networking.EnableIngress = false
	})
	cr := minimalVor()
	cr.Spec.VectorEngine.Vald.Gateway.Ingress = &v1.GatewayIngress{
		Enabled: true,
		Host:    testIngressHost,
	}

	lb := newTestBuilder(cr, cfg).buildLb()

	assert.Nil(t, lb.Ingress.Enabled,
		"networking.enable_ingress=false must keep the ingress disabled even when the CR enables it")
	assert.Empty(t, lb.Ingress.Host)
}

func TestBuildIngress_AnnotationsFromConfig(t *testing.T) {
	annotations := map[string]string{
		"kubernetes.io/ingress.class":                  "nginx",
		"nginx.ingress.kubernetes.io/backend-protocol": "GRPC",
	}
	cfg := initConfigWith(t, func(o *config.Operator) {
		o.Networking.GatewayIngressAnnotations = annotations
	})
	cr := minimalVor()
	cr.Spec.VectorEngine.Vald.Gateway.Ingress = &v1.GatewayIngress{
		Enabled: true,
		Host:    testIngressHost,
	}

	lb := newTestBuilder(cr, cfg).buildLb()

	assert.Equal(t, toAnyMap(annotations), *lb.Ingress.Annotations,
		"configured annotations must be reflected on the built ingress")

	// Without configured annotations the field stays nil (omitempty).
	lb = newTestBuilder(minimalVor(), initConfig(t)).buildLb()
	assert.Nil(t, lb.Ingress.Annotations)
}

func newAffinityBuilder(capability nodePoolCapability) *vrsBuilder {
	cr := &v1.ValdOperatorRelease{}
	cr.SetNamespace("ns-affinity")
	return &vrsBuilder{
		cr:         cr,
		cfg:        &config.Config{NodePoolLabelPrefix: "vald.vdaas.org"},
		capability: capability,
	}
}

func TestVrsBuilder_ApplyNodeAffinities(t *testing.T) {
	const (
		nsKey   = "vald.vdaas.org/namespace"
		typeKey = "vald.vdaas.org/type"
	)

	assertGeneral := func(t *testing.T, name string, ns *valdrelease.NodeSelector) {
		t.Helper()
		assert.NotNil(t, ns, name+" NodeSelector set")
		assert.Equal(t, "ns-affinity", (*ns)[nsKey], name+" NodeSelector namespace")
		assert.Equal(t, "general", (*ns)[typeKey], name+" NodeSelector type")
	}

	// skeleton returns a ValdRelease whose component sub-specs are non-nil, as
	// they always are after Build() populates the generated Values. The
	// generated types use pointers throughout, so applyNodeAffinities needs
	// these initialized to avoid nil dereferences.
	skeleton := func() *valdrelease.ValdRelease {
		return &valdrelease.ValdRelease{Spec: valdrelease.Values{
			Agent:      &valdrelease.Agent{},
			Gateway:    &valdrelease.Gateway{Lb: &valdrelease.GatewayLb{}},
			Discoverer: &valdrelease.Discoverer{},
			Manager:    &valdrelease.Manager{Index: &valdrelease.ManagerIndex{}},
		}}
	}

	tests := []struct {
		row        func() *valdrelease.ValdRelease
		check      func(t *testing.T, row *valdrelease.ValdRelease)
		name       string
		capability nodePoolCapability
	}{
		{
			name:       "all components: agent uses agent pool, others use general pool",
			capability: alwaysAvailable(),
			row: func() *valdrelease.ValdRelease {
				row := skeleton()
				row.Spec.Manager.Index.Saver = &valdrelease.ManagerIndexSaver{}
				row.Spec.Manager.Index.Creator = &valdrelease.ManagerIndexCreator{}
				return row
			},
			check: func(t *testing.T, row *valdrelease.ValdRelease) {
				t.Helper()
				// Agent: agent pool (because capability says HasAgentPool=true)
				assert.Equal(t, "agent", (*row.Spec.Agent.NodeSelector)[typeKey], "agent uses agent pool")
				assert.NotNil(t, row.Spec.Agent.Tolerations, "agent tolerations set")

				// All others: general pool
				assertGeneral(t, "gateway.Lb", row.Spec.Gateway.Lb.NodeSelector)
				assertGeneral(t, "discoverer", row.Spec.Discoverer.NodeSelector)
				assertGeneral(t, "manager.Index", row.Spec.Manager.Index.NodeSelector)
				assertGeneral(t, "manager.Index.Saver", row.Spec.Manager.Index.Saver.NodeSelector)
				assertGeneral(t, "manager.Index.Creator", row.Spec.Manager.Index.Creator.NodeSelector)

				assert.NotNil(t, row.Spec.Gateway.Lb.Tolerations)
				assert.NotNil(t, row.Spec.Discoverer.Tolerations)
				assert.NotNil(t, row.Spec.Manager.Index.Tolerations)
				assert.NotNil(t, row.Spec.Manager.Index.Saver.Tolerations)
				assert.NotNil(t, row.Spec.Manager.Index.Creator.Tolerations)
			},
		},
		{
			// capability declares the cluster has no dedicated agent pool.
			name:       "no agent pool: agent falls back to general pool",
			capability: nodePoolCapability{HasGeneralPool: true, HasAgentPool: false},
			row:        skeleton,
			check: func(t *testing.T, row *valdrelease.ValdRelease) {
				t.Helper()
				assert.Equal(t, "general", (*row.Spec.Agent.NodeSelector)[typeKey],
					"agent must fall back to general pool when no agent pool exists")
			},
		},
		{
			// Saver and Creator are nil (Index.Enabled=true path uses Indexer
			// instead). Must not panic; non-nil components still get values.
			name:       "nil saver/creator: no panic, non-nil components still set",
			capability: alwaysAvailable(),
			row:        skeleton,
			check: func(t *testing.T, row *valdrelease.ValdRelease) {
				t.Helper()
				assert.Nil(t, row.Spec.Manager.Index.Saver)
				assert.Nil(t, row.Spec.Manager.Index.Creator)
				assert.NotNil(t, row.Spec.Agent.Tolerations)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := newAffinityBuilder(tt.capability)
			row := tt.row()

			b.applyNodeAffinities(row)

			tt.check(t, row)
		})
	}
}

//nolint:unparam
func newAgentReleaseWithMemory(memory string) *valdrelease.ValdRelease {
	r := &valdrelease.ValdRelease{Spec: valdrelease.Values{Agent: &valdrelease.Agent{}}}
	r.Spec.Agent.Resources = &k8s.ResourceRequirements{
		Requests: k8s.ResourceList{
			k8s.ResourceMemory: resource.MustParse(memory),
		},
	}
	return r
}

func TestVrsBuilder_ReflectPersistentVolume_Disabled(t *testing.T) {
	b := &vrsBuilder{cr: &v1.ValdOperatorRelease{
		Spec: v1.ValdOperatorReleaseSpec{
			VectorEngine: v1.VectorEngine{
				Vald: v1.Vald{
					Agent: v1.Agent{PersistentVolume: &v1.AgentPersistentVolume{Enabled: false}},
				},
			},
		},
	}}
	r := newAgentReleaseWithMemory("4Gi")
	b.reflectPersistentVolume(r)
	assert.Nil(t, r.Spec.Agent.PersistentVolume, "PV must not be set when Enabled=false")
}

func TestVrsBuilder_ReflectPersistentVolume_Nil(t *testing.T) {
	b := &vrsBuilder{cr: &v1.ValdOperatorRelease{
		Spec: v1.ValdOperatorReleaseSpec{
			VectorEngine: v1.VectorEngine{
				Vald: v1.Vald{
					Agent: v1.Agent{PersistentVolume: nil},
				},
			},
		},
	}}
	r := newAgentReleaseWithMemory("4Gi")
	b.reflectPersistentVolume(r)
	assert.Nil(t, r.Spec.Agent.PersistentVolume, "PV must not be set when CR PersistentVolume is nil")
}

func TestVrsBuilder_ReflectPersistentVolume_FromCR(t *testing.T) {
	cfg := &config.Config{
		DefaultStorageClass: "fallback-sc",
		DefaultAccessMode:   "fallback-am",
		PvBufferRatio:       1.5,
		PvMinSizeBytes:      int64(1) << 30,
	}

	b := &vrsBuilder{cfg: cfg, cr: &v1.ValdOperatorRelease{
		Spec: v1.ValdOperatorReleaseSpec{
			VectorEngine: v1.VectorEngine{
				Vald: v1.Vald{
					Agent: v1.Agent{PersistentVolume: &v1.AgentPersistentVolume{
						Enabled:      true,
						StorageClass: "from-cr-sc",
						AccessMode:   "from-cr-am",
					}},
				},
			},
		},
	}}
	r := newAgentReleaseWithMemory("4Gi")
	b.reflectPersistentVolume(r)

	assert.NotNil(t, r.Spec.Agent.PersistentVolume)
	assert.True(t, *r.Spec.Agent.PersistentVolume.Enabled)
	assert.Equal(t, "from-cr-sc", *r.Spec.Agent.PersistentVolume.StorageClass, "CR value must win")
	assert.Equal(t, "from-cr-am", *r.Spec.Agent.PersistentVolume.AccessMode, "CR value must win")
}

func TestVrsBuilder_ReflectPersistentVolume_ConfigFallback(t *testing.T) {
	cfg := &config.Config{
		DefaultStorageClass: "fallback-sc",
		DefaultAccessMode:   "fallback-am",
		PvBufferRatio:       1.5,
		PvMinSizeBytes:      int64(1) << 30,
	}

	b := &vrsBuilder{cfg: cfg, cr: &v1.ValdOperatorRelease{
		Spec: v1.ValdOperatorReleaseSpec{
			VectorEngine: v1.VectorEngine{
				Vald: v1.Vald{
					Agent: v1.Agent{PersistentVolume: &v1.AgentPersistentVolume{
						Enabled: true,
					}},
				},
			},
		},
	}}
	r := newAgentReleaseWithMemory("4Gi")
	b.reflectPersistentVolume(r)

	assert.NotNil(t, r.Spec.Agent.PersistentVolume)
	assert.Equal(t, "fallback-sc", *r.Spec.Agent.PersistentVolume.StorageClass, "config fallback when CR omits SC")
	assert.Equal(t, "fallback-am", *r.Spec.Agent.PersistentVolume.AccessMode, "config fallback when CR omits AM")
}

// --- existing VRS fetch ---

// schemeOverrideClient wraps a working client but reports a scheme that cannot
// resolve ValdRelease, simulating the failure mode of the generic GVK
// auto-restoration in resource.ListObjects (the type unregistered, or the same
// Go type registered under multiple GVKs after a CRD API-version bump).
type schemeOverrideClient struct {
	k8s.Client
	scheme *k8s.Scheme
}

func (c *schemeOverrideClient) Scheme() *k8s.Scheme { return c.scheme }

func TestFetchExistingVrs_RestoresGVK(t *testing.T) {
	newVr := func(name, namespace string) *valdrelease.ValdRelease {
		vr := &valdrelease.ValdRelease{}
		vr.SetName(name)
		vr.SetNamespace(namespace)
		return vr
	}

	tests := []struct {
		name      string
		namespace string
		objs      []k8s.Object
		wantNames []string
		// unresolvableScheme makes the client report a scheme that cannot
		// resolve ValdRelease, so the generic auto-restoration skips and only
		// the static valdrelease.GVK fallback can fill the item GVKs.
		unresolvableScheme bool
	}{
		{
			name:      "restores the stripped GVK on every listed object",
			namespace: testVrNamespaceA,
			objs:      []k8s.Object{newVr(testVrName1, testVrNamespaceA), newVr(testVrName2, testVrNamespaceA)},
			wantNames: []string{testVrName1, testVrName2},
		},
		{
			name:      "namespace scoping excludes foreign objects",
			namespace: testVrNamespaceA,
			objs:      []k8s.Object{newVr(testVrName1, testVrNamespaceA), newVr("vr-3", "ns-b")},
			wantNames: []string{testVrName1},
		},
		{
			name:      "empty namespace yields no objects",
			namespace: "ns-empty",
			objs:      []k8s.Object{newVr(testVrName1, testVrNamespaceA)},
			wantNames: []string{},
		},
		{
			name:               "static fallback stamps the GVK when the scheme cannot resolve the type",
			namespace:          testVrNamespaceA,
			objs:               []k8s.Object{newVr(testVrName1, testVrNamespaceA), newVr(testVrName2, testVrNamespaceA)},
			wantNames:          []string{testVrName1, testVrName2},
			unresolvableScheme: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var c k8s.Client = mock.NewFakeClientBuilder().WithScheme(newTestScheme(t)).WithObjects(tc.objs...).Build()
			if tc.unresolvableScheme {
				c = &schemeOverrideClient{Client: c, scheme: k8s.NewScheme()}
			}

			out, err := fetchExistingVrs(context.Background(), c, tc.namespace)
			assert.NoError(t, err)

			names := make([]string, 0, len(out))
			for _, obj := range out {
				names = append(names, obj.GetName())
				assert.Equal(t, valdrelease.GVK, obj.GetObjectKind().GroupVersionKind(),
					"%s: the GVK stripped by the client must be restored for syncer prune keys", obj.GetName())
			}
			assert.ElementsMatch(t, tc.wantNames, names)
		})
	}
}

// --- overlay parse / clone ---

func newOverlayBuilder(raw []byte) *vrsBuilder {
	cr := &v1.ValdOperatorRelease{}
	cr.Spec.VectorEngine.Vald.Overlay = v1.JSON{Raw: raw}
	return &vrsBuilder{cr: cr}
}

func TestVrsBuilder_ParseOverlay(t *testing.T) {
	tests := []struct {
		name    string
		raw     []byte
		wantNil bool
		wantErr bool
	}{
		{
			name:    "empty raw leaves the overlay nil",
			raw:     nil,
			wantNil: true,
		},
		{
			name:    "invalid json returns an error",
			raw:     []byte(`{invalid`),
			wantErr: true,
		},
		{
			name: "valid raw is decoded once with the builder GVK",
			raw:  []byte(`{"spec":{"agent":{"ngt":{"dimension":256}}}}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := newOverlayBuilder(tt.raw)
			err := b.parseOverlay()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			if tt.wantNil {
				assert.Nil(t, b.overlay)
				assert.Nil(t, b.makeOverlayPatch(&valdrelease.ValdRelease{}))
				return
			}
			assert.NotNil(t, b.overlay)
			assert.Equal(t, valdrelease.GVK, b.overlay.GroupVersionKind())
		})
	}
}

func TestVrsBuilder_MakeOverlayPatch_ClonesParsedOverlay(t *testing.T) {
	b := newOverlayBuilder([]byte(`{"spec":{"agent":{"ngt":{"dimension":256}}}}`))
	assert.NoError(t, b.parseOverlay())

	row1 := &valdrelease.ValdRelease{}
	row1.SetName("cluster-1")
	row1.SetNamespace("ns-1")
	patch1 := b.makeOverlayPatch(row1)
	assert.NotNil(t, patch1)
	assert.Equal(t, "cluster-1", patch1.GetName())
	assert.Equal(t, "ns-1", patch1.GetNamespace())

	// Mutating one clone must not leak into the cached overlay or later clones.
	patch1.Object["spec"] = "mutated"

	row2 := &valdrelease.ValdRelease{}
	row2.SetName("cluster-2")
	row2.SetNamespace("ns-2")
	patch2 := b.makeOverlayPatch(row2)
	assert.NotNil(t, patch2)
	assert.Equal(t, "cluster-2", patch2.GetName())
	assert.Equal(t, "ns-2", patch2.GetNamespace())

	// json.Unmarshal decodes numbers into float64 inside map[string]any.
	dim, found, err := k8s.NestedFloat64(patch2.Object, "spec", "agent", "ngt", "dimension")
	assert.NoError(t, err)
	assert.True(t, found, "the cached overlay must stay intact after a clone is mutated")
	assert.Equal(t, float64(256), dim)
}

// TestVrsBuilder_Build_DoesNotMutateDefaultVrsCache pins the invariant that
// lets mergeOverlay pass the cached default VRS template directly into
// deepMergeMap without a per-row DeepCopy: neither Build nor a repeated Build
// on the same builder may mutate cfg.DefaultVrs.Us.
func TestVrsBuilder_Build_DoesNotMutateDefaultVrsCache(t *testing.T) {
	cfg := initConfig(t)
	cr := loadValdOperatorReleaseFromYAML(t)

	pristine := cfg.DefaultVrs.Us.DeepCopy()

	b := newVrsBuilder(cr, cfg, alwaysAvailable())
	first, err := b.Build(context.Background())
	assert.NoError(t, err)
	assert.NotEmpty(t, first)
	assert.Equal(t, pristine.Object, cfg.DefaultVrs.Us.Object,
		"the cached default VRS template must stay intact after Build")

	// A second Build through the same cached template must both leave the
	// cache intact and produce the same objects as the first one.
	second, err := newVrsBuilder(cr, cfg, alwaysAvailable()).Build(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, pristine.Object, cfg.DefaultVrs.Us.Object,
		"the cached default VRS template must stay intact after repeated Builds")
	assert.Equal(t, first, second, "repeated Builds must be deterministic")
}

// TestVrsBuilder_Build_WithinInfraOnlyNameVaries pins the invariant that lets
// mergeOverlay hoist the row's ToUnstructured conversion to once per
// infrastructure: within one infra, the built items differ only in
// metadata.name (every other field is derived from cr/cfg/infra, never from
// the individual cluster).
func TestVrsBuilder_Build_WithinInfraOnlyNameVaries(t *testing.T) {
	cfg := initConfig(t)
	cr := loadValdOperatorReleaseFromYAML(t)

	items, err := newVrsBuilder(cr, cfg, alwaysAvailable()).Build(context.Background())
	assert.NoError(t, err)
	if len(items) < 2 {
		t.Skipf("fixture yields %d items; need >=2 clusters in one infra to exercise the invariant", len(items))
	}

	names := make(map[string]bool, len(items))
	normalized := make([]map[string]any, 0, len(items))
	for _, it := range items {
		names[it.GetName()] = true
		us, uerr := resource.ToUnstructured(it)
		assert.NoError(t, uerr)
		us.SetName("normalized")
		normalized = append(normalized, us.Object)
	}
	assert.Len(t, names, len(items), "every built item must carry a distinct name (no aliasing across clusters)")
	for i := 1; i < len(normalized); i++ {
		if !reflect.DeepEqual(normalized[0], normalized[i]) {
			// Items from different infras may legitimately differ; the
			// invariant only holds within one infra. The fixture builds one
			// infra role per pair, so compare within pairs when totals allow.
			continue
		}
	}
	// Strict pairwise check within the first infra (fixture: clusters of the
	// same infra are adjacent in Build output).
	assert.True(t, reflect.DeepEqual(normalized[0], normalized[1]),
		"items of the same infra must be identical apart from metadata.name")
}

// BenchmarkVrsBuilder_Build tracks the Build hot path (called per reconcile);
// the overlay is parsed once per Build and cloned per cluster instead of
// re-unmarshalling the raw JSON on every iteration.
func BenchmarkVrsBuilder_Build(b *testing.B) {
	cfg := initConfig(b)
	cr := loadValdOperatorReleaseFromYAML(b)
	ctx := context.Background()

	b.ReportAllocs()
	for b.Loop() {
		if _, err := newVrsBuilder(cr, cfg, alwaysAvailable()).Build(ctx); err != nil {
			b.Fatal(err)
		}
	}
}

// NOT IMPLEMENTED BELOW
//
// func Test_newVrsBuilder(t *testing.T) {
// 	type args struct {
// 		cr         *v1.ValdOperatorRelease
// 		cfg        *config.Config
// 		capability nodePoolCapability
// 	}
// 	type want struct {
// 		want *vrsBuilder
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, *vrsBuilder) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got *vrsBuilder) error {
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
// 		       args: args {
// 		           cr:nil,
// 		           cfg:nil,
// 		           capability:nodePoolCapability{},
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
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
// 		           args: args {
// 		           cr:nil,
// 		           cfg:nil,
// 		           capability:nodePoolCapability{},
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
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
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got := newVrsBuilder(test.args.cr, test.args.cfg, test.args.capability)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_vrsBuilder_Build(t *testing.T) {
// 	type args struct {
// 		in0 context.Context
// 	}
// 	type fields struct {
// 		list       *valdrelease.ValdReleaseList
// 		cr         *v1.ValdOperatorRelease
// 		cfg        *config.Config
// 		overlay    *k8s.Unstructured
// 		capability nodePoolCapability
// 	}
// 	type want struct {
// 		want []k8s.Object
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, []k8s.Object, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got []k8s.Object, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
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
// 		       args: args {
// 		           in0:nil,
// 		       },
// 		       fields: fields {
// 		           list:nil,
// 		           cr:nil,
// 		           cfg:nil,
// 		           overlay:nil,
// 		           capability:nodePoolCapability{},
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
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
// 		           args: args {
// 		           in0:nil,
// 		           },
// 		           fields: fields {
// 		           list:nil,
// 		           cr:nil,
// 		           cfg:nil,
// 		           overlay:nil,
// 		           capability:nodePoolCapability{},
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
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
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			b := &vrsBuilder{
// 				list:       test.fields.list,
// 				cr:         test.fields.cr,
// 				cfg:        test.fields.cfg,
// 				overlay:    test.fields.overlay,
// 				capability: test.fields.capability,
// 			}
//
// 			got, err := b.Build(test.args.in0)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_fetchExistingVrs(t *testing.T) {
// 	type args struct {
// 		ctx       context.Context
// 		c         k8s.Client
// 		namespace string
// 	}
// 	type want struct {
// 		want []k8s.Object
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, []k8s.Object, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got []k8s.Object, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
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
// 		       args: args {
// 		           ctx:nil,
// 		           c:nil,
// 		           namespace:"",
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
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
// 		           args: args {
// 		           ctx:nil,
// 		           c:nil,
// 		           namespace:"",
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
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
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got, err := fetchExistingVrs(test.args.ctx, test.args.c, test.args.namespace)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_vrsBuilder_validate(t *testing.T) {
// 	type fields struct {
// 		list       *valdrelease.ValdReleaseList
// 		cr         *v1.ValdOperatorRelease
// 		cfg        *config.Config
// 		overlay    *k8s.Unstructured
// 		capability nodePoolCapability
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           list:nil,
// 		           cr:nil,
// 		           cfg:nil,
// 		           overlay:nil,
// 		           capability:nodePoolCapability{},
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
// 		           list:nil,
// 		           cr:nil,
// 		           cfg:nil,
// 		           overlay:nil,
// 		           capability:nodePoolCapability{},
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
// 			b := &vrsBuilder{
// 				list:       test.fields.list,
// 				cr:         test.fields.cr,
// 				cfg:        test.fields.cfg,
// 				overlay:    test.fields.overlay,
// 				capability: test.fields.capability,
// 			}
//
// 			err := b.validate()
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_vrsBuilder_buildDefaults(t *testing.T) {
// 	type fields struct {
// 		list       *valdrelease.ValdReleaseList
// 		cr         *v1.ValdOperatorRelease
// 		cfg        *config.Config
// 		overlay    *k8s.Unstructured
// 		capability nodePoolCapability
// 	}
// 	type want struct {
// 		want *valdrelease.Defaults
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *valdrelease.Defaults) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got *valdrelease.Defaults) error {
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
// 		           list:nil,
// 		           cr:nil,
// 		           cfg:nil,
// 		           overlay:nil,
// 		           capability:nodePoolCapability{},
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
// 		           list:nil,
// 		           cr:nil,
// 		           cfg:nil,
// 		           overlay:nil,
// 		           capability:nodePoolCapability{},
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
// 			b := &vrsBuilder{
// 				list:       test.fields.list,
// 				cr:         test.fields.cr,
// 				cfg:        test.fields.cfg,
// 				overlay:    test.fields.overlay,
// 				capability: test.fields.capability,
// 			}
//
// 			got := b.buildDefaults()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_vrsBuilder_buildAgent(t *testing.T) {
// 	type fields struct {
// 		list       *valdrelease.ValdReleaseList
// 		cr         *v1.ValdOperatorRelease
// 		cfg        *config.Config
// 		overlay    *k8s.Unstructured
// 		capability nodePoolCapability
// 	}
// 	type want struct {
// 		want *valdrelease.Agent
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *valdrelease.Agent) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got *valdrelease.Agent) error {
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
// 		           list:nil,
// 		           cr:nil,
// 		           cfg:nil,
// 		           overlay:nil,
// 		           capability:nodePoolCapability{},
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
// 		           list:nil,
// 		           cr:nil,
// 		           cfg:nil,
// 		           overlay:nil,
// 		           capability:nodePoolCapability{},
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
// 			b := &vrsBuilder{
// 				list:       test.fields.list,
// 				cr:         test.fields.cr,
// 				cfg:        test.fields.cfg,
// 				overlay:    test.fields.overlay,
// 				capability: test.fields.capability,
// 			}
//
// 			got := b.buildAgent()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_vrsBuilder_buildAgentNgt(t *testing.T) {
// 	type fields struct {
// 		list       *valdrelease.ValdReleaseList
// 		cr         *v1.ValdOperatorRelease
// 		cfg        *config.Config
// 		overlay    *k8s.Unstructured
// 		capability nodePoolCapability
// 	}
// 	type want struct {
// 		want *valdrelease.AgentNgt
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *valdrelease.AgentNgt) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got *valdrelease.AgentNgt) error {
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
// 		           list:nil,
// 		           cr:nil,
// 		           cfg:nil,
// 		           overlay:nil,
// 		           capability:nodePoolCapability{},
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
// 		           list:nil,
// 		           cr:nil,
// 		           cfg:nil,
// 		           overlay:nil,
// 		           capability:nodePoolCapability{},
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
// 			b := &vrsBuilder{
// 				list:       test.fields.list,
// 				cr:         test.fields.cr,
// 				cfg:        test.fields.cfg,
// 				overlay:    test.fields.overlay,
// 				capability: test.fields.capability,
// 			}
//
// 			got := b.buildAgentNgt()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_vrsBuilder_buildGateway(t *testing.T) {
// 	type fields struct {
// 		list       *valdrelease.ValdReleaseList
// 		cr         *v1.ValdOperatorRelease
// 		cfg        *config.Config
// 		overlay    *k8s.Unstructured
// 		capability nodePoolCapability
// 	}
// 	type want struct {
// 		want *valdrelease.Gateway
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *valdrelease.Gateway) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got *valdrelease.Gateway) error {
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
// 		           list:nil,
// 		           cr:nil,
// 		           cfg:nil,
// 		           overlay:nil,
// 		           capability:nodePoolCapability{},
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
// 		           list:nil,
// 		           cr:nil,
// 		           cfg:nil,
// 		           overlay:nil,
// 		           capability:nodePoolCapability{},
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
// 			b := &vrsBuilder{
// 				list:       test.fields.list,
// 				cr:         test.fields.cr,
// 				cfg:        test.fields.cfg,
// 				overlay:    test.fields.overlay,
// 				capability: test.fields.capability,
// 			}
//
// 			got := b.buildGateway()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_vrsBuilder_buildLb(t *testing.T) {
// 	type fields struct {
// 		list       *valdrelease.ValdReleaseList
// 		cr         *v1.ValdOperatorRelease
// 		cfg        *config.Config
// 		overlay    *k8s.Unstructured
// 		capability nodePoolCapability
// 	}
// 	type want struct {
// 		want *valdrelease.GatewayLb
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *valdrelease.GatewayLb) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got *valdrelease.GatewayLb) error {
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
// 		           list:nil,
// 		           cr:nil,
// 		           cfg:nil,
// 		           overlay:nil,
// 		           capability:nodePoolCapability{},
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
// 		           list:nil,
// 		           cr:nil,
// 		           cfg:nil,
// 		           overlay:nil,
// 		           capability:nodePoolCapability{},
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
// 			b := &vrsBuilder{
// 				list:       test.fields.list,
// 				cr:         test.fields.cr,
// 				cfg:        test.fields.cfg,
// 				overlay:    test.fields.overlay,
// 				capability: test.fields.capability,
// 			}
//
// 			got := b.buildLb()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_vrsBuilder_buildIngress(t *testing.T) {
// 	type args struct {
// 		in *v1.GatewayIngress
// 	}
// 	type fields struct {
// 		list       *valdrelease.ValdReleaseList
// 		cr         *v1.ValdOperatorRelease
// 		cfg        *config.Config
// 		overlay    *k8s.Unstructured
// 		capability nodePoolCapability
// 	}
// 	type want struct {
// 		want *valdrelease.GatewayLbIngress
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *valdrelease.GatewayLbIngress) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got *valdrelease.GatewayLbIngress) error {
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
// 		       args: args {
// 		           in:nil,
// 		       },
// 		       fields: fields {
// 		           list:nil,
// 		           cr:nil,
// 		           cfg:nil,
// 		           overlay:nil,
// 		           capability:nodePoolCapability{},
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
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
// 		           args: args {
// 		           in:nil,
// 		           },
// 		           fields: fields {
// 		           list:nil,
// 		           cr:nil,
// 		           cfg:nil,
// 		           overlay:nil,
// 		           capability:nodePoolCapability{},
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
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
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			b := &vrsBuilder{
// 				list:       test.fields.list,
// 				cr:         test.fields.cr,
// 				cfg:        test.fields.cfg,
// 				overlay:    test.fields.overlay,
// 				capability: test.fields.capability,
// 			}
//
// 			got := b.buildIngress(test.args.in)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_vrsBuilder_getGatewayIngressPathType(t *testing.T) {
// 	type fields struct {
// 		list       *valdrelease.ValdReleaseList
// 		cr         *v1.ValdOperatorRelease
// 		cfg        *config.Config
// 		overlay    *k8s.Unstructured
// 		capability nodePoolCapability
// 	}
// 	type want struct {
// 		want string
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, string) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got string) error {
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
// 		           list:nil,
// 		           cr:nil,
// 		           cfg:nil,
// 		           overlay:nil,
// 		           capability:nodePoolCapability{},
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
// 		           list:nil,
// 		           cr:nil,
// 		           cfg:nil,
// 		           overlay:nil,
// 		           capability:nodePoolCapability{},
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
// 			b := &vrsBuilder{
// 				list:       test.fields.list,
// 				cr:         test.fields.cr,
// 				cfg:        test.fields.cfg,
// 				overlay:    test.fields.overlay,
// 				capability: test.fields.capability,
// 			}
//
// 			got := b.getGatewayIngressPathType()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_vrsBuilder_getGatewayServiceType(t *testing.T) {
// 	type args struct {
// 		st string
// 	}
// 	type fields struct {
// 		list       *valdrelease.ValdReleaseList
// 		cr         *v1.ValdOperatorRelease
// 		cfg        *config.Config
// 		overlay    *k8s.Unstructured
// 		capability nodePoolCapability
// 	}
// 	type want struct {
// 		want *valdrelease.GatewayLbServiceType
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *valdrelease.GatewayLbServiceType) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got *valdrelease.GatewayLbServiceType) error {
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
// 		       args: args {
// 		           st:"",
// 		       },
// 		       fields: fields {
// 		           list:nil,
// 		           cr:nil,
// 		           cfg:nil,
// 		           overlay:nil,
// 		           capability:nodePoolCapability{},
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
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
// 		           args: args {
// 		           st:"",
// 		           },
// 		           fields: fields {
// 		           list:nil,
// 		           cr:nil,
// 		           cfg:nil,
// 		           overlay:nil,
// 		           capability:nodePoolCapability{},
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
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
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			b := &vrsBuilder{
// 				list:       test.fields.list,
// 				cr:         test.fields.cr,
// 				cfg:        test.fields.cfg,
// 				overlay:    test.fields.overlay,
// 				capability: test.fields.capability,
// 			}
//
// 			got := b.getGatewayServiceType(test.args.st)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_vrsBuilder_buildDiscoverer(t *testing.T) {
// 	type fields struct {
// 		list       *valdrelease.ValdReleaseList
// 		cr         *v1.ValdOperatorRelease
// 		cfg        *config.Config
// 		overlay    *k8s.Unstructured
// 		capability nodePoolCapability
// 	}
// 	type want struct {
// 		want *valdrelease.Discoverer
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *valdrelease.Discoverer) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got *valdrelease.Discoverer) error {
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
// 		           list:nil,
// 		           cr:nil,
// 		           cfg:nil,
// 		           overlay:nil,
// 		           capability:nodePoolCapability{},
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
// 		           list:nil,
// 		           cr:nil,
// 		           cfg:nil,
// 		           overlay:nil,
// 		           capability:nodePoolCapability{},
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
// 			b := &vrsBuilder{
// 				list:       test.fields.list,
// 				cr:         test.fields.cr,
// 				cfg:        test.fields.cfg,
// 				overlay:    test.fields.overlay,
// 				capability: test.fields.capability,
// 			}
//
// 			got := b.buildDiscoverer()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_vrsBuilder_buildManager(t *testing.T) {
// 	type fields struct {
// 		list       *valdrelease.ValdReleaseList
// 		cr         *v1.ValdOperatorRelease
// 		cfg        *config.Config
// 		overlay    *k8s.Unstructured
// 		capability nodePoolCapability
// 	}
// 	type want struct {
// 		want *valdrelease.Manager
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *valdrelease.Manager) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got *valdrelease.Manager) error {
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
// 		           list:nil,
// 		           cr:nil,
// 		           cfg:nil,
// 		           overlay:nil,
// 		           capability:nodePoolCapability{},
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
// 		           list:nil,
// 		           cr:nil,
// 		           cfg:nil,
// 		           overlay:nil,
// 		           capability:nodePoolCapability{},
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
// 			b := &vrsBuilder{
// 				list:       test.fields.list,
// 				cr:         test.fields.cr,
// 				cfg:        test.fields.cfg,
// 				overlay:    test.fields.overlay,
// 				capability: test.fields.capability,
// 			}
//
// 			got := b.buildManager()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_vrsBuilder_makeName(t *testing.T) {
// 	type args struct {
// 		str []string
// 	}
// 	type fields struct {
// 		list       *valdrelease.ValdReleaseList
// 		cr         *v1.ValdOperatorRelease
// 		cfg        *config.Config
// 		overlay    *k8s.Unstructured
// 		capability nodePoolCapability
// 	}
// 	type want struct {
// 		want string
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, string, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got string, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
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
// 		       args: args {
// 		           str:nil,
// 		       },
// 		       fields: fields {
// 		           list:nil,
// 		           cr:nil,
// 		           cfg:nil,
// 		           overlay:nil,
// 		           capability:nodePoolCapability{},
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
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
// 		           args: args {
// 		           str:nil,
// 		           },
// 		           fields: fields {
// 		           list:nil,
// 		           cr:nil,
// 		           cfg:nil,
// 		           overlay:nil,
// 		           capability:nodePoolCapability{},
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
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
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			b := &vrsBuilder{
// 				list:       test.fields.list,
// 				cr:         test.fields.cr,
// 				cfg:        test.fields.cfg,
// 				overlay:    test.fields.overlay,
// 				capability: test.fields.capability,
// 			}
//
// 			got, err := b.makeName(test.args.str...)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_vrsBuilder_buildLabels(t *testing.T) {
// 	type args struct {
// 		iKey int
// 	}
// 	type fields struct {
// 		list       *valdrelease.ValdReleaseList
// 		cr         *v1.ValdOperatorRelease
// 		cfg        *config.Config
// 		overlay    *k8s.Unstructured
// 		capability nodePoolCapability
// 	}
// 	type want struct {
// 		want map[string]string
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, map[string]string) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got map[string]string) error {
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
// 		       args: args {
// 		           iKey:0,
// 		       },
// 		       fields: fields {
// 		           list:nil,
// 		           cr:nil,
// 		           cfg:nil,
// 		           overlay:nil,
// 		           capability:nodePoolCapability{},
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
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
// 		           args: args {
// 		           iKey:0,
// 		           },
// 		           fields: fields {
// 		           list:nil,
// 		           cr:nil,
// 		           cfg:nil,
// 		           overlay:nil,
// 		           capability:nodePoolCapability{},
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
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
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			b := &vrsBuilder{
// 				list:       test.fields.list,
// 				cr:         test.fields.cr,
// 				cfg:        test.fields.cfg,
// 				overlay:    test.fields.overlay,
// 				capability: test.fields.capability,
// 			}
//
// 			got := b.buildLabels(test.args.iKey)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_mergeLabels(t *testing.T) {
// 	type args struct {
// 		base    map[string]string
// 		overlay map[string]string
// 	}
// 	type want struct {
// 		want map[string]string
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, map[string]string) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got map[string]string) error {
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
// 		       args: args {
// 		           base:nil,
// 		           overlay:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
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
// 		           args: args {
// 		           base:nil,
// 		           overlay:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
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
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got := mergeLabels(test.args.base, test.args.overlay)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_vrsBuilder_buildLogging(t *testing.T) {
// 	type args struct {
// 		ll string
// 	}
// 	type fields struct {
// 		list       *valdrelease.ValdReleaseList
// 		cr         *v1.ValdOperatorRelease
// 		cfg        *config.Config
// 		overlay    *k8s.Unstructured
// 		capability nodePoolCapability
// 	}
// 	type want struct {
// 		want *valdrelease.Logging
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *valdrelease.Logging) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got *valdrelease.Logging) error {
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
// 		       args: args {
// 		           ll:"",
// 		       },
// 		       fields: fields {
// 		           list:nil,
// 		           cr:nil,
// 		           cfg:nil,
// 		           overlay:nil,
// 		           capability:nodePoolCapability{},
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
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
// 		           args: args {
// 		           ll:"",
// 		           },
// 		           fields: fields {
// 		           list:nil,
// 		           cr:nil,
// 		           cfg:nil,
// 		           overlay:nil,
// 		           capability:nodePoolCapability{},
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
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
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			b := &vrsBuilder{
// 				list:       test.fields.list,
// 				cr:         test.fields.cr,
// 				cfg:        test.fields.cfg,
// 				overlay:    test.fields.overlay,
// 				capability: test.fields.capability,
// 			}
//
// 			got := b.buildLogging(test.args.ll)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_ptr(t *testing.T) {
// 	type args struct {
// 		v T
// 	}
// 	type want struct {
// 		want *T
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, *T) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got *T) error {
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
// 		       args: args {
// 		           v:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
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
// 		           args: args {
// 		           v:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
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
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got := ptr(test.args.v)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_toAnyMap(t *testing.T) {
// 	type args struct {
// 		m map[string]string
// 	}
// 	type want struct {
// 		want map[string]any
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, map[string]any) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got map[string]any) error {
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
// 		       args: args {
// 		           m:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
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
// 		           args: args {
// 		           m:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
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
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got := toAnyMap(test.args.m)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_nsPtr(t *testing.T) {
// 	type args struct {
// 		m map[string]string
// 	}
// 	type want struct {
// 		want *valdrelease.NodeSelector
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, *valdrelease.NodeSelector) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got *valdrelease.NodeSelector) error {
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
// 		       args: args {
// 		           m:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
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
// 		           args: args {
// 		           m:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
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
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got := nsPtr(test.args.m)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_vrsBuilder_labelKey(t *testing.T) {
// 	type args struct {
// 		suffix string
// 	}
// 	type fields struct {
// 		list       *valdrelease.ValdReleaseList
// 		cr         *v1.ValdOperatorRelease
// 		cfg        *config.Config
// 		overlay    *k8s.Unstructured
// 		capability nodePoolCapability
// 	}
// 	type want struct {
// 		want string
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, string) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got string) error {
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
// 		       args: args {
// 		           suffix:"",
// 		       },
// 		       fields: fields {
// 		           list:nil,
// 		           cr:nil,
// 		           cfg:nil,
// 		           overlay:nil,
// 		           capability:nodePoolCapability{},
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
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
// 		           args: args {
// 		           suffix:"",
// 		           },
// 		           fields: fields {
// 		           list:nil,
// 		           cr:nil,
// 		           cfg:nil,
// 		           overlay:nil,
// 		           capability:nodePoolCapability{},
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
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
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			b := &vrsBuilder{
// 				list:       test.fields.list,
// 				cr:         test.fields.cr,
// 				cfg:        test.fields.cfg,
// 				overlay:    test.fields.overlay,
// 				capability: test.fields.capability,
// 			}
//
// 			got := b.labelKey(test.args.suffix)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_vrsBuilder_mergeOverlay(t *testing.T) {
// 	type args struct {
// 		current *k8s.Unstructured
// 	}
// 	type fields struct {
// 		list       *valdrelease.ValdReleaseList
// 		cr         *v1.ValdOperatorRelease
// 		cfg        *config.Config
// 		overlay    *k8s.Unstructured
// 		capability nodePoolCapability
// 	}
// 	type want struct {
// 		want *valdrelease.ValdRelease
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *valdrelease.ValdRelease, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got *valdrelease.ValdRelease, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
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
// 		       args: args {
// 		           current:nil,
// 		       },
// 		       fields: fields {
// 		           list:nil,
// 		           cr:nil,
// 		           cfg:nil,
// 		           overlay:nil,
// 		           capability:nodePoolCapability{},
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
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
// 		           args: args {
// 		           current:nil,
// 		           },
// 		           fields: fields {
// 		           list:nil,
// 		           cr:nil,
// 		           cfg:nil,
// 		           overlay:nil,
// 		           capability:nodePoolCapability{},
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
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
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			b := &vrsBuilder{
// 				list:       test.fields.list,
// 				cr:         test.fields.cr,
// 				cfg:        test.fields.cfg,
// 				overlay:    test.fields.overlay,
// 				capability: test.fields.capability,
// 			}
//
// 			got, err := b.mergeOverlay(test.args.current)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_deepMergeMap(t *testing.T) {
// 	type args struct {
// 		dst map[string]any
// 		src map[string]any
// 	}
// 	type want struct {
// 		want map[string]any
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, map[string]any) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got map[string]any) error {
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
// 		       args: args {
// 		           dst:nil,
// 		           src:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
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
// 		           args: args {
// 		           dst:nil,
// 		           src:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
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
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got := deepMergeMap(test.args.dst, test.args.src)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_vrsBuilder_parseOverlay(t *testing.T) {
// 	type fields struct {
// 		list       *valdrelease.ValdReleaseList
// 		cr         *v1.ValdOperatorRelease
// 		cfg        *config.Config
// 		overlay    *k8s.Unstructured
// 		capability nodePoolCapability
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           list:nil,
// 		           cr:nil,
// 		           cfg:nil,
// 		           overlay:nil,
// 		           capability:nodePoolCapability{},
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
// 		           list:nil,
// 		           cr:nil,
// 		           cfg:nil,
// 		           overlay:nil,
// 		           capability:nodePoolCapability{},
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
// 			b := &vrsBuilder{
// 				list:       test.fields.list,
// 				cr:         test.fields.cr,
// 				cfg:        test.fields.cfg,
// 				overlay:    test.fields.overlay,
// 				capability: test.fields.capability,
// 			}
//
// 			err := b.parseOverlay()
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_vrsBuilder_makeOverlayPatch(t *testing.T) {
// 	type args struct {
// 		row k8s.Object
// 	}
// 	type fields struct {
// 		list       *valdrelease.ValdReleaseList
// 		cr         *v1.ValdOperatorRelease
// 		cfg        *config.Config
// 		overlay    *k8s.Unstructured
// 		capability nodePoolCapability
// 	}
// 	type want struct {
// 		want *k8s.Unstructured
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *k8s.Unstructured) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got *k8s.Unstructured) error {
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
// 		       args: args {
// 		           row:nil,
// 		       },
// 		       fields: fields {
// 		           list:nil,
// 		           cr:nil,
// 		           cfg:nil,
// 		           overlay:nil,
// 		           capability:nodePoolCapability{},
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
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
// 		           args: args {
// 		           row:nil,
// 		           },
// 		           fields: fields {
// 		           list:nil,
// 		           cr:nil,
// 		           cfg:nil,
// 		           overlay:nil,
// 		           capability:nodePoolCapability{},
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
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
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			b := &vrsBuilder{
// 				list:       test.fields.list,
// 				cr:         test.fields.cr,
// 				cfg:        test.fields.cfg,
// 				overlay:    test.fields.overlay,
// 				capability: test.fields.capability,
// 			}
//
// 			got := b.makeOverlayPatch(test.args.row)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_vrsBuilder_buildNodeSelector(t *testing.T) {
// 	type args struct {
// 		nt v1.NodePoolType
// 	}
// 	type fields struct {
// 		list       *valdrelease.ValdReleaseList
// 		cr         *v1.ValdOperatorRelease
// 		cfg        *config.Config
// 		overlay    *k8s.Unstructured
// 		capability nodePoolCapability
// 	}
// 	type want struct {
// 		want map[string]string
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, map[string]string) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got map[string]string) error {
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
// 		       args: args {
// 		           nt:nil,
// 		       },
// 		       fields: fields {
// 		           list:nil,
// 		           cr:nil,
// 		           cfg:nil,
// 		           overlay:nil,
// 		           capability:nodePoolCapability{},
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
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
// 		           args: args {
// 		           nt:nil,
// 		           },
// 		           fields: fields {
// 		           list:nil,
// 		           cr:nil,
// 		           cfg:nil,
// 		           overlay:nil,
// 		           capability:nodePoolCapability{},
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
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
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			b := &vrsBuilder{
// 				list:       test.fields.list,
// 				cr:         test.fields.cr,
// 				cfg:        test.fields.cfg,
// 				overlay:    test.fields.overlay,
// 				capability: test.fields.capability,
// 			}
//
// 			got := b.buildNodeSelector(test.args.nt)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_vrsBuilder_buildToleration(t *testing.T) {
// 	type args struct {
// 		nt v1.NodePoolType
// 	}
// 	type fields struct {
// 		list       *valdrelease.ValdReleaseList
// 		cr         *v1.ValdOperatorRelease
// 		cfg        *config.Config
// 		overlay    *k8s.Unstructured
// 		capability nodePoolCapability
// 	}
// 	type want struct {
// 		want *[]k8s.Toleration
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *[]k8s.Toleration) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got *[]k8s.Toleration) error {
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
// 		       args: args {
// 		           nt:nil,
// 		       },
// 		       fields: fields {
// 		           list:nil,
// 		           cr:nil,
// 		           cfg:nil,
// 		           overlay:nil,
// 		           capability:nodePoolCapability{},
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
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
// 		           args: args {
// 		           nt:nil,
// 		           },
// 		           fields: fields {
// 		           list:nil,
// 		           cr:nil,
// 		           cfg:nil,
// 		           overlay:nil,
// 		           capability:nodePoolCapability{},
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
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
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			b := &vrsBuilder{
// 				list:       test.fields.list,
// 				cr:         test.fields.cr,
// 				cfg:        test.fields.cfg,
// 				overlay:    test.fields.overlay,
// 				capability: test.fields.capability,
// 			}
//
// 			got := b.buildToleration(test.args.nt)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_vrsBuilder_effectiveNodePoolType(t *testing.T) {
// 	type args struct {
// 		nt v1.NodePoolType
// 	}
// 	type fields struct {
// 		list       *valdrelease.ValdReleaseList
// 		cr         *v1.ValdOperatorRelease
// 		cfg        *config.Config
// 		overlay    *k8s.Unstructured
// 		capability nodePoolCapability
// 	}
// 	type want struct {
// 		want v1.NodePoolType
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, v1.NodePoolType) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got v1.NodePoolType) error {
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
// 		       args: args {
// 		           nt:nil,
// 		       },
// 		       fields: fields {
// 		           list:nil,
// 		           cr:nil,
// 		           cfg:nil,
// 		           overlay:nil,
// 		           capability:nodePoolCapability{},
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
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
// 		           args: args {
// 		           nt:nil,
// 		           },
// 		           fields: fields {
// 		           list:nil,
// 		           cr:nil,
// 		           cfg:nil,
// 		           overlay:nil,
// 		           capability:nodePoolCapability{},
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
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
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			b := &vrsBuilder{
// 				list:       test.fields.list,
// 				cr:         test.fields.cr,
// 				cfg:        test.fields.cfg,
// 				overlay:    test.fields.overlay,
// 				capability: test.fields.capability,
// 			}
//
// 			got := b.effectiveNodePoolType(test.args.nt)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_vrsBuilder_applyNodeAffinities(t *testing.T) {
// 	type args struct {
// 		v *valdrelease.ValdRelease
// 	}
// 	type fields struct {
// 		list       *valdrelease.ValdReleaseList
// 		cr         *v1.ValdOperatorRelease
// 		cfg        *config.Config
// 		overlay    *k8s.Unstructured
// 		capability nodePoolCapability
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           v:nil,
// 		       },
// 		       fields: fields {
// 		           list:nil,
// 		           cr:nil,
// 		           cfg:nil,
// 		           overlay:nil,
// 		           capability:nodePoolCapability{},
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
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
// 		           args: args {
// 		           v:nil,
// 		           },
// 		           fields: fields {
// 		           list:nil,
// 		           cr:nil,
// 		           cfg:nil,
// 		           overlay:nil,
// 		           capability:nodePoolCapability{},
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
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
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			b := &vrsBuilder{
// 				list:       test.fields.list,
// 				cr:         test.fields.cr,
// 				cfg:        test.fields.cfg,
// 				overlay:    test.fields.overlay,
// 				capability: test.fields.capability,
// 			}
//
// 			b.applyNodeAffinities(test.args.v)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_vrsBuilder_reflectPersistentVolume(t *testing.T) {
// 	type args struct {
// 		v *valdrelease.ValdRelease
// 	}
// 	type fields struct {
// 		list       *valdrelease.ValdReleaseList
// 		cr         *v1.ValdOperatorRelease
// 		cfg        *config.Config
// 		overlay    *k8s.Unstructured
// 		capability nodePoolCapability
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           v:nil,
// 		       },
// 		       fields: fields {
// 		           list:nil,
// 		           cr:nil,
// 		           cfg:nil,
// 		           overlay:nil,
// 		           capability:nodePoolCapability{},
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
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
// 		           args: args {
// 		           v:nil,
// 		           },
// 		           fields: fields {
// 		           list:nil,
// 		           cr:nil,
// 		           cfg:nil,
// 		           overlay:nil,
// 		           capability:nodePoolCapability{},
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
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
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			b := &vrsBuilder{
// 				list:       test.fields.list,
// 				cr:         test.fields.cr,
// 				cfg:        test.fields.cfg,
// 				overlay:    test.fields.overlay,
// 				capability: test.fields.capability,
// 			}
//
// 			b.reflectPersistentVolume(test.args.v)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }

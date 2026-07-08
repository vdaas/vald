//go:build e2e

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

// package crud provides end-to-end tests using ann-benchmarks datasets.
package crud

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/k8s/resource"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/tests/v2/e2e/config"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	yaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/util/jsonpath"
)

const (
	// fieldManager identifies this test harness in server-side apply operations.
	fieldManager = "vald-e2e"
	// customResourcePollInterval is the polling interval for custom resource wait actions.
	customResourcePollInterval = 5 * time.Second
	// manifestDecodeBufferSize is the initial buffer size for the multi-document yaml decoder.
	manifestDecodeBufferSize = 4096
)

func (r *runner) processKubernetes(t *testing.T, ctx context.Context, plan *config.Execution) {
	t.Helper()
	if plan == nil || plan.Kubernetes == nil {
		t.Fatal("kubernetes plan is nil")
		return
	}
	k := plan.Kubernetes
	// Manifest-based apply/delete operates on arbitrary resources through the
	// dynamic client, so it is dispatched before the typed kind switch.
	if k.Manifest != "" {
		r.handleKubernetesManifest(t, ctx, k)
		return
	}
	if k.Kind == config.CustomResource {
		r.handleKubernetesCustomResource(t, ctx, k)
		return
	}
	// The create action is special-cased because it always creates a job from
	// the cronjob named k.Name regardless of the requested kind.
	if k.Action == config.KubernetesActionCreate {
		if k.Kind != config.Job {
			t.Errorf("kubernetes action create is only supported for creating job from cronjob")
			return
		}
		if _, err := resource.CronJob(r.k8s, k.Namespace).CreateJob(ctx, k.Name, resource.EmptyGetOptions, resource.EmptyCreateOptions); err != nil {
			t.Errorf("failed to create job from cronjob: %v", err)
		}
		return
	}
	switch k.Kind {
	case config.ConfigMap:
		handleKubernetesAction(t, ctx, k, resource.ConfigMap(r.k8s, k.Namespace))
	case config.CronJob:
		handleKubernetesWorkloadAction(t, ctx, k, resource.CronJob(r.k8s, k.Namespace))
	case config.DaemonSet:
		handleKubernetesWorkloadAction(t, ctx, k, resource.DaemonSet(r.k8s, k.Namespace))
	case config.Deployment:
		handleKubernetesWorkloadAction(t, ctx, k, resource.Deployment(r.k8s, k.Namespace))
	case config.Job:
		handleKubernetesWorkloadAction(t, ctx, k, resource.Job(r.k8s, k.Namespace))
	case config.MutatingWebhookConfiguration:
		handleKubernetesAction(t, ctx, k, resource.MutatingWebhookConfiguration(r.k8s))
	case config.Pod:
		handleKubernetesAction(t, ctx, k, resource.Pod(r.k8s, k.Namespace))
	case config.Secret:
		handleKubernetesAction(t, ctx, k, resource.Secret(r.k8s, k.Namespace))
	case config.Service:
		handleKubernetesAction(t, ctx, k, resource.Service(r.k8s, k.Namespace))
	case config.StatefulSet:
		handleKubernetesWorkloadAction(t, ctx, k, resource.StatefulSet(r.k8s, k.Namespace))
	case config.ValidatingWebhookConfiguration:
		handleKubernetesAction(t, ctx, k, resource.ValidatingWebhookConfiguration(r.k8s))
	default:
		t.Errorf("unsupported kubernetes kind: %s", k.Kind)
	}
}

// handleKubernetesAction executes the actions applicable to every resource kind
// (get, delete and wait) through the generic resource client.
func handleKubernetesAction[T resource.Object, L resource.ObjectList, C resource.NamedObject, I resource.ResourceInterface[T, L, C]](
	t *testing.T, ctx context.Context, k *config.KubernetesConfig, client I,
) {
	t.Helper()
	switch k.Action {
	case config.KubernetesActionGet:
		obj, err := client.Get(ctx, k.Name, resource.EmptyGetOptions)
		if err != nil {
			t.Errorf("failed to get %s: %v", k.Kind, err)
			return
		}
		log.Infof("kubernetes object: %v", obj)
	case config.KubernetesActionDelete:
		if err := client.Delete(ctx, k.Name, resource.EmptyDeleteOptions); err != nil {
			t.Errorf("failed to delete %s: %v", k.Kind, err)
		}
	case config.KubernetesActionWait:
		ok, err := resource.WaitForStatus(ctx, client, k.Name, k.LabelSelector, k.Status.Status())
		if !ok {
			t.Errorf("failed to wait for %s: %v", k.Kind, err)
		}
	default:
		t.Errorf("kubernetes action %s is not supported for kind %s", k.Action, k.Kind)
	}
}

// handleKubernetesWorkloadAction adds the workload-controller specific rollout
// action on top of the common resource actions.
func handleKubernetesWorkloadAction[T resource.Object, L resource.ObjectList, C resource.NamedObject, I resource.WorkloadControllerResourceClient[T, L, C]](
	t *testing.T, ctx context.Context, k *config.KubernetesConfig, client I,
) {
	t.Helper()
	switch k.Action {
	case config.KubernetesActionRollout:
		if err := resource.RolloutRestart(ctx, client, k.Name); err != nil {
			t.Errorf("failed to rollout restart %s: %v", k.Kind, err)
		}
	default:
		handleKubernetesAction(t, ctx, k, client)
	}
}

// newDynamicClient builds a dynamic client and a RESTMapper from the runner's
// kubernetes REST config so arbitrary GVKs can be resolved to GVRs.
func (r *runner) newDynamicClient() (dynamic.Interface, meta.RESTMapper, error) {
	if r.k8s == nil {
		return nil, nil, errors.New("kubernetes client is not initialized")
	}
	cfg := r.k8s.GetRESTConfig()
	if cfg == nil {
		return nil, nil, errors.New("kubernetes rest config is not initialized")
	}
	dyn, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create dynamic client")
	}
	dc, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create discovery client")
	}
	groups, err := restmapper.GetAPIGroupResources(dc)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to fetch api group resources")
	}
	return dyn, restmapper.NewDiscoveryRESTMapper(groups), nil
}

// decodeManifest reads a multi-document yaml manifest file and decodes each
// document into an unstructured object. Empty documents are skipped.
func decodeManifest(path string) (objs []*unstructured.Unstructured, err error) {
	data, err := file.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read manifest file %s", path)
	}
	dec := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(data), manifestDecodeBufferSize)
	for {
		obj := new(unstructured.Unstructured)
		if err := dec.Decode(obj); err != nil {
			if errors.Is(err, io.EOF) {
				return objs, nil
			}
			return nil, errors.Wrapf(err, "failed to decode manifest %s", path)
		}
		if len(obj.Object) == 0 {
			continue
		}
		objs = append(objs, obj)
	}
}

// resourceFor resolves the dynamic resource interface for the given object via
// the RESTMapper, honoring namespace scoping. For namespaced resources the
// namespace is taken from the object, then the config, then "default".
func resourceFor(
	dyn dynamic.Interface,
	mapper meta.RESTMapper,
	obj *unstructured.Unstructured,
	fallbackNamespace string,
) (dynamic.ResourceInterface, error) {
	gvk := obj.GroupVersionKind()
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to resolve rest mapping for %s", gvk.String())
	}
	if mapping.Scope.Name() != meta.RESTScopeNameNamespace {
		return dyn.Resource(mapping.Resource), nil
	}
	ns := obj.GetNamespace()
	if ns == "" {
		ns = fallbackNamespace
	}
	if ns == "" {
		ns = "default"
	}
	return dyn.Resource(mapping.Resource).Namespace(ns), nil
}

// handleKubernetesManifest applies or deletes every document in the configured
// manifest through the dynamic client. Apply uses server-side apply so the
// operation is idempotent across repeated e2e runs.
func (r *runner) handleKubernetesManifest(
	t *testing.T, ctx context.Context, k *config.KubernetesConfig,
) {
	t.Helper()
	objs, err := decodeManifest(k.Manifest)
	if err != nil {
		t.Errorf("failed to decode manifest %s: %v", k.Manifest, err)
		return
	}
	if len(objs) == 0 {
		t.Errorf("manifest %s contains no objects", k.Manifest)
		return
	}
	dyn, mapper, err := r.newDynamicClient()
	if err != nil {
		t.Errorf("failed to create dynamic client: %v", err)
		return
	}
	for _, obj := range objs {
		client, err := resourceFor(dyn, mapper, obj, k.Namespace)
		if err != nil {
			t.Errorf("failed to resolve resource for %s %s: %v", obj.GetKind(), obj.GetName(), err)
			return
		}
		switch k.Action {
		case config.KubernetesActionApply:
			data, err := obj.MarshalJSON()
			if err != nil {
				t.Errorf("failed to marshal %s %s: %v", obj.GetKind(), obj.GetName(), err)
				return
			}
			force := true
			if _, err := client.Patch(ctx, obj.GetName(), types.ApplyPatchType, data, metav1.PatchOptions{
				FieldManager: fieldManager,
				Force:        &force,
			}); err != nil {
				t.Errorf("failed to apply %s %s: %v", obj.GetKind(), obj.GetName(), err)
				return
			}
			log.Infof("applied %s %s from manifest %s", obj.GetKind(), obj.GetName(), k.Manifest)
		case config.KubernetesActionDelete:
			if err := client.Delete(ctx, obj.GetName(), resource.EmptyDeleteOptions); err != nil {
				if apierrors.IsNotFound(err) {
					log.Warnf("%s %s from manifest %s is already deleted", obj.GetKind(), obj.GetName(), k.Manifest)
					continue
				}
				t.Errorf("failed to delete %s %s: %v", obj.GetKind(), obj.GetName(), err)
				return
			}
			log.Infof("deleted %s %s from manifest %s", obj.GetKind(), obj.GetName(), k.Manifest)
		default:
			t.Errorf("kubernetes action %s is not supported with manifest", k.Action)
			return
		}
	}
}

// handleKubernetesCustomResource executes get/delete/wait actions on an
// arbitrary custom resource identified by group/version/resource.
func (r *runner) handleKubernetesCustomResource(
	t *testing.T, ctx context.Context, k *config.KubernetesConfig,
) {
	t.Helper()
	dyn, _, err := r.newDynamicClient()
	if err != nil {
		t.Errorf("failed to create dynamic client: %v", err)
		return
	}
	ri := dyn.Resource(schema.GroupVersionResource{
		Group:    k.Group,
		Version:  k.Version,
		Resource: k.Resource,
	})
	var client dynamic.ResourceInterface = ri
	if k.Namespace != "" {
		client = ri.Namespace(k.Namespace)
	}
	switch k.Action {
	case config.KubernetesActionGet:
		obj, err := client.Get(ctx, k.Name, resource.EmptyGetOptions)
		if err != nil {
			t.Errorf("failed to get custom resource %s/%s: %v", k.Resource, k.Name, err)
			return
		}
		log.Infof("custom resource object: %v", obj)
	case config.KubernetesActionDelete:
		if err := client.Delete(ctx, k.Name, resource.EmptyDeleteOptions); err != nil {
			if apierrors.IsNotFound(err) {
				log.Warnf("custom resource %s/%s is already deleted", k.Resource, k.Name)
				return
			}
			t.Errorf("failed to delete custom resource %s/%s: %v", k.Resource, k.Name, err)
		}
	case config.KubernetesActionWait:
		if err := waitForCustomResourceStatus(ctx, client, k); err != nil {
			t.Errorf("failed to wait for custom resource %s/%s: %v", k.Resource, k.Name, err)
		}
	default:
		t.Errorf("kubernetes action %s is not supported for kind %s", k.Action, k.Kind)
	}
}

// waitForCustomResourceStatus polls the custom resource until the jsonpath
// status_path evaluates to status_value. NotFound and missing-path results
// keep the poll running; ctx cancellation (execution timeout) stops it.
func waitForCustomResourceStatus(
	ctx context.Context, client dynamic.ResourceInterface, k *config.KubernetesConfig,
) error {
	jp := jsonpath.New(fieldManager).AllowMissingKeys(true)
	if err := jp.Parse(k.StatusPath); err != nil {
		return errors.Wrapf(err, "failed to parse status_path %s", k.StatusPath)
	}
	tick := time.NewTicker(customResourcePollInterval)
	defer tick.Stop()
	for {
		obj, err := client.Get(ctx, k.Name, resource.EmptyGetOptions)
		switch {
		case err == nil:
			var buf strings.Builder
			if err := jp.Execute(&buf, obj.UnstructuredContent()); err != nil {
				return errors.Wrapf(err, "failed to evaluate status_path %s", k.StatusPath)
			}
			got := buf.String()
			if got == k.StatusValue {
				log.Infof("custom resource %s/%s reached desired status %s=%s", k.Resource, k.Name, k.StatusPath, got)
				return nil
			}
			log.Infof("custom resource %s/%s status %s is %q, waiting for %q", k.Resource, k.Name, k.StatusPath, got, k.StatusValue)
		case apierrors.IsNotFound(err):
			log.Infof("custom resource %s/%s is not found yet, waiting", k.Resource, k.Name)
		default:
			return errors.Wrapf(err, "failed to get custom resource %s/%s", k.Resource, k.Name)
		}
		select {
		case <-ctx.Done():
			return errors.Wrapf(ctx.Err(), "timed out waiting for custom resource %s/%s to reach %s=%s",
				k.Resource, k.Name, k.StatusPath, k.StatusValue)
		case <-tick.C:
		}
	}
}

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
	"fmt"
	"math/rand/v2"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/k8s"
	kclient "github.com/vdaas/vald/internal/k8s/client"
	"github.com/vdaas/vald/internal/k8s/resource"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/tests/v2/e2e/config"
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

	rawCfg := r.k8s.GetRESTConfig()
	patcher, _ := kclient.NewPatcher(fieldManager)

	// The create action is special-cased because it always creates a job from
	// the cronjob named k.Name regardless of the requested kind.
	if k.Action == config.KubernetesActionCreate {
		if k.Kind != config.Job {
			t.Errorf("kubernetes action create is only supported for creating job from cronjob")
			return
		}
		c, err := kclient.NewWithConfig[*k8s.CronJob, *k8s.CronJobList](rawCfg, new(k8s.CronJob), new(k8s.CronJobList))
		if err != nil {
			t.Errorf("failed to create cronjob client: %v", err)
			return
		}
		cj, err := c.Get(ctx, k.Name, k.Namespace)
		if err != nil {
			t.Errorf("failed to get cronjob: %v", err)
			return
		}
		jc, err := kclient.NewWithConfig[*k8s.Job, *k8s.JobList](rawCfg, new(k8s.Job), new(k8s.JobList))
		if err != nil {
			t.Errorf("failed to create job client: %v", err)
			return
		}
		suffix := strconv.FormatInt(time.Now().UnixNano(), 36) + "-" + strconv.FormatUint(uint64(rand.Uint32()), 36)
		job := &k8s.Job{
			ObjectMeta: k8s.ObjectMeta{
				Name:      fmt.Sprintf("%s-%s", k.Name, suffix),
				Namespace: k.Namespace,
			},
			Spec: k8s.JobSpec{
				Template: cj.Spec.JobTemplate.Spec.Template,
			},
		}
		if err := jc.Create(ctx, job); err != nil {
			t.Errorf("failed to create job from cronjob: %v", err)
		}
		return
	}

	switch k.Kind {
	case config.ConfigMap:
		c, _ := kclient.NewWithConfig[*k8s.ConfigMap, *k8s.ConfigMapList](rawCfg, new(k8s.ConfigMap), new(k8s.ConfigMapList))
		handleKubernetesAction(t, ctx, k, c)
	case config.CronJob:
		c, _ := kclient.NewWithConfig[*k8s.CronJob, *k8s.CronJobList](rawCfg, new(k8s.CronJob), new(k8s.CronJobList))
		handleKubernetesWorkloadAction(t, ctx, k, c, patcher)
	case config.DaemonSet:
		c, _ := kclient.NewWithConfig[*k8s.DaemonSet, *k8s.DaemonSetList](rawCfg, new(k8s.DaemonSet), new(k8s.DaemonSetList))
		handleKubernetesWorkloadAction(t, ctx, k, c, patcher)
	case config.Deployment:
		c, _ := kclient.NewWithConfig[*k8s.Deployment, *k8s.DeploymentList](rawCfg, new(k8s.Deployment), new(k8s.DeploymentList))
		handleKubernetesWorkloadAction(t, ctx, k, c, patcher)
	case config.Job:
		c, _ := kclient.NewWithConfig[*k8s.Job, *k8s.JobList](rawCfg, new(k8s.Job), new(k8s.JobList))
		handleKubernetesWorkloadAction(t, ctx, k, c, patcher)
	case config.MutatingWebhookConfiguration:
		c, _ := kclient.NewWithConfig[*k8s.MutatingWebhookConfiguration, *k8s.MutatingWebhookConfigurationList](rawCfg, new(k8s.MutatingWebhookConfiguration), new(k8s.MutatingWebhookConfigurationList))
		handleKubernetesAction(t, ctx, k, c)
	case config.Pod:
		c, _ := kclient.NewWithConfig[*k8s.Pod, *k8s.PodList](rawCfg, new(k8s.Pod), new(k8s.PodList))
		handleKubernetesAction(t, ctx, k, c)
	case config.Secret:
		c, _ := kclient.NewWithConfig[*k8s.Secret, *k8s.SecretList](rawCfg, new(k8s.Secret), new(k8s.SecretList))
		handleKubernetesAction(t, ctx, k, c)
	case config.Service:
		c, _ := kclient.NewWithConfig[*k8s.Service, *k8s.ServiceList](rawCfg, new(k8s.Service), new(k8s.ServiceList))
		handleKubernetesAction(t, ctx, k, c)
	case config.StatefulSet:
		c, _ := kclient.NewWithConfig[*k8s.StatefulSet, *k8s.StatefulSetList](rawCfg, new(k8s.StatefulSet), new(k8s.StatefulSetList))
		handleKubernetesWorkloadAction(t, ctx, k, c, patcher)
	case config.ValidatingWebhookConfiguration:
		c, _ := kclient.NewWithConfig[*k8s.ValidatingWebhookConfiguration, *k8s.ValidatingWebhookConfigurationList](rawCfg, new(k8s.ValidatingWebhookConfiguration), new(k8s.ValidatingWebhookConfigurationList))
		handleKubernetesAction(t, ctx, k, c)
	default:
		t.Errorf("unsupported kubernetes kind: %s", k.Kind)
	}
}

// handleKubernetesAction executes the actions applicable to every resource kind
// (get, delete and wait) through the generic resource client.
func handleKubernetesAction[T k8s.Object, L k8s.ObjectList](
	t *testing.T, ctx context.Context, k *config.KubernetesConfig, client kclient.Client[T, L],
) {
	t.Helper()
	switch k.Action {
	case config.KubernetesActionGet:
		obj, err := client.Get(ctx, k.Name, k.Namespace)
		if err != nil {
			t.Errorf("failed to get %s: %v", k.Kind, err)
			return
		}
		log.Infof("kubernetes object: %v", obj)
	case config.KubernetesActionDelete:
		var tVar T
		obj := reflect.New(reflect.TypeFor[T]().Elem()).Interface().(T)
		obj.SetName(k.Name)
		obj.SetNamespace(k.Namespace)
		if err := client.Delete(ctx, obj); err != nil {
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
func handleKubernetesWorkloadAction[T k8s.Object, L k8s.ObjectList](
	t *testing.T, ctx context.Context, k *config.KubernetesConfig, client kclient.Client[T, L], patcher kclient.Patcher,
) {
	t.Helper()
	switch k.Action {
	case config.KubernetesActionRollout:
		if err := resource.RolloutRestart(ctx, patcher, k.Name, k.Namespace); err != nil {
			t.Errorf("failed to rollout restart %s: %v", k.Kind, err)
		}
	default:
		handleKubernetesAction(t, ctx, k, client)
	}
}

// decodeManifest reads a multi-document yaml manifest file and decodes each
// document into an unstructured object. Empty documents are skipped.
func decodeManifest(path string) (objs []*k8s.Unstructured, err error) {
	data, err := file.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read manifest file %s", path)
	}
	dec := k8s.NewYAMLOrJSONDecoder(bytes.NewReader(data), manifestDecodeBufferSize)
	for {
		obj := new(k8s.Unstructured)
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
	dyn kclient.Dynamic, mapper kclient.RESTMapper, obj *k8s.Unstructured, fallbackNamespace string,
) (kclient.DynamicResource, error) {
	gvk := obj.GroupVersionKind()
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to resolve rest mapping for %s", gvk.String())
	}
	if mapping.Scope.Name() != kclient.RESTScopeNameNamespace {
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
	dyn, mapper, err := kclient.NewDynamicClient(r.k8s)
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
			if _, err := client.Patch(ctx, obj.GetName(), k8s.ApplyPatchType, data, k8s.PatchOptions{
				FieldManager: fieldManager,
				Force:        &force,
			}); err != nil {
				t.Errorf("failed to apply %s %s: %v", obj.GetKind(), obj.GetName(), err)
				return
			}
			log.Infof("applied %s %s from manifest %s", obj.GetKind(), obj.GetName(), k.Manifest)
		case config.KubernetesActionDelete:
			if err := client.Delete(ctx, obj.GetName(), resource.EmptyDeleteOptions); err != nil {
				if k8s.IsNotFound(err) {
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
	dyn, _, err := kclient.NewDynamicClient(r.k8s)
	if err != nil {
		t.Errorf("failed to create dynamic client: %v", err)
		return
	}
	ri := dyn.Resource(k8s.GroupVersionResource{
		Group:    k.Group,
		Version:  k.Version,
		Resource: k.Resource,
	})
	var client kclient.DynamicResource = ri
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
			if k8s.IsNotFound(err) {
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
	ctx context.Context, client kclient.DynamicResource, k *config.KubernetesConfig,
) error {
	jp := k8s.NewJSONPath(fieldManager).AllowMissingKeys(true)
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
		case k8s.IsNotFound(err):
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

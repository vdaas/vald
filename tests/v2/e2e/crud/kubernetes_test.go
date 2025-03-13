//go:build e2e

//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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
	"context"
	"testing"

	"github.com/vdaas/vald/tests/v2/e2e/config"
	"github.com/vdaas/vald/tests/v2/e2e/kubernetes"
)

func (r *runner) proscessKubernetes(t *testing.T, ctx context.Context, plan *config.Execution) {
	t.Helper()
	if plan == nil {
		t.Fatal("kubernetes plan is nil")
		return
	}
	switch plan.Type {
	case config.OpKubernetes:
		if plan.Kubernetes == nil {
			t.Fatal("kubernetes field is nil")
			return
		}
		switch plan.Kubernetes.Kind {
		case config.ConfigMap:
			client := kubernetes.ConfigMap(r.k8s, plan.Kubernetes.Namespace)
			resourceKubernetes(t, ctx, client, plan)
		case config.CronJob:
			client := kubernetes.CronJob(r.k8s, plan.Kubernetes.Namespace)
			executeKubernetes(t, ctx, client, plan)
		case config.DaemonSet:
			client := kubernetes.DaemonSet(r.k8s, plan.Kubernetes.Namespace)
			executeKubernetes(t, ctx, client, plan)
		case config.Deployment:
			client := kubernetes.Deployment(r.k8s, plan.Kubernetes.Namespace)
			executeKubernetes(t, ctx, client, plan)
		case config.Job:
			client := kubernetes.Job(r.k8s, plan.Kubernetes.Namespace)
			executeKubernetes(t, ctx, client, plan)
		case config.Pod:
			client := kubernetes.Pod(r.k8s, plan.Kubernetes.Namespace)
			resourceKubernetes(t, ctx, client, plan)
		case config.Secret:
			client := kubernetes.Secret(r.k8s, plan.Kubernetes.Namespace)
			resourceKubernetes(t, ctx, client, plan)
		case config.Service:
			client := kubernetes.Service(r.k8s, plan.Kubernetes.Namespace)
			resourceKubernetes(t, ctx, client, plan)
		case config.StatefulSet:
			client := kubernetes.StatefulSet(r.k8s, plan.Kubernetes.Namespace)
			executeKubernetes(t, ctx, client, plan)
		}
	default:
		t.Errorf("unsupported kubernetes operation: %s", plan.Type)
	}
}

func executeKubernetes[
	T kubernetes.Object,
	L kubernetes.ObjectList,
	C kubernetes.NamedObject,
	I kubernetes.WorkloadControllerResourceClient[T, L, C],
](
	t *testing.T,
	ctx context.Context,
	client I,
	plan *config.Execution,
) {
	t.Helper()
	switch plan.Kubernetes.Action {
	case config.KubernetesActionRollout:
		err := kubernetes.RolloutRestart(ctx, client, plan.Kubernetes.Name)
		if err != nil {
			t.Errorf("failed to %s %s: name = %s, namespace = %s: %#v", plan.Kubernetes.Action, plan.Kubernetes.Kind, plan.Kubernetes.Namespace, plan.Kubernetes.Name, err)
			return
		}
		return
	case config.KubernetesActionWait:
		obj, _, err := kubernetes.WaitForStatus(ctx, client, plan.Kubernetes.Name, plan.Kubernetes.Status.Status())
		if err != nil {
			t.Errorf("failed to %s %s: name = %s, namespace = %s: %#v", plan.Kubernetes.Action, plan.Kubernetes.Kind, plan.Kubernetes.Namespace, plan.Kubernetes.Name, err)
			return
		}
		t.Logf("complete wait for statefulset ready: %#v", obj)
		return
	}
	return
}

func resourceKubernetes[
	T kubernetes.Object,
	L kubernetes.ObjectList,
	C kubernetes.NamedObject,
	I kubernetes.ResourceClient[T, L, C],
](
	t *testing.T,
	ctx context.Context,
	client I,
	plan *config.Execution,
) {
	t.Helper()
	// TODO:
	return
}

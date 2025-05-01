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

	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/tests/v2/e2e/config"
	"github.com/vdaas/vald/tests/v2/e2e/kubernetes"
)

func (r *runner) processKubernetes(t *testing.T, ctx context.Context, plan *config.Execution) {
	t.Helper()
	if plan == nil || plan.Kubernetes == nil {
		t.Fatal("kubernetes plan is nil")
		return
	}
	var err error
	switch plan.Kubernetes.Action {
	case config.KubernetesActionRollout:
		switch plan.Kubernetes.Kind {
		case config.Deployment:
			err = kubernetes.RolloutRestart(ctx, kubernetes.Deployment(r.k8s, plan.Kubernetes.Namespace), plan.Kubernetes.Name)
		case config.DaemonSet:
			err = kubernetes.RolloutRestart(ctx, kubernetes.DaemonSet(r.k8s, plan.Kubernetes.Namespace), plan.Kubernetes.Name)
		case config.StatefulSet:
			err = kubernetes.RolloutRestart(ctx, kubernetes.StatefulSet(r.k8s, plan.Kubernetes.Namespace), plan.Kubernetes.Name)
		case config.CronJob:
			err = kubernetes.RolloutRestart(ctx, kubernetes.CronJob(r.k8s, plan.Kubernetes.Namespace), plan.Kubernetes.Name)
		case config.Job:
			err = kubernetes.RolloutRestart(ctx, kubernetes.Job(r.k8s, plan.Kubernetes.Namespace), plan.Kubernetes.Name)
		default:
		}
		if err != nil {
			t.Errorf("failed to rollout restart %s: %v", plan.Kubernetes.Kind, err)
		}
	case config.KubernetesActionDelete:
		switch plan.Kubernetes.Kind {
		case config.ConfigMap:
			err = kubernetes.ConfigMap(r.k8s, plan.Kubernetes.Namespace).Delete(ctx, plan.Kubernetes.Name, kubernetes.EmptyDeleteOptions)
		case config.CronJob:
			err = kubernetes.CronJob(r.k8s, plan.Kubernetes.Namespace).Delete(ctx, plan.Kubernetes.Name, kubernetes.EmptyDeleteOptions)
		case config.DaemonSet:
			err = kubernetes.DaemonSet(r.k8s, plan.Kubernetes.Namespace).Delete(ctx, plan.Kubernetes.Name, kubernetes.EmptyDeleteOptions)
		case config.Deployment:
			err = kubernetes.Deployment(r.k8s, plan.Kubernetes.Namespace).Delete(ctx, plan.Kubernetes.Name, kubernetes.EmptyDeleteOptions)
		case config.Job:
			err = kubernetes.Job(r.k8s, plan.Kubernetes.Namespace).Delete(ctx, plan.Kubernetes.Name, kubernetes.EmptyDeleteOptions)
		case config.Pod:
			err = kubernetes.Pod(r.k8s, plan.Kubernetes.Namespace).Delete(ctx, plan.Kubernetes.Name, kubernetes.EmptyDeleteOptions)
		case config.Secret:
			err = kubernetes.Secret(r.k8s, plan.Kubernetes.Namespace).Delete(ctx, plan.Kubernetes.Name, kubernetes.EmptyDeleteOptions)
		case config.Service:
			err = kubernetes.Service(r.k8s, plan.Kubernetes.Namespace).Delete(ctx, plan.Kubernetes.Name, kubernetes.EmptyDeleteOptions)
		case config.StatefulSet:
			err = kubernetes.StatefulSet(r.k8s, plan.Kubernetes.Namespace).Delete(ctx, plan.Kubernetes.Name, kubernetes.EmptyDeleteOptions)
		default:
		}
		if err != nil {
			t.Errorf("failed to delete %s: %v", plan.Kubernetes.Kind, err)
		}
	case config.KubernetesActionGet:
		var obj kubernetes.Object
		switch plan.Kubernetes.Kind {
		case config.ConfigMap:
			obj, err = kubernetes.ConfigMap(r.k8s, plan.Kubernetes.Namespace).Get(ctx, plan.Kubernetes.Name, kubernetes.EmptyGetOptions)
		case config.CronJob:
			obj, err = kubernetes.CronJob(r.k8s, plan.Kubernetes.Namespace).Get(ctx, plan.Kubernetes.Name, kubernetes.EmptyGetOptions)
		case config.DaemonSet:
			obj, err = kubernetes.DaemonSet(r.k8s, plan.Kubernetes.Namespace).Get(ctx, plan.Kubernetes.Name, kubernetes.EmptyGetOptions)
		case config.Deployment:
			obj, err = kubernetes.Deployment(r.k8s, plan.Kubernetes.Namespace).Get(ctx, plan.Kubernetes.Name, kubernetes.EmptyGetOptions)
		case config.Job:
			obj, err = kubernetes.Job(r.k8s, plan.Kubernetes.Namespace).Get(ctx, plan.Kubernetes.Name, kubernetes.EmptyGetOptions)
		case config.Pod:
			obj, err = kubernetes.Pod(r.k8s, plan.Kubernetes.Namespace).Get(ctx, plan.Kubernetes.Name, kubernetes.EmptyGetOptions)
		case config.Secret:
			obj, err = kubernetes.Secret(r.k8s, plan.Kubernetes.Namespace).Get(ctx, plan.Kubernetes.Name, kubernetes.EmptyGetOptions)
		case config.Service:
			obj, err = kubernetes.Service(r.k8s, plan.Kubernetes.Namespace).Get(ctx, plan.Kubernetes.Name, kubernetes.EmptyGetOptions)
		case config.StatefulSet:
			obj, err = kubernetes.StatefulSet(r.k8s, plan.Kubernetes.Namespace).Get(ctx, plan.Kubernetes.Name, kubernetes.EmptyGetOptions)
		default:
		}
		if err != nil {
			t.Errorf("failed to get %s: %v", plan.Kubernetes.Kind, err)
		}
		if obj != nil {
			log.Infof("kubernetes object: %v", obj)
		}
	case config.KubernetesActionWait:
		var ok bool
		switch plan.Kubernetes.Kind {
		case config.ConfigMap:
			ok, err = kubernetes.WaitForStatus(ctx, kubernetes.ConfigMap(r.k8s, plan.Kubernetes.Namespace), plan.Kubernetes.Name, plan.Kubernetes.Selector, plan.Kubernetes.Status.Status())
		case config.CronJob:
			ok, err = kubernetes.WaitForStatus(ctx, kubernetes.CronJob(r.k8s, plan.Kubernetes.Namespace), plan.Kubernetes.Name, plan.Kubernetes.Selector, plan.Kubernetes.Status.Status())
		case config.DaemonSet:
			ok, err = kubernetes.WaitForStatus(ctx, kubernetes.DaemonSet(r.k8s, plan.Kubernetes.Namespace), plan.Kubernetes.Name, plan.Kubernetes.Selector, plan.Kubernetes.Status.Status())
		case config.Deployment:
			ok, err = kubernetes.WaitForStatus(ctx, kubernetes.Deployment(r.k8s, plan.Kubernetes.Namespace), plan.Kubernetes.Name, plan.Kubernetes.Selector, plan.Kubernetes.Status.Status())
		case config.Job:
			ok, err = kubernetes.WaitForStatus(ctx, kubernetes.Job(r.k8s, plan.Kubernetes.Namespace), plan.Kubernetes.Name, plan.Kubernetes.Selector, plan.Kubernetes.Status.Status())
		case config.Pod:
			ok, err = kubernetes.WaitForStatus(ctx, kubernetes.Pod(r.k8s, plan.Kubernetes.Namespace), plan.Kubernetes.Name, plan.Kubernetes.Selector, plan.Kubernetes.Status.Status())
		case config.Secret:
			ok, err = kubernetes.WaitForStatus(ctx, kubernetes.Secret(r.k8s, plan.Kubernetes.Namespace), plan.Kubernetes.Name, plan.Kubernetes.Selector, plan.Kubernetes.Status.Status())
		case config.Service:
			ok, err = kubernetes.WaitForStatus(ctx, kubernetes.Service(r.k8s, plan.Kubernetes.Namespace), plan.Kubernetes.Name, plan.Kubernetes.Selector, plan.Kubernetes.Status.Status())
		case config.StatefulSet:
			ok, err = kubernetes.WaitForStatus(ctx, kubernetes.StatefulSet(r.k8s, plan.Kubernetes.Namespace), plan.Kubernetes.Name, plan.Kubernetes.Selector, plan.Kubernetes.Status.Status())
		default:
		}
		if !ok {
			t.Errorf("failed to wait for %s: %v", plan.Kubernetes.Kind, err)
		}
	case config.KubernetesActionExec:
		t.Errorf("kubernetes action %s is not supported yet", plan.Kubernetes.Action)
	case config.KubernetesActionApply:
		t.Errorf("kubernetes action %s is not supported yet", plan.Kubernetes.Action)
	case config.KubernetesActionCreate:
		t.Errorf("kubernetes action %s is not supported yet", plan.Kubernetes.Action)
	case config.KubernetesActionPatch:
		t.Errorf("kubernetes action %s is not supported yet", plan.Kubernetes.Action)
	case config.KubernetesActionScale:
		t.Errorf("kubernetes action %s is not supported yet", plan.Kubernetes.Action)
	default:
		t.Errorf("kubernetes action %s is not supported yet", plan.Kubernetes.Action)
	}
	return
}

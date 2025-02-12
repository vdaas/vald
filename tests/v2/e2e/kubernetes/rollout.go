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

// package kubernetes provides kubernetes e2e tests
package kubernetes

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
)

const (
	rolloutAnnotationKey = "kubectl.kubernetes.io/restartedAt"
)

func (c *client) RolloutRestartDeployment(wait bool) (err error) {
	res := Deployment{
		Name:      "something",
		Namespace: "default",
	}
	err = RolloutRestart(context.TODO(), c.clientset, res)
	if err != nil {
		return err
	}
	if !wait {
		return nil
	}
	return WaitForRestart(context.TODO(), c.clientset, res)
}

func RolloutRestart[R metav1.Object, T ResourceType[R]](
	ctx context.Context,
	clientset *kubernetes.Clientset,
	rt T,
) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		obj, err := rt.Get(ctx, clientset)
		if err != nil {
			return err
		}
		annotations := obj.GetAnnotations()
		if annotations == nil {
			annotations = make(map[string]string, 1)
		}

		annotations[rolloutAnnotationKey] = time.Now().UTC().Format(time.RFC3339)

		obj.SetAnnotations(annotations)

		_, err = rt.Update(ctx, clientset, obj)
		return err
	})
}

func WaitForRestart[R metav1.Object, T ResourceType[R]](
	ctx context.Context,
	clientset *kubernetes.Clientset,
	rt T,
) error {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			obj, err := rt.Get(ctx, clientset)
			if err != nil {
				return err
			}
			done, info, err := checkRolloutStatus(obj)
			if err != nil {
				return err
			}
			fmt.Println(info)
			if done {
				return nil
			}
		}
	}
}

func checkRolloutStatus[T any](obj T) (bool, string, error) {
	switch r := any(obj).(type) {
	case *appsv1.Deployment:
		desired := int32(1)
		if r.Spec.Replicas != nil {
			desired = *r.Spec.Replicas
		}
		if r.Status.UpdatedReplicas == desired &&
			r.Status.AvailableReplicas == desired &&
			r.Status.ObservedGeneration >= r.Generation {
			return true, fmt.Sprintf("Deployment rollout complete: %d/%d updated, %d available", r.Status.UpdatedReplicas, desired, r.Status.AvailableReplicas), nil
		}
		return false, fmt.Sprintf("Deployment rollout in progress: %d/%d updated, %d available", r.Status.UpdatedReplicas, desired, r.Status.AvailableReplicas), nil
	case *appsv1.StatefulSet:
		desired := int32(1)
		if r.Spec.Replicas != nil {
			desired = *r.Spec.Replicas
		}
		if r.Status.UpdatedReplicas == desired &&
			r.Status.ReadyReplicas == desired &&
			r.Status.CurrentReplicas == desired {
			return true, fmt.Sprintf("StatefulSet rollout complete: %d/%d updated, %d ready", r.Status.UpdatedReplicas, desired, r.Status.ReadyReplicas), nil
		}
		return false, fmt.Sprintf("StatefulSet rollout in progress: %d/%d updated, %d ready", r.Status.UpdatedReplicas, desired, r.Status.ReadyReplicas), nil
	case *appsv1.DaemonSet:
		if r.Status.UpdatedNumberScheduled == r.Status.DesiredNumberScheduled &&
			r.Status.NumberAvailable == r.Status.DesiredNumberScheduled {
			return true, fmt.Sprintf("DaemonSet rollout complete: %d/%d updated, %d available", r.Status.UpdatedNumberScheduled, r.Status.DesiredNumberScheduled, r.Status.NumberAvailable), nil
		}
		return false, fmt.Sprintf("DaemonSet rollout in progress: %d/%d updated, %d available", r.Status.UpdatedNumberScheduled, r.Status.DesiredNumberScheduled, r.Status.NumberAvailable), nil
	default:
		return false, "", fmt.Errorf("unsupported resource type in checkRolloutStatus")
	}
}

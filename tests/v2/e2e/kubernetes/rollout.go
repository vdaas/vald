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

package kubernetes

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
)

const (
	rolloutAnnotationKey = "kubectl.kubernetes.io/restartedAt"
)

// RolloutRestart restarts a kubernetes resources (Deployment, DaemonSet, StatefulSet, Job, CronJob).
//
// # Example
// ```go
//
//	client, err := kubernetes.NewClient("/path/to/kubeconfig", "current context") // create a kubernetes client
//	if err != nil {
//		return err
//	}
//
//	deploymentClient := kubernetes.Deployment(client, "default") // create a deployment client
//	err = kubernetes.RolloutRestart(ctx, deploymentClient, "some deployment") // restart the deployment
//	if err != nil {
//		return err
//	}
//
//	statefulSetClient := kubernetes.StatefulSet(client, "default") // create a statefulset client
//	err = kubernetes.RolloutRestart(ctx, statefulSetClient, "some statefulset") // restart the statefulset
//	if err != nil {
//		return err
//	}
//
// ```
func RolloutRestart[T Object, L ObjectList, C NamedObject, I WorkloadControllerResourceClient[T, L, C]](
	ctx context.Context, client I, name string,
) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() (err error) {
		_, err = client.SetPodAnnotations(ctx, name, map[string]string{
			rolloutAnnotationKey: time.Now().UTC().Format(time.RFC3339),
		}, metav1.GetOptions{}, metav1.UpdateOptions{})
		return err
	})
}

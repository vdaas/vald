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
	"encoding/json"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
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
func RolloutRestart[T Object, L ObjectList, C NamedObject, I ObjectInterface[T, L, C]](
	ctx context.Context, client I, name string,
) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		_, err := client.Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return err
		}

		// For RolloutRestart, the annotation under the "T.spec.template.metadata" should be updated.
		patchData := map[string]interface{}{
			"spec": map[string]interface{}{
				"template": map[string]interface{}{
					"metadata": map[string]interface{}{
						"annotations": map[string]string{
							rolloutAnnotationKey: time.Now().UTC().Format(time.RFC3339),
						},
					},
				},
			},
		}
		patchBytes, err := json.Marshal(patchData)
		if err != nil {
			return err
		}

		_, err = client.Patch(ctx, name, types.MergePatchType, patchBytes, metav1.PatchOptions{})
		return err
	})
}

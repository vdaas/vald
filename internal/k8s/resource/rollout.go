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

package resource

import (
	"context"
	"time"

	"github.com/vdaas/vald/internal/k8s/client"
	"k8s.io/client-go/util/retry"
)

const (
	rolloutAnnotationKey = "kubectl.kubernetes.io/restartedAt"
)

// RolloutRestart triggers a rolling restart of the named workload by stamping
// the kubectl restartedAt annotation onto its pod template.
func RolloutRestart(
	ctx context.Context, patcher client.Patcher, name, namespace string,
) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() (err error) {
		err = patcher.ApplyPodAnnotations(ctx, name, namespace, map[string]string{
			rolloutAnnotationKey: time.Now().UTC().Format(time.RFC3339),
		})
		return err
	})
}

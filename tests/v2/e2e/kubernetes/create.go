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

package kubernetes

import (
	"context"
	"fmt"
	"time"

	"github.com/vdaas/vald/internal/errors"
	batchv1 "k8s.io/api/batch/v1" // For Job and CronJob
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
)

func createJobFromCronJob(
	ctx context.Context, j JobClient, c CronJobClient, from string,
) (*batchv1.Job, error) {
	cobj, err := c.Get(ctx, from, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	tmpl, err := c.GetPodTemplate(cobj)
	if err != nil {
		return nil, err
	}
	if tmpl == nil {
		return nil, errors.ErrPodTemplateNotFound
	}
	t := time.Now()
	return j.Create(ctx, &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", from, t.Format("2006-01-02-15-04-05")),
			Namespace: cobj.Namespace,
		},
		Spec: batchv1.JobSpec{
			Template: *tmpl,
		},
	}, metav1.CreateOptions{})
}

// Create creates a kubernetes resources (Only Job from CronJob is supported).
//
// # Example
// ```go
//
//	client, err := kubernetes.NewClient("/path/to/kubeconfig", "current context") // create a kubernetes client
//	if err != nil {
//		return err
//	}
//
//	jobClient := kubernetes.Job(client, "default") // create a job client
//	cronClient := kubernetes.CronJob(client, "default") // create a cronjob client
//	err = kubernetes.Create(ctx, jobClient, cronClient, "cronjob name") // create a job from cronjob
//	if err != nil {
//		return err
//	}
//
// ```
func Create(ctx context.Context, jobClient JobClient, cronClient CronJobClient, from string) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() (err error) {
		_, err = createJobFromCronJob(ctx, jobClient, cronClient, from)
		return err
	})
}

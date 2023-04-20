//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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

// Package job manages the main logic of benchmark job.
package job

import (
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
)

type benchmarkJobTemplate = batchv1.Job

const (
	SvcAccountName = "vald-benchmark-operator"
	ContainerName  = "vald-benchmark-job"
	ContainerImage = "local-registry:5000/vdaas/vald-benchmark-job:latest"

	RestartPolicyAlways    corev1.RestartPolicy = "Always"
	RestartPolicyOnFailure corev1.RestartPolicy = "OnFailure"
	RestartPolicyNever     corev1.RestartPolicy = "Never"
)

// NewBenchmarkJobTemplate creates the job template for crating k8s job resource.
func NewBenchmarkJobTemplate(opts ...BenchmarkJobOption) (benchmarkJobTemplate, error) {
	jobTmpl := new(benchmarkJobTemplate)
	for _, opt := range append(defaultBenchmarkJobOpts, opts...) {
		err := opt(jobTmpl)
		if err != nil {
			return *jobTmpl, err
		}
	}
	jobTmpl.Spec.Template.Spec.Containers = []corev1.Container{
		{
			Name:            ContainerName,
			Image:           ContainerImage,
			ImagePullPolicy: corev1.PullAlways,
			LivenessProbe: &corev1.Probe{
				InitialDelaySeconds: int32(60),
				PeriodSeconds:       int32(10),
				TimeoutSeconds:      int32(300),
				ProbeHandler: corev1.ProbeHandler{
					Exec: &corev1.ExecAction{
						Command: []string{
							"/go/bin/job",
							"-v",
						},
					},
				},
			},
			StartupProbe: &corev1.Probe{
				FailureThreshold: int32(30),
				PeriodSeconds:    int32(10),
				TimeoutSeconds:   int32(300),
				ProbeHandler: corev1.ProbeHandler{
					Exec: &corev1.ExecAction{
						Command: []string{
							"/go/bin/job",
							"-v",
						},
					},
				},
			},
			Ports: []corev1.ContainerPort{
				{
					Name:          "liveness",
					Protocol:      corev1.ProtocolTCP,
					ContainerPort: int32(3000),
				},
				{
					Name:          "readiness",
					Protocol:      corev1.ProtocolTCP,
					ContainerPort: int32(3001),
				},
			},
			Env: []corev1.EnvVar{
				{
					Name: "CRD_NAMESPACE",
					ValueFrom: &corev1.EnvVarSource{
						FieldRef: &corev1.ObjectFieldSelector{
							FieldPath: "metadata.namespace",
						},
					},
				},
				{
					Name: "CRD_NAME",
					ValueFrom: &corev1.EnvVarSource{
						FieldRef: &corev1.ObjectFieldSelector{
							FieldPath: "metadata.labels['job-name']",
						},
					},
				},
			},
		},
	}
	return *jobTmpl, nil
}

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
	jobs "github.com/vdaas/vald/internal/k8s/job"
	corev1 "k8s.io/api/core/v1"
)

type (
	ImagePullPolicy corev1.PullPolicy
	RestartPolicy   corev1.RestartPolicy
)

const (
	PullAlways       ImagePullPolicy = "Always"
	PullNever        ImagePullPolicy = "Never"
	PullIfNotPresent ImagePullPolicy = "PullIfNotPresent"

	RestartPolicyAlways    RestartPolicy = "Always"
	RestartPolicyOnFailure RestartPolicy = "OnFailure"
	RestartPolicyNever     RestartPolicy = "Never"
)

const (
	svcAccount = "vald-benchmark-operator"
)

type BenchmarkJobTpl interface {
	CreateJobTpl(opts ...BenchmarkJobOption) (jobs.Job, error)
}

type benchmarkJobTpl struct {
	containerName      string
	containerImageName string
	imagePullPolicy    ImagePullPolicy
	jobTpl             jobs.Job
}

func NewBenchmarkJob(opts ...BenchmarkJobTplOption) (BenchmarkJobTpl, error) {
	bjTpl := new(benchmarkJobTpl)
	for _, opt := range append(defaultBenchmarkJobTplOpts, opts...) {
		err := opt(bjTpl)
		if err != nil {
			return nil, err
		}
	}
	return bjTpl, nil
}

func (b *benchmarkJobTpl) CreateJobTpl(opts ...BenchmarkJobOption) (jobs.Job, error) {
	for _, opt := range append(defaultBenchmarkJobOpts, opts...) {
		err := opt(&b.jobTpl)
		if err != nil {
			return b.jobTpl, err
		}
	}
	// TODO: check enable pprof flag
	b.jobTpl.Spec.Template.Annotations = map[string]string{
		"pyroscope.io/scrape":              "true",
		"pyroscope.io/application-name":    "benchmark-job",
		"pyroscope.io/profile-cpu-enabled": "true",
		"pyroscope.io/profile-mem-enabled": "true",
		"pyroscope.io/port":                "6060",
	}
	b.jobTpl.Spec.Template.Spec.Containers = []corev1.Container{
		{
			Name:            b.containerName,
			Image:           b.containerImageName,
			ImagePullPolicy: corev1.PullPolicy(b.imagePullPolicy),
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
	return b.jobTpl, nil
}

// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package job manages the main logic of benchmark job.
package job

import (
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/k8s/vald"
	corev1 "k8s.io/api/core/v1"
)

type (
	ImagePullPolicy corev1.PullPolicy
	RestartPolicy   corev1.RestartPolicy
)

const (
	PullAlways ImagePullPolicy = "Always"

	RestartPolicyNever RestartPolicy = "Never"

	volumeName = "vald-benchmark-job-config"
	svcAccount = "vald-benchmark-operator"

	configVolumeDefaultMode = 420

	pyroscopeScrapeEnabled = "true"

	livenessProbeInitialDelaySeconds = 60
	livenessProbePeriodSeconds       = 10
	livenessProbeTimeoutSeconds      = 300

	startupProbeFailureThreshold = 30
	startupProbePeriodSeconds    = 10
	startupProbeTimeoutSeconds   = 300

	containerPortLiveness  = 3000
	containerPortReadiness = 3001
)

type BenchmarkTpl interface {
	CreateJobTpl(opts ...BenchmarkOption) (k8s.Job, error)
}

// fieldRefEnvVar builds an EnvVar sourced from the downward-API field at
// fieldPath.
func fieldRefEnvVar(name, fieldPath string) corev1.EnvVar {
	return corev1.EnvVar{
		Name: name,
		ValueFrom: &corev1.EnvVarSource{
			FieldRef: &corev1.ObjectFieldSelector{
				FieldPath: fieldPath,
			},
		},
	}
}

type benchmarkJobTemplate struct {
	jobTpl             k8s.Job
	containerName      string
	containerImageName string
	configMapName      string
	imagePullPolicy    ImagePullPolicy
}

// NewBenchmarkJob builds a benchmark job template from the given options.
// Option application errors abort construction only when critical (an empty
// required field such as the container name or image); any other option error
// is logged as a warning and skipped.
func NewBenchmarkJob(opts ...BenchmarkTemplateOption) (BenchmarkTpl, error) {
	template := new(benchmarkJobTemplate)
	for _, opt := range append(defaultBenchmarkJobTemplateOptions, opts...) {
		if err := opt(template); err != nil {
			if abort, oerr := vald.SkipNonCriticalOptionError(err, opt); abort {
				return nil, oerr
			}
		}
	}
	return template, nil
}

// CreateJobTpl materializes the k8s Job manifest from the template and the
// given options, with the same severity handling as NewBenchmarkJob: only
// critical option errors (an empty job name) abort, any other option error is
// logged as a warning and skipped.
func (b *benchmarkJobTemplate) CreateJobTpl(opts ...BenchmarkOption) (k8s.Job, error) {
	for _, opt := range append(defaultBenchmarkJobOpts, opts...) {
		if err := opt(&b.jobTpl); err != nil {
			if abort, oerr := vald.SkipNonCriticalOptionError(err, opt); abort {
				return b.jobTpl, oerr
			}
		}
	}
	// TODO: check enable pprof flag
	b.jobTpl.Spec.Template.Annotations = map[string]string{
		"pyroscope.io/scrape":              pyroscopeScrapeEnabled,
		"pyroscope.io/application-name":    "benchmark-job",
		"pyroscope.io/profile-cpu-enabled": pyroscopeScrapeEnabled,
		"pyroscope.io/profile-mem-enabled": pyroscopeScrapeEnabled,
		"pyroscope.io/port":                "6060",
	}
	b.jobTpl.Spec.Template.Spec.Containers = []corev1.Container{
		{
			Name:            b.containerName,
			Image:           b.containerImageName,
			ImagePullPolicy: corev1.PullPolicy(b.imagePullPolicy),
			LivenessProbe: &corev1.Probe{
				InitialDelaySeconds: int32(livenessProbeInitialDelaySeconds),
				PeriodSeconds:       int32(livenessProbePeriodSeconds),
				TimeoutSeconds:      int32(livenessProbeTimeoutSeconds),
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
				FailureThreshold: int32(startupProbeFailureThreshold),
				PeriodSeconds:    int32(startupProbePeriodSeconds),
				TimeoutSeconds:   int32(startupProbeTimeoutSeconds),
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
					ContainerPort: int32(containerPortLiveness),
				},
				{
					Name:          "readiness",
					Protocol:      corev1.ProtocolTCP,
					ContainerPort: int32(containerPortReadiness),
				},
			},
			Env: []corev1.EnvVar{
				fieldRefEnvVar("CRD_NAMESPACE", "metadata.namespace"),
				fieldRefEnvVar("CRD_NAME", "metadata.labels['job-name']"),
				fieldRefEnvVar("MY_NODE_NAME", "spec.nodeName"),
				fieldRefEnvVar("MY_POD_NAMESPACE", "metadata.namespace"),
				fieldRefEnvVar("MY_POD_NAME", "metadata.name"),
			},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      volumeName,
					MountPath: "/etc/server",
				},
			},
		},
	}
	// mount benchmark operator config map.
	// It is used for bind only observability config for each benchmark job
	mode := int32(configVolumeDefaultMode)
	b.jobTpl.Spec.Template.Spec.Volumes = []corev1.Volume{
		{
			Name: volumeName,
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: b.configMapName,
					},
					DefaultMode: &mode,
				},
			},
		},
	}
	return b.jobTpl, nil
}

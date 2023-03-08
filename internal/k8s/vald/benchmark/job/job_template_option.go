//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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

package job

import (
	"github.com/vdaas/vald/internal/k8s"
	corev1 "k8s.io/api/core/v1"
)

// BenchmarkJobOption represents the option for create benchmark job template.
type BenchmarkJobOption func(b *benchmarkJobTemplate) error

var defaultBenchmarkJobOpts = []BenchmarkJobOption{
	WithSvcAccountName(SvcAccountName),
	WithRestartPolicy(RestartPolicyNever),
}

// WithSvcAccountName sets the service account name for benchmark job.
func WithSvcAccountName(name string) BenchmarkJobOption {
	return func(b *benchmarkJobTemplate) error {
		if len(name) > 0 {
			b.Spec.Template.Spec.ServiceAccountName = name
		}
		return nil
	}
}

// WithRestartPolicy sets the job restart policy for benchmark job.
func WithRestartPolicy(rp corev1.RestartPolicy) BenchmarkJobOption {
	return func(b *benchmarkJobTemplate) error {
		if len(rp) > 0 {
			b.Spec.Template.Spec.RestartPolicy = rp
		}
		return nil
	}
}

// WithBackoffLimit sets the job backoff limit for benchmark job.
func WithBackoffLimit(bo int32) BenchmarkJobOption {
	return func(b *benchmarkJobTemplate) error {
		b.Spec.BackoffLimit = &bo
		return nil
	}
}

// WithName sets the job name of benchmark job.
func WithName(name string) BenchmarkJobOption {
	return func(b *benchmarkJobTemplate) error {
		b.Name = name
		return nil
	}
}

// WithNamespace specify namespace where job will execute.
func WithNamespace(ns string) BenchmarkJobOption {
	return func(b *benchmarkJobTemplate) error {
		b.Namespace = ns
		return nil
	}
}

// WithOwnerRef sets the OwnerReference to the job resource.
func WithOwnerRef(refs []k8s.OwnerReference) BenchmarkJobOption {
	return func(b *benchmarkJobTemplate) error {
		if len(refs) > 0 {
			b.OwnerReferences = refs
		}
		return nil
	}
}

// WithCompletions sets the job completion.
func WithCompletions(com int32) BenchmarkJobOption {
	return func(b *benchmarkJobTemplate) error {
		if com > 1 {
			b.Spec.Completions = &com
		}
		return nil
	}
}

// WithParallelism sets the job parallelism.
func WithParallelism(parallelism int32) BenchmarkJobOption {
	return func(b *benchmarkJobTemplate) error {
		if parallelism > 1 {
			b.Spec.Parallelism = &parallelism
		}
		return nil
	}
}

// WithLabel sets the label to the job resource.
func WithLabel(label map[string]string) BenchmarkJobOption {
	return func(b *benchmarkJobTemplate) error {
		if len(label) > 0 {
			b.Labels = label
		}
		return nil
	}
}
